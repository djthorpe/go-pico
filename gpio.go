//go:build !rpi

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

type gpio struct {
	init    []bool
	adcinit bool
	intr    interrupt.Interrupt
}

//////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new GPIO object
//
func _NewGPIO() *gpio {
	g := &gpio{}
	g.init = make([]bool, NUM_BANK0_GPIOS)
	g.intr = interrupt.New(rp.IRQ_IO_IRQ_BANK0, gpio_intr_handler)
	return g
}

// Close GPIO device, return each pin to NULL state
//
func (g *gpio) Close() error {
	for pin := Pin(0); pin < NUM_BANK0_GPIOS; pin++ {
		if g.init[pin] {
			g.deinit(pin)
		}
	}

	// Return success
	return nil
}

//////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - INTERRUPTS

// Handle interrupts
func gpio_intr_handler(interrupt.Interrupt) {
	fmt.Println("Got intr")
	for pin := GPIO_pin(0); pin < NUM_BANK0_GPIOS; pin++ {
		// Get the change for this pin

		GPIO_acknowledge_irq(pin, 0xF)
	}
	// TODO
	/*
			io_irq_ctrl_hw_t *irq_ctrl_base = get_core_num() ?
		                                           &iobank0_hw->proc1_irq_ctrl : &iobank0_hw->proc0_irq_ctrl;
		    for (uint gpio = 0; gpio < NUM_BANK0_GPIOS; gpio++) {
		        io_ro_32 *status_reg = &irq_ctrl_base->ints[gpio / 8];
		        uint events = (*status_reg >> 4 * (gpio % 8)) & 0xf;
		        if (events) {
		            // TODO: If both cores care about this event then the second core won't get the irq?
		            gpio_acknowledge_irq(gpio, events);
		            gpio_irq_callback_t callback = _callbacks[get_core_num()];
		            if (callback) {
		                callback(gpio, events);
		            }
		        }
		    }
	*/
}

//////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Initialise a single pin to a specific mode
//
func (g *gpio) setmode(pin Pin, mode Mode) error {
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
func (g *gpio) deinit(pin Pin) {
	if g.init[pin] {
		g.setInterrupt(pin, nil)
		GPIO_deinit(GPIO_pin(pin))
		g.init[pin] = false
	}
}

// Get mode on a pin
//
func (g *gpio) mode(pin Pin) (Mode, error) {
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

// Get pin state
//
func (g *gpio) get(pin Pin) (bool, error) {
	if err := assert(pin < NUM_BANK0_GPIOS, ErrBadParameter); err != nil {
		return false, err
	}
	if !g.init[pin] {
		if err := g.setmode(pin, ModeInput); err != nil {
			return false, err
		}
	}
	return GPIO_get(GPIO_pin(pin)), nil
}

// Set pin state
//
func (g *gpio) set(pin Pin, value bool) error {
	if err := assert(pin < NUM_BANK0_GPIOS, ErrBadParameter); err != nil {
		return err
	}
	if !g.init[pin] {
		if err := g.setmode(pin, ModeOutput); err != nil {
			return err
		}
	}
	GPIO_put(GPIO_pin(pin), value)
	return nil
}

// Return PWM device on a pin
//
func (g *gpio) pwm(pin Pin) (*PWM, error) {
	// Check parameters
	if err := assert(pin < NUM_BANK0_GPIOS, ErrBadParameter); err != nil {
		return nil, err
	}
	// Set mode
	if mode, err := g.mode(pin); err != nil {
		return nil, err
	} else if mode != ModePWM {
		if err := g.setmode(pin, ModePWM); err != nil {
			return nil, err
		}
	}
	// Return PWM
	return pwm[PWM_gpio_to_slice_num(GPIO_pin(pin))], nil
}

// Return ADC device on a pin
//
func (g *gpio) adc(pin Pin) (*ADC, error) {
	// Check parameters
	if err := assert(pin < NUM_BANK0_GPIOS, ErrBadParameter.With(pin)); err != nil {
		return nil, err
	}

	// Initialise ADC device
	if !g.adcinit {
		ADC_init()
		g.adcinit = true
	}

	// Get input channel, returns
	ch, exists := ADC_gpio_to_input(GPIO_pin(pin))
	if !exists {
		return nil, ErrBadParameter.With(pin)
	} else if err := assert(ch < NUM_ADC_CHANNELS, ErrBadParameter.With(pin)); err != nil {
		return nil, err
	}

	// Set mode
	if mode, err := g.mode(pin); err != nil {
		return nil, err
	} else if mode != ModeOff {
		if err := g.setmode(pin, ModeOff); err != nil {
			return nil, err
		}
	}

	// Initialise pin
	ADC_gpio_init(GPIO_pin(pin))

	// Return channel
	return adc[ch], nil
}

// Return UART device on a pin
//
// TODO

// Add pin handler
//
func (g *gpio) setInterrupt(pin Pin, handler func(pin Pin)) {
	// TODO
	if handler != nil {
		// Enable interrupt handler
		GPIO_set_irq_enabled(GPIO_pin(pin), GPIO_IRQ_EDGE_RISE|GPIO_IRQ_EDGE_FALL, true)
		// Enable ARM interrupt
		g.intr.Enable()
	} else {
		// Diable ARM interrupt
		g.intr.Disable()
		// Disable interrupt handler
		GPIO_set_irq_enabled(GPIO_pin(pin), GPIO_IRQ_EDGE_RISE|GPIO_IRQ_EDGE_FALL, false)
	}
}
