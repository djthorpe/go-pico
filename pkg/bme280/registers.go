package bme280

import (
	// Namespace imports
	"time"

	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type register uint8

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	BME280_REG_DIG_T1       register = 0x88
	BME280_REG_DIG_T2       register = 0x8A
	BME280_REG_DIG_T3       register = 0x8C
	BME280_REG_DIG_P1       register = 0x8E
	BME280_REG_DIG_P2       register = 0x90
	BME280_REG_DIG_P3       register = 0x92
	BME280_REG_DIG_P4       register = 0x94
	BME280_REG_DIG_P5       register = 0x96
	BME280_REG_DIG_P6       register = 0x98
	BME280_REG_DIG_P7       register = 0x9A
	BME280_REG_DIG_P8       register = 0x9C
	BME280_REG_DIG_P9       register = 0x9E
	BME280_REG_DIG_H1       register = 0xA1
	BME280_REG_DIG_H2       register = 0xE1
	BME280_REG_DIG_H3       register = 0xE3
	BME280_REG_DIG_H4       register = 0xE4
	BME280_REG_DIG_H5       register = 0xE5
	BME280_REG_DIG_H6       register = 0xE7
	BME280_REG_CHIPID       register = 0xD0
	BME280_REG_VERSION      register = 0xD1
	BME280_REG_SOFTRESET    register = 0xE0
	BME280_REG_CAL26        register = 0xE1 // R calibration stored in 0xE1-0xF0
	BME280_REG_CONTROLHUMID register = 0xF2
	BME280_REG_STATUS       register = 0xF3
	BME280_REG_CONTROL      register = 0xF4
	BME280_REG_CONFIG       register = 0xF5
	BME280_REG_PRESSUREDATA register = 0xF7
	BME280_REG_TEMPDATA     register = 0xFA
	BME280_REG_HUMIDDATA    register = 0xFD
)

