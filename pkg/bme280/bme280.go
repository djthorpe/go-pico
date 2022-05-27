package bme280

import (
	"fmt"
	"time"

	// Package imports
	pico "github.com/djthorpe/go-pico"
	i2c "github.com/djthorpe/go-pico/pkg/i2c"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Bus   uint   // I2C Bus (0 or 1)
	Speed uint32 // I2C Communication Speed in Hz, optional
	Slave uint8  // BME280 Slave address, optional
}

type I2C interface{}

type device struct {
	i2c                    pico.I2C
	slave                  uint8
	chipid, version        uint8
	mode                   Mode
	osrs_t, osrs_p, osrs_h Oversample
	t_sb                   StandbyTime
	filter                 Filter
	spi3w_en               bool
}

// calibration data
type cal struct {
	T1                             uint16
	T2, T3                         int16
	P1                             uint16
	P2, P3, P4, P5, P6, P7, P8, P9 int16
	H1                             uint8
	H2                             int16
	H3                             uint8
	H4, H5                         int16
	H6                             int8
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DEFAULT_SPEED = 100 * 1000 // 100 kHz
	DEFAULT_SLAVE = 0x77
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	BME280_SOFTRESET_VALUE    = 0xB6
	BME280_SKIPTEMP_VALUE     = 0x80000
	BME280_SKIPPRESSURE_VALUE = 0x80000
	BME280_SKIPHUMID_VALUE    = 0x8000
	BME280_CALIBRATION_SIZE   = 33
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) New() (*device, error) {
	this := new(device)

	// Create I2C device
	device, err := i2c.Config{
		Bus:       cfg.Bus,
		Frequency: cfg.Speed | DEFAULT_SPEED,
	}.New()
	if err != nil {
		return nil, err
	} else {
		this.i2c = device
		this.slave = cfg.Slave | DEFAULT_SLAVE
	}

	if err := this.sync(); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d *device) String() string {
	str := "<bme280"
	str += fmt.Sprintf(" slave=0x%02X", d.slave)
	str += fmt.Sprintf(" chipid=0x%02X", d.chipid)
	str += fmt.Sprintf(" version=0x%02X", d.version)
	str += fmt.Sprint(" mode=", d.mode)
	str += fmt.Sprint(" osrs_t=", d.osrs_t)
	str += fmt.Sprint(" osrs_p=", d.osrs_p)
	str += fmt.Sprint(" osrs_h=", d.osrs_h)
	str += fmt.Sprint(" t_sb=", d.t_sb)
	str += fmt.Sprint(" filter=", d.filter)
	str += fmt.Sprint(" spi3w_en=", d.spi3w_en)
	return str + ">"
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

	// Read calibration data
	t1, err := d.i2c.ReadRegister_Uint16(d.slave, uint8(BME280_REG_DIG_T1))
	if err != nil {
		return err
	}
	fmt.Println("t1=", t1)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (d *device) Sample() error {
	// Save mode, if it's SLEEP then return to sleep afterwards
	mode := d.mode

	// Wait for no measuring or updating
	// TODO: Timeout
	for {
		if measuring, updating, err := d.Status(); err != nil {
			return err
		} else if measuring == false && updating == false {
			break
		}
	}

	// Set mode of operation if we're in FORCED or SLEEP mode, and wait until we
	// can read the measurement for the correct amount of time
	if mode == BME280_MODE_FORCED || mode == BME280_MODE_SLEEP {
		if err := d.SetMode(BME280_MODE_FORCED); err != nil {
			return err
		}
		// Wait until we can measure
		time.Sleep(toMeasurementTime(d.osrs_t, d.osrs_p, d.osrs_h))
	}

	// Read temperature, return error if temperature reading is skipped
	adc_t, err := d.Temperature()
	if err != nil {
		return err
	}
	adc_p, err := d.Pressure()
	if err != nil {
		return err
	}
	adc_h, err := d.Humidity()
	if err != nil {
		return err
	}

	fmt.Printf("temp=%04X pressure=%04X humidity=%04X\n", adc_t, adc_p, adc_h)

	// Return success
	return nil
}
