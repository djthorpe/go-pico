package pico

//////////////////////////////////////////////////////////////////////////////
// TYPES

type Pin uint
type Pin_callback_t func(Pin)

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Get pin mode
func (p Pin) Mode() Mode {
	if mode, err := _GPIO.mode(p); err != nil {
		return 0
	} else {
		return mode
	}
}

// Set pin mode
//
func (p Pin) SetMode(mode Mode) error {
	return _GPIO.setmode(p, mode)
}

// Set pin state
//
func (p Pin) Set(value bool) {
	_GPIO.set(p, value)
}

// Get pin state
//
func (p Pin) Get() bool {
	v, _ := _GPIO.get(p)
	return v
}

// Get PWM for pin
func (p Pin) PWM() *PWM {
	if pwm, err := _GPIO.pwm(p); err != nil {
		return nil
	} else {
		return pwm
	}
}

// Get ADC for pin
func (p Pin) ADC() *ADC {
	if adc, err := _GPIO.adc(p); err != nil {
		return nil
	} else {
		return adc
	}
}

// Set pin interrupt
//
func (p Pin) SetInterrupt(callback Pin_callback_t) {
	_GPIO.setInterrupt(p, callback)
}
