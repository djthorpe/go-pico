package bme280

import (
	"errors"
	"fmt"
	"time"

	// Package imports
	pico "github.com/djthorpe/go-pico"
	multierror "github.com/hashicorp/go-multierror"

	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type I2CConfig struct {
	Bus   uint   // I2C Bus (0 or 1)
	Slave uint8  // BME280 Slave address, uses DEFAULT_I2C_SLAVE if not set
	Speed uint32 // I2C Communication Speed in Hz, uses DEFAULT_I2C_SPEED if not set
}

type SPIConfig struct {
	Bus   uint   // SPI Bus (0, 1 or 2)
	Slave uint   // SPI Slave (0 or 1) - not used on Pico
	Speed uint32 // SPI Communication Speed in Hz, uses DEFAULT_SPI_SPEED if not set
}

type device struct {
	i2c                    pico.I2C
	spi                    pico.SPI
	slave                  uint8
	chipid, version        uint8
	mode                   Mode
	osrs_t, osrs_p, osrs_h Oversample
	t_sb                   StandbyTime
	filter                 Filter
	spi3w_en               bool
	coefficients           cal
	ch                     chan Event
	d16                    [2]byte
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DEFAULT_I2C_SPEED = 100 * 1000 // 100 kHz
	DEFAULT_I2C_SLAVE = 0x77
	DEFAULT_SPI_SPEED = 4 * 1000 * 1000 // 4Mhz
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

// Create BME280 device on the I2C bus
func (cfg I2CConfig) New() (*device, error) {
	this := new(device)

	// Create I2C device
	if device, err := NewI2C(cfg); err != nil {
		return nil, err
	} else {
		this.i2c = device
		this.slave = cfg.Slave | DEFAULT_I2C_SLAVE
	}

	if err := this.sync(); err != nil {
		return nil, err
	}

	// Set channel
	this.ch = make(chan Event)

	// Return success
	return this, nil
}

// Create BME280 device on the SPI bus
func (cfg SPIConfig) New() (*device, error) {
	this := new(device)

	// Create SPI device
	if device, err := NewSPI(cfg); err != nil {
		return nil, err
	} else {
		this.spi = device
	}

	if err := this.sync(); err != nil {
		return nil, err
	}

	// Set channel
	this.ch = make(chan Event)

	// Return success
	return this, nil
}

func (d *device) Close() error {
	var result error

	// Close the channel
	if d.ch != nil {
		close(d.ch)
		d.ch = nil
	}

	// Put device into sleep state
	if err := d.SetMode(BME280_MODE_SLEEP); err != nil {
		result = multierror.Append(result, err)
	}

	// Close underlying devices
	if d.i2c != nil {
		if err := d.i2c.Close(); err != nil {
			result = multierror.Append(result, err)
		}
		d.i2c = nil
	}
	if d.spi != nil {
		if err := d.spi.Close(); err != nil {
			result = multierror.Append(result, err)
		}
		d.spi = nil
	}

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

// Return event channel
func (d *device) C() <-chan Event {
	return d.ch
}

// Sample sensor and emit an event on channel which includes the temperature,
// humidity and pressure. ErrSampleSkipped is returned if no sample was taken,
// ErrTimeout if either the device timed out or if the channel was blocked
func (d *device) Sample() error {
	var event Event

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
	tvalue, tfine, err := toTemperature(data, d.coefficients)
	if errors.Is(err, ErrSampleSkipped) {
		return ErrSampleSkipped
	} else {
		event.Type |= Temperature
		event.Temperature = tvalue
	}
	pvalue, err := toPressure(data, tfine, d.coefficients)
	if err == nil {
		event.Type |= Pressure
		event.Pressure = pvalue
		avalue, err := toAltitude(tvalue, pvalue, 101325*1000)
		if err == nil {
			event.Type |= Altitude
			event.Altitude = avalue
		}
	}
	hvalue, err := toHumidity(data, tfine, d.coefficients)
	if err == nil {
		event.Type |= Humidity
		event.Humidity = hvalue
	}

	// Emit the event, or return ErrTimeout
	select {
	case d.ch <- event:
		return nil
	default:
		return ErrTimeout
	}
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
