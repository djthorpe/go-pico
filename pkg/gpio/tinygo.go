//go:build tinygo

package gpio

import (
	"machine"

	// Package imports
	event "github.com/djthorpe/go-pico/pkg/event"
	"github.com/hashicorp/go-multierror"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type device struct {
	in  map[Pin]machine.Pin
	out map[Pin]machine.Pin
	ch  chan<- Event
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) New(ch chan<- Event) (*device, error) {
	this := new(device)
	this.in = make(map[Pin]machine.Pin, len(cfg.In))
	this.out = make(map[Pin]machine.Pin, len(cfg.Out))

	// Check channel
	if ch == nil {
		return nil, ErrBadParameter.With("gpio")
	} else {
		this.ch = ch
	}

	// Make map for input pins
	for _, pin := range cfg.In {
		_pin := machine.Pin(pin)
		if _, exists := this.in[pin]; exists {
			return nil, ErrBadParameter
		} else {
			this.in[pin] = _pin
		}
	}

	// Check for watch pins being the same as in pins
	for _, pin := range cfg.Watch {
		if _, exists := this.in[pin]; !exists {
			return nil, ErrBadParameter.With("watch:", pin)
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
	for _, pin := range cfg.Watch {
		pin_ := this.in[pin]
		if err := pin_.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {
			event.New(Pin(pin_)).Set(Sample, UnitNone, pin_.Get()).Emit(this.ch)
		}); err != nil {
			//return nil, ErrBadParameter.With("gpio:", pin, err)
		}
	}

	// Return success
	return this, nil
}

func (d *device) Close() error {
	var result error

	// Unset interrupts
	for _, pin := range d.in {
		if err := pin.SetInterrupt(0, nil); err != nil {
			result = multierror.Append(result, err)
		}
	}

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

// Set output pins to low
func (d *device) Get(pin Pin) bool {
	if pin_, exists := d.in[pin]; !exists {
		return false
	} else {
		return pin_.Get()
	}
}
