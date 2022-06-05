//go:build rpi

package gpio

import (
	"os"
	"reflect"
	"syscall"
	"time"
	"unsafe"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
	. "github.com/djthorpe/go-pico/pkg/errors"

	// Package imports
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: bcm_host
#include "bcm_host.h"
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type device struct {
	mem8  []uint8  // access GPIO as bytes
	mem32 []uint32 // access GPIO as uint32
	pins  map[Pin]Mode
}

type Mode uint
type Pull uint

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	GPIO_DEV            = "/dev/gpiomem"
	GPIO_BASE    uint32 = 0x200000
	GPIO_MAXPINS        = 54 // GPIO0 to GPIO53
)

const (
	// GPIO Registers
	GPIO_GPLVL0    = 0x0034 // Register to read pins GPIO0-GPIO31
	GPIO_GPLVL1    = 0x0038 // Register to read pins GPIO32-GPIO53
	GPIO_GPSET0    = 0x001C // Register to write HIGH to pins GPIO0-GPIO31
	GPIO_GPSET1    = 0x0020 // Register to write HIGH to pins GPIO32-GPIO53
	GPIO_GPCLR0    = 0x0028 // Register to write LOW to pins GPIO0-GPIO31
	GPIO_GPCLR1    = 0x002C // Register to write LOW to pins GPIO32-GPIO53
	GPIO_GPFSEL0   = 0x0000 // Pin modes for GPIO0-GPIO9
	GPIO_GPFSEL1   = 0x0004 // Pin modes for GPIO10-GPIO19
	GPIO_GPFSEL2   = 0x0008 // Pin modes for GPIO20-GPIO29
	GPIO_GPFSEL3   = 0x000C // Pin modes for GPIO30-GPIO39
	GPIO_GPFSEL4   = 0x0010 // Pin modes for GPIO40-GPIO49
	GPIO_GPFSEL5   = 0x0014 // Pin modes for GPIO50-GPIO53
	GPIO_GPPUD     = 0x0094 // GPIO Pin Pull-up/down Enable
	GPIO_GPPUDCLK0 = 0x0098 // GPIO Pin Pull-up/down Enable Clock 0
	GPIO_GPPUDCLK1 = 0x009c // GPIO Pin Pull-up/down Enable Clock 1
)

const (
	INPUT Mode = iota
	OUTPUT
	ALT5
	ALT4
	ALT0
	ALT1
	ALT2
	ALT3
	NONE
)

const (
	PULL_OFF Pull = iota
	PULL_DOWN
	PULL_UP
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) New() (*device, error) {
	this := new(device)
	this.pins = make(map[Pin]Mode, len(cfg.In)+len(cfg.Out))

	// Open GPIO device for read/write
	f, err := os.OpenFile(GPIO_DEV, os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Memory map GPIO registers to byte array
	offset := int64(bcm_peripheral_address() + GPIO_BASE)
	size := int(bcm_peripheral_size())
	if mem8, err := syscall.Mmap(int(f.Fd()), offset, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED); err != nil {
		return nil, err
	} else {
		this.mem8 = mem8
	}

	// Convert mapped byte memory to unsafe []uint32 pointer, adjust length as needed
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&this.mem8))
	header.Len /= (32 / 8)
	header.Cap /= (32 / 8)
	this.mem32 = *(*[]uint32)(unsafe.Pointer(&header))

	// Check length of arrays
	if len(this.mem8) == 0 || len(this.mem32) == 0 {
		return nil, ErrUnexpectedValue
	}

	// Set output pins
	var result error
	for _, pin := range cfg.Out {
		if pin > GPIO_MAXPINS {
			result = multierror.Append(result, ErrBadParameter.With("gpio pin", pin))
		} else if _, exists := this.pins[pin]; exists {
			result = multierror.Append(result, ErrDuplicateValue.With("gpio pin", pin))
		} else {
			this.bcm_pin_setmode(pin, OUTPUT)
			this.pins[pin] = OUTPUT
		}
	}

	// Set input pins
	for _, pin := range cfg.In {
		if pin > GPIO_MAXPINS {
			result = multierror.Append(result, ErrBadParameter.With("gpio pin", pin))
		} else if _, exists := this.pins[pin]; exists {
			result = multierror.Append(result, ErrDuplicateValue.With("gpio pin", pin))
		} else {
			this.bcm_pin_setmode(pin, INPUT)
			this.pins[pin] = INPUT
		}
	}

	// Return any errors
	return this, result
}

