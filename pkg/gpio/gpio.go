package gpio

import (
	"machine"

	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Pin uint8

type Config struct {
	In  []Pin
	Out []Pin
}

type device struct {
	in  map[Pin]machine.Pin
	out map[Pin]machine.Pin
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) New() (*device, error) {
	this := new(device)
	this.in = make(map[Pin]machine.Pin, len(cfg.In))
	this.out = make(map[Pin]machine.Pin, len(cfg.Out))

	// Make map for input pins
	for _, pin := range cfg.In {
		_pin := machine.Pin(pin)
		if _, exists := this.in[pin]; exists {
			return nil, ErrBadParameter
		} else {
			this.in[pin] = _pin
		}
	}

	// Make map for output pins
	for _, pin := range cfg.Out {
		_pin := machine.Pin(pin)
		if _, exists := this.in[pin]; exists {
			return nil, ErrBadParameter
		} else if _, exists := this.out[pin]; exists {
			return nil, ErrBadParameter
		} else {
			this.out[pin] = _pin
		}
	}

	// Configure pins
	for _, pin := range this.out {
		pin.Configure(machine.PinConfig{
			Mode: machine.PinOutput,
		})
	}
	for _, pin := range this.in {
		pin.Configure(machine.PinConfig{
			Mode: machine.PinInput,
		})
	}

	// Return success
	return this, nil
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
		if pin_, exists := d.out[pin]; exists {
			pin_.High()
		}
	}
}

// Set output pins to low
func (d *device) Low(pins ...Pin) {
	for _, pin := range pins {
		if pin_, exists := d.out[pin]; exists {
			pin_.Low()
		}
	}
}
