//go:build rpi

package i2c

import (
	"fmt"
	"os"

	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/errors"
	. "github.com/djthorpe/go-pico/pkg/linux"

	// Packages
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Bus uint
}

type device struct {
	bus   uint
	f     *os.File
	fd    uintptr
	funcs I2CFunction
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) New() (*device, error) {
	this := new(device)

	// Open device
	if f, err := I2COpenDevice(cfg.Bus); err != nil {
		return nil, err
	} else {
		this.f = f
		this.fd = uintptr(f.Fd())
		this.bus = cfg.Bus
	}

	// Get features for this device
	if funcs, err := I2CFunctions(this.fd); err != nil {
		this.Close()
		return nil, err
	} else {
		this.funcs = funcs
	}

	// Return success
	return this, nil
}

func (d *device) Close() error {
	var result error

	if d.f != nil {
		if err := d.f.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	d.f = nil
	d.fd = 0
	d.funcs = 0

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d *device) String() string {
	str := "<i2c"
	if d.f != nil {
		str += fmt.Sprintf(" bus=%q", I2CDevice(d.bus))
	}
	if d.funcs != 0 {
		str += fmt.Sprint(" funcs=", d.funcs)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// ReadRegister reads a register into a block. Block size should be no
// more than about 30 bytes
func (d *device) ReadRegister(slave, register uint8, data []uint8) error {
	if len(data) > 0xFF {
		return ErrBadParameter.With("ReadRegister")
	}
	if err := I2CSetSlave(d.fd, slave); err != nil {
		return err
	}
	if r, err := I2CReadBlock(d.fd, register, uint8(len(data)), d.funcs); err != nil {
		return err
	} else {
		copy(data, r)
	}

	// Return success
	return nil
}

// ReadRegister_Uint8 reads a byte register
func (d *device) ReadRegister_Uint8(slave, register uint8) (uint8, error) {
	if err := I2CSetSlave(d.fd, slave); err != nil {
		return 0, err
	}
	return I2CReadUint8(d.fd, register, d.funcs)
}

// ReadRegister_Uint16 reads a word register
func (d *device) ReadRegister_Uint16(slave, register uint8) (uint16, error) {
	if err := I2CSetSlave(d.fd, slave); err != nil {
		return 0, err
	}
	return I2CReadUint16(d.fd, register, d.funcs)
}

// WriteRegister_Uint8 writes a byte to a register
func (d *device) WriteRegister_Uint8(slave, register, data uint8) error {
	if err := I2CSetSlave(d.fd, slave); err != nil {
		return err
	}
	return I2CWriteUint8(d.fd, register, data, d.funcs)
}
