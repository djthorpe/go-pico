//go:build rpi

package spi

import (
	"fmt"
	"os"

	// Namespace imports
	//. "github.com/djthorpe/go-pico/pkg/errors"
	. "github.com/djthorpe/go-pico/pkg/linux"

	// Package imports
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Bus, Slave uint
	Speed      uint32
	Mode       Mode
}

type device struct {
	bus, slave uint
	f          *os.File
	fd         uintptr
	speed      uint32
	delay      uint16
	bits       uint8
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) New() (*device, error) {
	this := new(device)

	// Open device
	if f, err := SPIOpenDevice(cfg.Bus, cfg.Slave); err != nil {
		return nil, err
	} else {
		this.f = f
		this.fd = uintptr(f.Fd())
		this.bus = cfg.Bus
		this.slave = cfg.Slave
	}

	// Set speed
	if cfg.Speed != 0 {
		if err := SPISetSpeedHz(this.fd, cfg.Speed); err != nil {
			this.Close()
			return nil, err
		}
	}

	// Set mode
	if err := SPISetMode(this.fd, SPIMode(cfg.Mode)&SPI_MODE_MASK); err != nil {
		this.Close()
		return nil, err
	}

	// Always set 8 bits per word
	if err := SPISetBitsPerWord(this.fd, 8); err != nil {
		this.Close()
		return nil, err
	}

	// Read back speed, delay and bits values
	if speed, err := SPISpeedHz(this.fd); err != nil {
		this.Close()
		return nil, err
	} else if bits, err := SPIBitsPerWord(this.fd); err != nil {
		this.Close()
		return nil, err
	} else {
		this.speed = speed
		this.bits = bits
		this.delay = 0
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

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d *device) String() string {
	str := "<spi"
	if d.f != nil {
		str += fmt.Sprintf(" dev=%q", SPIDevice(d.bus, d.slave))
		str += fmt.Sprintf(" speed=%vHz", d.speed)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Transfer writes then reads from SPI bus
func (d *device) Transfer(w, r []byte) error {
	if data, err := SPITransfer(d.fd, w, d.speed, d.delay, d.bits); err != nil {
		return err
	} else {
		copy(r, data)
	}

	// Return success
	return nil
}
