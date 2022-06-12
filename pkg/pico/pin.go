package pico

//////////////////////////////////////////////////////////////////////////////
// TYPES

type Pin uint
type Mode uint

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ModeOutput Mode = iota
	ModeInput
	ModeInputPulldown
	ModeInputPullup
	ModeUART
	ModePWM
	ModeI2C
	ModeSPI
	ModeOff
)

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Set pin mode
//
func (p Pin) SetMode(mode Mode) error {
	return gpio.initpin(p, mode)
}

// Set pin state
//
func (p Pin) Set(value bool) {
	gpio.set(p, value)
}

// Get pin state
//
func (p Pin) Get() bool {
	v, _ := gpio.get(p)
	return v
}

// Get PWM for pin
func (p Pin) PWM() *PWM {
	v, _ := gpio.pwm(p)
	return v
}