const (
	// Write mask
	BME280_REG_SPI_WRITE register = 0x7F

	// Timeout for measuring
	BME280_TIMEOUT = time.Millisecond * 500
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (d *device) ChipID() (uint8, error) {
	return d.i2c.ReadRegister_Uint8(d.slave, uint8(BME280_REG_CHIPID))
}

func (d *device) Version() (uint8, error) {
	return d.i2c.ReadRegister_Uint8(d.slave, uint8(BME280_REG_VERSION))
}

// Read values mode, osrs_t, osrs_p, osrs_h
func (d *device) Control() (Mode, Oversample, Oversample, Oversample, error) {
	// Read control values
	ctrl_meas, err := d.i2c.ReadRegister_Uint8(d.slave, uint8(BME280_REG_CONTROL))
	if err != nil {
		return 0, 0, 0, 0, err
	}
	ctrl_hum, err := d.i2c.ReadRegister_Uint8(d.slave, uint8(BME280_REG_CONTROLHUMID))
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Convert values
	mode := Mode(ctrl_meas) & BME280_MODE_MAX
	osrs_t := Oversample(ctrl_meas>>5) & BME280_OVERSAMPLE_MAX
	osrs_p := Oversample(ctrl_meas>>2) & BME280_OVERSAMPLE_MAX
	osrs_h := Oversample(ctrl_hum) & BME280_OVERSAMPLE_MAX

	// Return success
	return mode, osrs_t, osrs_p, osrs_h, nil
}

// Read values t_sb, filter, spi3w_en
func (d *device) Config() (StandbyTime, Filter, bool, error) {
	config, err := d.i2c.ReadRegister_Uint8(d.slave, uint8(BME280_REG_CONFIG))
	if err != nil {
		return 0, 0, false, err
	}

	// Convert values
	filter := Filter(config>>2) & BME280_FILTER_MAX
	t_sb := StandbyTime(config>>5) & BME280_STANDBY_MAX
	spi3w_en := bool(config&0x01 != 0x00)

	// Return success
	return t_sb, filter, spi3w_en, nil
}

// Return current measuring and updating value
func (d *device) Status() (bool, bool, error) {
	status, err := d.i2c.ReadRegister_Uint8(d.slave, uint8(BME280_REG_STATUS))
	if err != nil {
		return false, false, err
	}

	// Convert values
	measuring := ((status>>3)&0x01 != 0x00)
	updating := (status&0x01 != 0x00)

	// Return success
	return measuring, updating, nil
}

// Reset the device using the complete power-on-reset procedure
func (d *device) Reset() error {
	if err := d.i2c.WriteRegister_Uint8(d.slave, uint8(BME280_REG_SOFTRESET), BME280_SOFTRESET_VALUE); err != nil {
		return err
	}
	if err := d.wait(); err != nil {
		return err
	}

	// Read registers and return
	return d.sync()
}

// Read raw sample values
func (d *device) Read() ([8]byte, error) {
	var data [8]byte
	if err := d.i2c.ReadRegister(d.slave, uint8(BME280_REG_PRESSUREDATA), data[:]); err != nil {
		return data, err
	} else {
		return data, nil
	}
}

// Set mode
func (d *device) SetMode(mode Mode) error {
	ctrl_meas := uint8(d.osrs_t)<<5 | uint8(d.osrs_p)<<2 | uint8(mode)
	if err := d.i2c.WriteRegister_Uint8(d.slave, uint8(BME280_REG_CONTROL), ctrl_meas); err != nil {
		return err
	}
	if mode_, _, _, _, err := d.Control(); err != nil {
		return err
	} else if mode != mode_ {
		return ErrUnexpectedValue
	} else {
		d.mode = mode_
	}

	// Return success
	return nil
}

func (d *device) SetTempOversample(osrs_t Oversample) error {
	osrs_p := d.osrs_p
	ctrl_meas := uint8(osrs_t&BME280_OVERSAMPLE_MAX)<<5 | uint8(osrs_p&BME280_OVERSAMPLE_MAX)<<2 | uint8(d.mode&BME280_MODE_MAX)
	if err := d.i2c.WriteRegister_Uint8(d.slave, uint8(BME280_REG_CONTROL), ctrl_meas); err != nil {
		return err
	} else if err := d.wait(); err != nil {
		return err
	} else if _, osrs_t_, osrs_p_, osrs_h_, err := d.Control(); err != nil {
		return err
	} else if osrs_t_ != osrs_t {
		return ErrUnexpectedValue
	} else {
		d.osrs_t = osrs_t_
		d.osrs_p = osrs_p_
		d.osrs_h = osrs_h_
		return nil
	}
}

func (d *device) SetPressureOversample(osrs_p Oversample) error {
	osrs_t := d.osrs_t
	ctrl_meas := uint8(osrs_t&BME280_OVERSAMPLE_MAX)<<5 | uint8(osrs_p&BME280_OVERSAMPLE_MAX)<<2 | uint8(d.mode&BME280_MODE_MAX)
	if err := d.i2c.WriteRegister_Uint8(d.slave, uint8(BME280_REG_CONTROL), ctrl_meas); err != nil {
		return err
	} else if err := d.wait(); err != nil {
		return err
	} else if _, osrs_t_, osrs_p_, osrs_h_, err := d.Control(); err != nil {
		return err
	} else if osrs_p_ != osrs_p {
		return ErrUnexpectedValue
	} else {
		d.osrs_t = osrs_t_
		d.osrs_p = osrs_p_
		d.osrs_h = osrs_h_
		return nil
	}
}

func (d *device) SetHumidityOversample(osrs_h Oversample) error {
	if err := d.i2c.WriteRegister_Uint8(d.slave, uint8(BME280_REG_CONTROLHUMID), uint8(osrs_h&BME280_OVERSAMPLE_MAX)); err != nil {
		return err
	} else if err := d.wait(); err != nil {
		return err
	} else if _, osrs_t_, osrs_p_, osrs_h_, err := d.Control(); err != nil {
		return err
	} else if osrs_h_ != osrs_h {
		return ErrUnexpectedValue
	} else {
		d.osrs_t = osrs_t_
		d.osrs_p = osrs_p_
		d.osrs_h = osrs_h_
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (d *device) wait() error {
	// Wait for no measuring or updating
	// TODO: Timeout
	for {
		select {
		default:
			if measuring, updating, err := d.Status(); err != nil {
				return err
			} else if measuring == false && updating == false {
				return nil
			}
			time.Sleep(time.Millisecond)
		}
	}
}
