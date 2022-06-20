//go:build pico

package pico

import (
	// Module imports
	rp "device/rp"
	interrupt "runtime/interrupt"

	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/errors"
	. "github.com/djthorpe/go-pico/pkg/sdk"
)

//////////////////////////////////////////////////////////////////////////////
// TYPES

type PWM struct {
	slice_num uint32
	config    *PWM_config
	intr      interrupt.Interrupt
}

type PWM_callback_t func(pwm *PWM)

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	_PWM_MAX_TOP      = 0xFFFF
	_PWM_DEFAULT_TOP  = 95 * _PWM_MAX_TOP / 100 // start algorithm at 95% Top. This allows us to undershoot period with prescale.
	_PWM_MILLISECONDS = 1_000_000_000
	_PWM_MIN_PERIOD   = 8                       // Minimum period is 8ns
	_PWM_MAX_PERIOD   = 268 * _PWM_MILLISECONDS // Maximum Period is 268369920ns on rp2040, given by (16*255+15)*8*(1+0xffff)*(1+1)/16
)

var (
	pwm           = [NUM_PWM_SLICES]*PWM{}
	pwm_callbacks = [NUM_PWM_SLICES]PWM_callback_t{}
)

//////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func _NewPWM(slice_num uint32) *PWM {
	if slice_num >= NUM_PWM_SLICES {
		return nil
	} else if pwm := pwm[slice_num]; pwm != nil {
		return pwm
	}

	// Initialise a new PWM
	pwm[slice_num] = &PWM{
		slice_num: slice_num,
		config:    PWM_get_default_config(),
		intr:      interrupt.New(rp.IRQ_PWM_IRQ_WRAP, intr_handler),
	}

	// Return the PWM
	return pwm[slice_num]
}

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (p *PWM) SetEnabled(enabled bool) {
	if enabled {
		PWM_config_set_clkdiv(p.config, 16)
		PWM_init(p.slice_num, p.config, true)
	} else {
		PWM_set_enabled(p.slice_num, enabled)
	}
}

func (p *PWM) Enabled() bool {
	return PWM_is_enabled(p.slice_num)
}

// Set level
func (p *PWM) Set(pin Pin, level uint16) {
	PWM_set_gpio_level(GPIO_pin(pin), level)
}

// Get level
func (p *PWM) Get(pin Pin) uint16 {
	return PWM_get_gpio_level(GPIO_pin(pin))
}

// Get counter value
//
func (p *PWM) Counter() uint16 {
	return PWM_get_counter(p.slice_num)
}

// Set counter value
//
func (p *PWM) SetCounter(value uint16) {
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

// Set wrapping value
//
func (p *PWM) SetWrap(wrap uint16) {
	PWM_set_wrap(p.slice_num, wrap)
	PWM_config_set_wrap(p.config, wrap)
}

// Get wrapping value
//
func (p *PWM) Wrap() uint16 {
	return PWM_get_wrap(p.slice_num)
}

// Set period in nanoseconds
//
func (p *PWM) SetPeriod(period uint64) error {
	if err := assert(period >= _PWM_MIN_PERIOD && period <= _PWM_MAX_PERIOD, ErrBadParameter.With("SetPeriod:", period)); err != nil {
		return err
	}

	// Must enable Phase correct to reach large periods.
	if period > (_PWM_MAX_PERIOD >> 1) {
		PWM_set_phase_correct(p.slice_num, true)
	}

	// clearing above expression:
	//  DIV_INT + DIV_FRAC/16 = cycles / ( (TOP+1) * (CSRPHCorrect+1) )  // DIV_FRAC/16 is always 0 in this equation
	// where cycles must be converted to time:
	//  target_period = cycles * period_per_cycle ==> cycles = target_period/period_per_cycle
	period_per_cycle := uint64(cpu_period())
	phc := uint64(PWM_get_phase_correct(p.slice_num))
	wrap := PWM_get_wrap(p.slice_num)
	rhs := 16 * period / ((1 + phc) * period_per_cycle * (1 + wrap)) // right-hand-side of equation, scaled so frac is not divided
	whole := rhs >> 4
	frac := rhs & 0x0F
	switch {
	case whole > 0xFF:
		whole = 0xFF
	case whole == 0:
		whole = 1
		frac = 0
	}

	// Step 2 is acquiring a better top value. Clearing the equation:
	// TOP =  cycles / ( (DIVINT+DIVFRAC/16) * (CSRPHCorrect+1) ) - 1
	top := 16*period/((16*whole+frac)*periodPerCycle*(1+phc)) - 1
	if top > maxTop {
		top = maxTop
	}
	pwm.SetTop(uint32(top))
	pwm.setClockDiv(uint8(whole), uint8(frac))
	return nil
}

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - INTERRUPTS

// Set interrupt handler
//
// If called with nil then handler is disabled
//
func (p *PWM) SetInterrupt(handler PWM_callback_t) {
	// Enable interrupt handler
	PWM_clear_irq(p.slice_num)
	if handler == nil {
		PWM_set_irq_enabled(p.slice_num, false)
	} else {
		PWM_set_irq_enabled(p.slice_num, true)
	}

	// Set callback
	pwm_callbacks[p.slice_num] = handler

	// Enable ARM interrupt
	if handler == nil {
		p.intr.Disable()
	} else {
		p.intr.Enable()
	}
}

//////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - INTERRUPTS

// Interrupt handler
//
func intr_handler(interrupt.Interrupt) {
	mask := PWM_get_irq_mask()
	PWM_clear_irq_mask(mask)
	for slice_num := uint32(0); slice_num < NUM_PWM_SLICES; slice_num++ {
		if mask&1 != 0 {
			if fn := pwm_callbacks[slice_num]; fn != nil {
				fn(pwm[slice_num])
			}
		}
		mask >>= 1
	}
}
