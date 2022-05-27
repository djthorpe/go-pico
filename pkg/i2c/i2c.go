package i2c

import (
	"fmt"
	"machine"

	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Bus       uint
	Frequency uint32
}

type device struct {
	*machine.I2C
	sda, scl machine.Pin
	d8, d16  []byte
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) New() (*device, error) {
	this := new(device)

	// Pico has two I2C devices
	switch cfg.Bus {
	case 0:
		this.I2C = machine.I2C0
		this.sda = machine.I2C0_SDA_PIN
		this.scl = machine.I2C0_SCL_PIN
	case 1:
		this.I2C = machine.I2C1
		this.sda = machine.I2C1_SDA_PIN
		this.scl = machine.I2C1_SCL_PIN
	default:
		return nil, ErrBadParameter
	}

	// Initialise the device
	if err := this.I2C.Configure(machine.I2CConfig{
		Frequency: cfg.Frequency,
		SDA:       this.sda,
		SCL:       this.scl,
	}); err != nil {
		return nil, err
	}

	// Create buffers for data
	this.d8 = make([]uint8, 1)
	this.d16 = make([]uint8, 2)

	// Return success
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d *device) String() string {
	str := "<i2c"
	str += fmt.Sprint(" sda=", d.sda)
	str += fmt.Sprint(" scl=", d.scl)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// ReadRegister reads a register
func (d *device) ReadRegister(slave, register uint8, data []uint8) error {
	return d.I2C.ReadRegister(slave, register, data)
}

// ReadRegister_Uint8 reads a byte register
func (d *device) ReadRegister_Uint8(slave, register uint8) (uint8, error) {
	if err := d.ReadRegister(slave, register, d.d8); err != nil {
		return 0, err
	} else {
		return d.d8[0], nil
	}
}

// ReadRegister_Uint16 reads a word register
func (d *device) ReadRegister_Uint16(slave, register uint8) (uint16, error) {
	if err := d.ReadRegister(slave, register, d.d16); err != nil {
		return 0, err
	} else {
		return uint16(d.d16[0])<<8 | uint16(d.d16[1]), nil
	}
}

// WriteRegister_Uint8 writes a byte to a register
func (d *device) WriteRegister_Uint8(slave, register, data uint8) error {
	d.d8[0] = data
	return d.I2C.WriteRegister(slave, register, d.d8)
}
