//go:build tinygo

package gpio

import (
	"fmt"
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
	pwm map[Pin]pwmChannel
	ch  chan<- Event
}

type pwmChannel interface {
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) New(ch chan<- Event) (*device, error) {
	this := new(device)
	this.in = make(map[Pin]machine.Pin, len(cfg.In))
	this.out = make(map[Pin]machine.Pin, len(cfg.Out))
	this.pwm = make(map[Pin]pwmChannel, len(cfg.PWM))

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

	// Make map for PWM pins
	for _, pin := range cfg.PWM {
		if pin_, exists := this.out[pin]; !exists {
			return nil, ErrBadParameter.With("pwm:", pin)
		} else if ch, err := getPwmChannel(pin_); err != nil {
			return nil, err
		} else {
			this.pwm[pin] = ch
		}
	}

	// Configure pins
	for _, pin := range this.out {
		_pin := Pin(pin)
		if freq, exists := this.pwm[_pin]; exists {
			pin.Configure(machine.PinConfig{
				Mode: machine.PinPWM,
			})
			if err := this.SetFrequency(_pin, freq); err != nil {
				return nil, ErrBadParameter.With("pwm:", pin, err)
			}
		} else {
			pin.Configure(machine.PinConfig{
				Mode: machine.PinOutput,
			})
		}
	}
	for _, pin := range this.in {
		pin.Configure(machine.PinConfig{
			Mode: machine.PinInput,
		})
	}
	for _, pin := range cfg.Watch {
		pin_ := this.in[pin]
		if err := pin_.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {
			fmt.Println("evt:", Pin(pin_), p.Get())
			event.New(Pin(pin_)).Set(Sample, UnitNone, pin_.Get()).Emit(this.ch)
		}); err != nil {
			fmt.Println("err:", pin, err)
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

// Get output pin value
func (d *device) Get(pin Pin) bool {
	if pin_, exists := d.in[pin]; !exists {
		return false
	} else {
		return pin_.Get()
	}
}

// Set output pin to value
func (d *device) Set(pin Pin, v bool) {
	if pin_, exists := d.out[pin]; exists {
		if v {
			pin_.High()
		} else {
			pin_.Low()
		}
	}
}

// Set PWM square wave frequency
func (d *device) SetFrequency(pin Pin, f uint32) error {
	if _, exists := d.pwm[pin]; !exists {
		return ErrBadParameter
	} else {
		d.pwm[pin] = f
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func getPwmChannel(p machine.Pin) (pwmChannel, error) {
	slice, err := machine.PWMPeripheral(p)
	if err != nil {
		return nil, ErrBadParameter.With("pwm:", p, err)
	}
	switch slice {
	case 0:
		return machine.PWM0.Channel(p)
	case 1:
		return machine.PWM1.Channel(p)
	case 2:
		return machine.PWM2.Channel(p)
	case 3:
		return machine.PWM3.Channel(p)
	case 4:
		return machine.PWM4.Channel(p)
	case 5:
		return machine.PWM5.Channel(p)
	case 6:
		return machine.PWM6.Channel(p)
	case 7:
		return machine.PWM7.Channel(p)
	}
	return nil, ErrBadParameter.With("pwm:", p)
}
