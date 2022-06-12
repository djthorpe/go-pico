package pico

import (
	//	. "github.com/djthorpe/go-pico/pkg/errors"
	. "github.com/djthorpe/go-pico/pkg/sdk"
)

//////////////////////////////////////////////////////////////////////////////
// TYPES

type PWM struct {
	pin       Pin
	slice_num uint32
	ch        PWM_chan
	config    *PWM_config
}

//////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewPWM(pin Pin) *PWM {
	pwm := new(PWM)
	pwm.pin = pin
	pwm.slice_num = PWM_gpio_to_slice_num(GPIO_pin(pin))
	pwm.ch = PWM_gpio_to_channel(GPIO_pin(pin))
	pwm.config = PWM_get_default_config()
	return pwm
}

func (p *PWM) SetEnabled(enabled bool) {
	if enabled {
		PWM_init(p.slice_num, p.config, true)
	} else {
		PWM_set_enabled(p.slice_num, enabled)
	}
}

// Get counter value
//
func (p *PWM) Get() uint16 {
	return PWM_get_counter(p.slice_num)
}

// Set counter value
//
func (p *PWM) Set(value uint16) {
	PWM_set_counter(p.slice_num, value)
}

// Increment counter
//
func (p *PWM) Inc() {
	PWM_advance_count(p.slice_num)
}

// Decrement counter
//
func (p *PWM) Dec() {
	PWM_retard_count(p.slice_num)
}
