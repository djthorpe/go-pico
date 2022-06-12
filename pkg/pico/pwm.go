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

func (p *PWM) SetEnabled(enable bool) {
	PWM_config_set_clkdiv(p.config, 4)
	if enable {
		PWM_init(p.slice_num, p.config, true)
	}
}

func (p *PWM) Get() uint16 {
	return PWM_get_counter(p.slice_num)
}

func (p *PWM) Set(value uint16) {
	PWM_set_counter(p.slice_num, value)
}
