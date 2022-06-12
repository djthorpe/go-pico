package pico

import (
	"fmt"

	// Module imports
	rp "device/rp"
	interrupt "runtime/interrupt"

	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/errors"
	. "github.com/djthorpe/go-pico/pkg/sdk"
)

//////////////////////////////////////////////////////////////////////////////
// TYPES

type GPIO struct {
	// Determine which pins have been initialised
	init [NUM_BANK0_GPIOS]bool
}

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
// CONSTANTS

var (
	gpio      = NewGPIO()
	gpio_intr = interrupt.New(rp.IRQ_IO_IRQ_BANK0, gpio_intr_handler)
)

//////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewGPIO() *GPIO {
	// Initialise GPIO
	gpio := new(GPIO)

	// Initialise PWM
	for slice_num := uint32(0); slice_num < NUM_PWM_SLICES; slice_num++ {
		NewPWM(slice_num)
	}

	// Return GPIO
	return gpio
}

//////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Initialise a single pin to a specific mode
//
func (g *GPIO) initpin(pin Pin, mode Mode) error {
	if err := assert(pin < NUM_BANK0_GPIOS, ErrBadParameter); err != nil {
		return err
	}
	if err := assert(mode <= ModeOff, ErrBadParameter); err != nil {
		return err
	}

	// Init pin
	_pin := GPIO_pin(pin)
	if !g.init[pin] {
		GPIO_init(_pin)
		g.init[pin] = true
	}

	// Set mode
	switch mode {
	case ModeOutput:
		GPIO_set_function(_pin, GPIO_FUNC_SIO)
		GPIO_set_dir(_pin, GPIO_DIR_OUT)
		GPIO_set_output_enabled(_pin, true)
		if err := assert(GPIO_get_function(_pin) == GPIO_FUNC_SIO && GPIO_get_dir(_pin) == GPIO_DIR_OUT && GPIO_get_output_enabled(_pin), ErrUnexpectedValue); err != nil {
			return err
		}
	case ModeInput:
		GPIO_set_function(_pin, GPIO_FUNC_SIO)
		GPIO_set_dir(_pin, GPIO_DIR_IN)
		GPIO_disable_pulls(_pin)
	case ModeInputPulldown:
		GPIO_set_function(_pin, GPIO_FUNC_SIO)
		GPIO_set_dir(_pin, GPIO_DIR_IN)
		GPIO_pull_down(_pin)
	case ModeInputPullup:
		GPIO_set_function(_pin, GPIO_FUNC_SIO)
		GPIO_set_dir(_pin, GPIO_DIR_IN)
		GPIO_pull_up(_pin)
	case ModeI2C:
		// IO config according to 4.3.1.3 of rp2040 datasheet
		GPIO_set_function(_pin, GPIO_FUNC_I2C)
		GPIO_pull_up(_pin)
		GPIO_set_input_hysteresis_enabled(_pin, true)
		GPIO_set_slew_rate(_pin, GPIO_SLEW_RATE_FAST)
	case ModeSPI:
		GPIO_set_function(_pin, GPIO_FUNC_SPI)
	case ModePWM:
		GPIO_set_function(_pin, GPIO_FUNC_PWM)
	case ModeUART:
		GPIO_set_function(_pin, GPIO_FUNC_UART)
	case ModeOff:
		GPIO_set_function(_pin, GPIO_FUNC_NULL)
		GPIO_disable_pulls(_pin)
	}

	// Return success
	return nil
}

// Resets a GPIO back to the NULL function
//
func (g *GPIO) deinit(pin Pin) {
	if g.init[pin] {
		GPIO_deinit(GPIO_pin(pin))
		g.init[pin] = false
	}
}

// Get mode on a pin
//
func (g *GPIO) mode(pin Pin) (Mode, error) {
	if err := assert(pin < NUM_BANK0_GPIOS, ErrBadParameter); err != nil {
		return 0, err
	}
	fn := GPIO_get_function(GPIO_pin(pin))
	switch fn {
	case GPIO_FUNC_SPI:
		return ModeSPI, nil
	case GPIO_FUNC_I2C:
		return ModeI2C, nil
	case GPIO_FUNC_PWM:
		return ModePWM, nil
	case GPIO_FUNC_UART:
		return ModeUART, nil
	case GPIO_FUNC_SIO:
		if GPIO_get_dir(GPIO_pin(pin)) == GPIO_DIR_OUT {
			return ModeOutput, nil
		} else if GPIO_is_pulled_up(GPIO_pin(pin)) {
			return ModeInputPullup, nil
		} else if GPIO_is_pulled_down(GPIO_pin(pin)) {
			return ModeInputPulldown, nil
		} else {
			return ModeInput, nil
		}
	case GPIO_FUNC_NULL:
		return ModeOff, nil
	default:
		return 0, assert(false, ErrUnexpectedValue.With(fn))
	}
}

// Get current value on a pin
//
func (g *GPIO) get(pin Pin) (bool, error) {
	if err := assert(pin < NUM_BANK0_GPIOS, ErrBadParameter); err != nil {
		return false, err
	}
	if err := assert(g.init[pin], ErrNotInitialised); err != nil {
		return false, err
	}
	return GPIO_get(GPIO_pin(pin)), nil
}

// Set current value on a pin
//
func (g *GPIO) set(pin Pin, value bool) error {
	if err := assert(pin < NUM_BANK0_GPIOS, ErrBadParameter); err != nil {
		return err
	}
	if err := assert(g.init[pin], ErrNotInitialised); err != nil {
		return err
	}
	GPIO_put(GPIO_pin(pin), value)
	return nil
}

// Return PWM device on a pin
//
func (g *GPIO) pwm(pin Pin) (*PWM, error) {
	if err := assert(pin < NUM_BANK0_GPIOS, ErrBadParameter); err != nil {
		return nil, err
	}
	if err := assert(g.init[pin], ErrNotInitialised); err != nil {
		return nil, err
	}
	if err := assert(GPIO_get_function(GPIO_pin(pin)) == GPIO_FUNC_PWM, ErrUnexpectedValue); err != nil {
		return nil, err
	}
	return pwm[PWM_gpio_to_slice_num(GPIO_pin(pin))], nil
}

// Return UART device on a pin
//

// Add pin handler
func (g *GPIO) SetInterrupt(pin Pin, handler func(pin Pin)) {
	// Enable interrupt handler

	// Enable ARM interrupt
	gpio_intr.Enable()
}

// Handle interrupts
func gpio_intr_handler(interrupt.Interrupt) {
	fmt.Println("Got intr")
}