func (d *device) Close() error {
	var result error

	if d.mem8 != nil {
		if err := syscall.Munmap(d.mem8); err != nil {
			result = multierror.Append(result, os.NewSyscallError("munmap", err))
		}
	}

	// Release resources
	d.mem8 = nil
	d.mem32 = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d *device) String() string {
	str := "<gpio"
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Set output pins to high
func (d *device) High(pins ...Pin) {
	for _, pin := range pins {
		d.bcm_pin_write(pin, true)
	}
}

// Set output pins to low
func (d *device) Low(pins ...Pin) {
	for _, pin := range pins {
		d.bcm_pin_write(pin, false)
	}
}

func (d *device) Get(Pin) bool {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (d *device) bcm_pin_write(p Pin, v bool) {
	// Silently ignore invalid pins
	if p >= GPIO_MAXPINS {
		return
	}

	// Shift in value
	value := uint32(1 << (uint8(p) & 31))
	switch v {
	case false:
		if uint8(p) <= 31 {
			d.mem32[GPIO_GPCLR0>>2] = value
		} else {
			d.mem32[GPIO_GPCLR1>>2] = value
		}
	case true:
		if uint8(p) <= 31 {
			d.mem32[GPIO_GPSET0>>2] = value
		} else {
			d.mem32[GPIO_GPSET1>>2] = value
		}
	}
}

func (d *device) bcm_pin_setmode(p Pin, v Mode) {
	// Silently ignore invalid pins
	if p >= GPIO_MAXPINS {
		return
	}

	// get register and the number of bits to shift to access the current mode
	register, shift := bcm_pin_to_register(p)

	// set register
	d.mem32[register>>2] = (d.mem32[register>>2] &^ (7 << shift)) | (uint32(v) << shift)
}

func (d *device) bcm_pin_setpull(p Pin, v Pull) {
	// Silently ignore invalid pins
	if p >= GPIO_MAXPINS {
		return
	}

	// Set the low two bits of register to 0 (off) 1 (down) or 2 (up)
	switch v {
	case PULL_UP, PULL_DOWN:
		d.mem32[GPIO_GPPUD] |= uint32(v)
	case PULL_OFF:
		d.mem32[GPIO_GPPUD] &^= 3
	}

	// Wait for 150 cycles
	time.Sleep(time.Microsecond)

	// Determine clock register
	clockReg := GPIO_GPPUDCLK0
	if p >= Pin(32) {
		clockReg = GPIO_GPPUDCLK1
	}

	// Clock it in
	d.mem32[clockReg] = 1 << (p % 32)

	// Wait for value to clock in
	time.Sleep(time.Microsecond)

	// Write 00 to the register to clear it
	d.mem32[GPIO_GPPUD] &^= 3

	// Wait for value to clock in
	time.Sleep(time.Microsecond)

	// Remove the clock
	d.mem32[clockReg] = 0
}

func (d *device) bcm_pin_getmode(p Pin) Mode {
	// Silently ignore invalid pins
	if p >= GPIO_MAXPINS {
		return NONE
	}

	// return the register and the number of bits to shift to
	// access the current mode
	register, shift := bcm_pin_to_register(p)

	// Retrieve register, shift to the right, and return last three bits
	return Mode((d.mem32[register>>2] >> shift) & 7)
}

func bcm_peripheral_address() uint32 {
	return uint32(C.bcm_host_get_peripheral_address())
}

func bcm_peripheral_size() uint32 {
	return uint32(C.bcm_host_get_peripheral_size())
}

func bcm_pin_to_register(pin Pin) (uint, uint) {
	p := uint(pin)
	switch {
	case p >= 0 && p <= 9:
		return GPIO_GPFSEL0, uint(p * 3)
	case p >= 10 && p <= 19:
		return GPIO_GPFSEL1, uint((p - 10) * 3)
	case p >= 20 && p <= 29:
		return GPIO_GPFSEL2, uint((p - 20) * 3)
	case p >= 30 && p <= 39:
		return GPIO_GPFSEL3, uint((p - 30) * 3)
	case p >= 40 && p <= 49:
		return GPIO_GPFSEL4, uint((p - 40) * 3)
	default:
		return GPIO_GPFSEL5, uint((p - 50) * 3)
	}
}
