package bme280

import (
	"errors"
	"fmt"
	"time"

	// Package imports
	event "github.com/djthorpe/go-pico/pkg/event"
	multierror "github.com/hashicorp/go-multierror"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	SPI   SPI
	I2C   I2C
	Slave uint8 // I2C Slave address, optional
}

type device struct {
	i2c                    I2C
	spi                    SPI
	slave                  uint8
	chipid, version        uint8
	mode                   Mode
	osrs_t, osrs_p, osrs_h Oversample
	t_sb                   StandbyTime
	filter                 Filter
	spi3w_en               bool
	coefficients           cal
	ch                     chan<- Event
	d16                    [2]byte
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DEFAULT_SPI_SPEED = 4 * 1e6   // 4Mhz
	DEFAULT_I2C_SPEED = 100 * 1e3 // 100 kHz
	DEFAULT_I2C_SLAVE = 0x77
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	BME280_SOFTRESET_VALUE    = 0xB6
	BME280_SKIPTEMP_VALUE     = 0x80000
	BME280_SKIPPRESSURE_VALUE = 0x80000
	BME280_SKIPHUMID_VALUE    = 0x8000
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create new BME280 device, with channel. Panics on error
func New(cfg Config, ch chan<- Event) *device {
	if d, err := cfg.New(ch); err != nil {
		panic(err)
	} else {
		return d
	}
}

// Create new BME280 device, with channel
func (cfg Config) New(ch chan<- Event) (*device, error) {
	d := new(device)

	// Set communciation device
	switch {
	case cfg.I2C != nil:
		d.i2c = cfg.I2C
		if cfg.Slave == 0 {
			d.slave = DEFAULT_I2C_SLAVE
		} else {
			d.slave = cfg.Slave
		}
	case cfg.SPI != nil:
		d.spi = cfg.SPI
	default:
		return nil, ErrBadParameter.With("bme280")
	}

	// Set channel
	if ch == nil {
		return nil, ErrBadParameter.With("bme280")
	} else {
		d.ch = ch
	}

	// Sync registers
	if err := d.sync(); err != nil {
		return nil, err
	}

	// Return success
	return d, nil
}

func (d *device) Close() error {
	var result error

	// Put device into sleep state
	if err := d.SetMode(BME280_MODE_SLEEP); err != nil {
		result = multierror.Append(result, err)
	}

	// Close underlying devices
	if d.i2c != nil {
		if err := d.i2c.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if d.spi != nil {
		if err := d.spi.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	d.slave = 0
	d.ch = nil
	d.i2c = nil
	d.spi = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d *device) String() string {
	str := "<bme280"
	if d.slave != 0 {
		str += fmt.Sprintf(" slave=0x%02X", d.slave)
	}
	str += fmt.Sprintf(" chipid=0x%02X", d.chipid)
	str += fmt.Sprintf(" version=0x%02X", d.version)
	str += fmt.Sprint(" mode=", d.mode)
	str += fmt.Sprint(" osrs_t=", d.osrs_t)
	str += fmt.Sprint(" osrs_p=", d.osrs_p)
	str += fmt.Sprint(" osrs_h=", d.osrs_h)
	str += fmt.Sprint(" t_sb=", d.t_sb)
	str += fmt.Sprint(" filter=", d.filter)
	str += fmt.Sprint(" spi3w_en=", d.spi3w_en)
	str += fmt.Sprint(" coefficients=", d.coefficients)

	if d.i2c != nil {
		str += fmt.Sprint(" i2c=", d.i2c)
	}
	if d.spi != nil {
		str += fmt.Sprint(" spi=", d.spi)
	}

	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Sample sensor and emit an event on channel which includes the temperature,
// humidity and pressure. ErrSampleSkipped is returned if no sample was taken,
// ErrTimeout if either the device timed out or if the channel was blocked
func (d *device) Sample() error {
	// Wait for no measuring or updating
	if err := d.wait(); err != nil {
		return err
	}

	// Set mode of operation if we're in FORCED or SLEEP mode, and wait until we
	// can read the measurement for the correct amount of time
	mode := d.mode
	if mode == BME280_MODE_FORCED || mode == BME280_MODE_SLEEP {
		if err := d.SetMode(BME280_MODE_FORCED); err != nil {
			return err
		}
		// Wait until we can measure
		time.Sleep(toMeasurementTime(d.osrs_t, d.osrs_p, d.osrs_h))
	}

	// Read samples
	data, err := d.Read()
	if err != nil {
		return err
	}

	// Convert to calibrated values
	event := event.New(d)
	tvalue, tfine, err := toTemperature(data, d.coefficients)
	if errors.Is(err, ErrSampleSkipped) {
		return ErrSampleSkipped
	} else {
		event.SetTemperature(tvalue)
	}
	pvalue, err := toPressure(data, tfine, d.coefficients)
	if err == nil {
		event.SetPressure(pvalue)
		avalue, err := toAltitude(tvalue, pvalue, 101325*1000)
		if err == nil {
			event.SetAltitude(avalue)
		}
	}
	hvalue, err := toHumidity(data, tfine, d.coefficients)
	if err == nil {
		event.SetHumidity(hvalue)
	}

	// Emit the event, and return any errors
	return event.Emit(d.ch)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (d *device) sync() error {
	// Read ChipId and Version
	if chipid, err := d.ChipID(); err != nil {
		return err
	} else if version, err := d.Version(); err != nil {
		return err
	} else {
		d.chipid = chipid
		d.version = version
	}

	// Read control values
	if mode, osrs_t, osrs_p, osrs_h, err := d.Control(); err != nil {
		return err
	} else {
		d.mode = mode
		d.osrs_t = osrs_t
		d.osrs_h = osrs_h
		d.osrs_p = osrs_p
	}

	// Read config values
	if t_sb, filter, spi3w_en, err := d.Config(); err != nil {
		return err
	} else {
		d.t_sb = t_sb
		d.filter = filter
		d.spi3w_en = spi3w_en
	}

	// Read calibration coefficients
	if coefficients, err := d.calibrate(); err != nil {
		return err
	} else {
		d.coefficients = coefficients
	}

	// Return success
	return nil
}
