//go:build rp2040

package sdk

import (
	// Module imports
	rp "device/rp"
	volatile "runtime/volatile"
	"unsafe"
)

// SDK documentation
// https://github.com/raspberrypi/pico-sdk/tree/master/src/rp2_common/hardware_pwm

//////////////////////////////////////////////////////////////////////////////
// TYPES

type PWM_clkdiv_mode uint32

type PWM_chan uint32

type PWM_config struct {
	csr uint32
	div uint32
	top uint32
}

type pwm_group struct {
	csr volatile.Register32
	div volatile.Register32
	ctr volatile.Register32
	cc  volatile.Register32
	top volatile.Register32
}

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	PWM_DIV_FREE_RUNNING PWM_clkdiv_mode = 0 // Free-running counting at rate dictated by fractional divider
	PWM_DIV_B_HIGH       PWM_clkdiv_mode = 1 // Fractional divider is gated by the PWM B pin
	PWM_DIV_B_RISING     PWM_clkdiv_mode = 2 // Fractional divider advances with each rising edge of the PWM B pin
	PWM_DIV_B_FALLING    PWM_clkdiv_mode = 3 // Fractional divider advances with each falling edge of the PWM B pin
)

const (
	PWM_CHAN_A PWM_chan = 0
	PWM_CHAN_B PWM_chan = 1
)

const (
	_PWM_CH0_CTR_RESET = 0
	_PWM_CH0_CC_RESET  = 0
)

var (
	pwm_groups = *(*[NUM_PWM_SLICES]pwm_group)(unsafe.Pointer(rp.PWM))
)

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Determine the PWM slice that is attached to the specified GPIO
//
//go:inline
func PWM_gpio_to_slice_num(pin GPIO_pin) uint32 {
	assert(pin < NUM_BANK0_GPIOS)
	return uint32(pin>>1) & 7
}

// Determine the PWM channel that is attached to the specified GPIO
//
//go:inline
func PWM_gpio_to_channel(pin GPIO_pin) PWM_chan {
	assert(pin < NUM_BANK0_GPIOS)
	return PWM_chan(uint32(pin) & 1)
}

// Set phase correction in a PWM configuration
//
// Setting phase control to true means that instead of wrapping back to
// zero when the wrap point is reached, the PWM starts counting back down.
// The output frequency is halved when phase-correct mode is enabled.
//
func PWM_config_set_phase_correct(c *PWM_config, phase_correct bool) {
	assert(c != nil)
	c.csr = (c.csr & ^uint32(rp.PWM_CH0_CSR_PH_CORRECT_Msk)) | (bool_to_bit(phase_correct) << rp.PWM_CH0_CSR_PH_CORRECT_Pos)
}

// Set PWM clock divider in a PWM configuration
//
// If the divide mode is free-running, the PWM counter runs at clk_sys / div.
// Otherwise, the divider reduces the rate of events seen on the B pin input (level or edge)
// before passing them on to the PWM counter.
//
func PWM_config_set_clkdiv(c *PWM_config, div float32) {
	assert(c != nil)
	assert(div >= 1.0 && div < 256.0)
	c.div = (uint32)(div * (float32)(1<<rp.PWM_CH0_DIV_INT_Pos))
}

// Set PWM clock divider in a PWM configuration using an 8:4 fractional value
//
// If the divide mode is free-running, the PWM counter runs at clk_sys / div.
// Otherwise, the divider reduces the rate of events seen on the B pin input (level or edge)
// before passing them on to the PWM counter.
//
func PWM_config_set_clkdiv_int_frac(c *PWM_config, integer, fract uint8) {
	assert(c != nil)
	assert(integer >= 1)
	assert(fract < 16)
	c.div = (uint32(integer) << rp.PWM_CH0_DIV_INT_Pos) | (uint32(fract) << rp.PWM_CH0_DIV_FRAC_Pos)
}

// Set PWM clock divider in a PWM configuration
//
// If the divide mode is free-running, the PWM counter runs at clk_sys / div.
// Otherwise, the divider reduces the rate of events seen on the B pin input (level or edge)
// before passing them on to the PWM counter.
//
func PWM_config_set_clkdiv_int(c *PWM_config, div uint32) {
	assert(c != nil)
	assert(div >= 1 && div < 256)
	PWM_config_set_clkdiv_int_frac(c, uint8(div), 0)
}

// Set PWM counting mode in a PWM configuration
//
// Configure which event gates the operation of the fractional divider.
// The default is always-on (free-running PWM). Can also be configured to count on
// high level, rising edge or falling edge of the B pin input.
//
func PWM_config_set_clkdiv_mode(c *PWM_config, mode PWM_clkdiv_mode) {
	assert(c != nil)
	assert(mode == PWM_DIV_FREE_RUNNING || mode == PWM_DIV_B_RISING || mode == PWM_DIV_B_HIGH || mode == PWM_DIV_B_FALLING)
	c.csr = (c.csr & ^uint32(rp.PWM_CH0_CSR_DIVMODE_Msk)) | (uint32(mode) << rp.PWM_CH0_CSR_DIVMODE_Pos)
}

// Set output polarity in a PWM configuration
//
// Set a or b to true to inverse the polarity of the output on channel a or b.
//
func PWM_config_set_output_polarity(c *PWM_config, a, b bool) {
	assert(c != nil)
	c.csr = (c.csr & ^uint32(rp.PWM_CH0_CSR_A_INV|rp.PWM_CH0_CSR_B_INV))
	c.csr |= (bool_to_bit(a) << rp.PWM_CH0_CSR_A_INV_Pos)
	c.csr |= (bool_to_bit(b) << rp.PWM_CH0_CSR_B_INV_Pos)
}

// Set PWM counter wrap value in a PWM configuration
//
// Set the highest value the counter will reach before returning to 0. Also known as TOP.
//
func PWM_config_set_wrap(c *PWM_config, wrap uint16) {
	assert(c != nil)
	c.top = uint32(wrap)
}

// Initialise a PWM with settings from a configuration object
//
// If start is set the PWM will be started running once configured. If false you will need to start
// manually using PWM_set_enabled() or PWM_set_mask_enabled()
//
func PWM_init(slice_num uint32, c *PWM_config, start bool) {
	assert(slice_num < NUM_PWM_SLICES)
	assert(c != nil)

	pwm_groups[slice_num].csr.Set(0)
	pwm_groups[slice_num].ctr.Set(_PWM_CH0_CTR_RESET)
	pwm_groups[slice_num].cc.Set(_PWM_CH0_CC_RESET)
	pwm_groups[slice_num].top.Set(c.top)
	pwm_groups[slice_num].div.Set(c.div)
	pwm_groups[slice_num].csr.Set(c.csr | (bool_to_bit(start) << rp.PWM_CH0_CSR_EN_Pos))
}

// Get a set of default values for PWM configuration
//
func PWM_get_default_config() *PWM_config {
	c := new(PWM_config)
	PWM_config_set_phase_correct(c, false)
	PWM_config_set_clkdiv_int(c, 1)
	PWM_config_set_clkdiv_mode(c, PWM_DIV_FREE_RUNNING)
	PWM_config_set_output_polarity(c, false, false)
	PWM_config_set_wrap(c, 0xFFFF)
	return c
}

// Set the current PWM counter wrap value
//
func PWM_set_wrap(slice_num uint32, wrap uint16) {
	assert(slice_num < NUM_PWM_SLICES)
	pwm_groups[slice_num].top.Set(uint32(wrap))
}

// Set the current PWM counter compare value for one channel
//
func PWM_set_chan_level(slice_num uint32, ch PWM_chan, level uint16) {
	assert(slice_num < NUM_PWM_SLICES)
	assert(ch == PWM_CHAN_A || ch == PWM_CHAN_B)
	switch ch {
	case PWM_CHAN_A:
		pwm_groups[slice_num].cc.ReplaceBits(uint32(level)<<rp.PWM_CH0_CC_A_Pos, rp.PWM_CH0_CC_A_Msk, 0)
	case PWM_CHAN_B:
		pwm_groups[slice_num].cc.ReplaceBits(uint32(level)<<rp.PWM_CH0_CC_B_Pos, rp.PWM_CH0_CC_B_Msk, 0)
	}
}

// Set PWM counter compare values
//
func PWM_set_both_levels(slice_num uint32, levela, levelb uint16) {
	assert(slice_num < NUM_PWM_SLICES)
	pwm_groups[slice_num].cc.Set((uint32(levela) << rp.PWM_CH0_CC_A_Pos) | (uint32(levelb) << rp.PWM_CH0_CC_B_Pos))
}

// Helper function to set the PWM level for the slice and channel associated with a GPIO.
//
// This PWM slice should already have been configured and set running. Also be
// careful of multiple GPIOs mapping to the same slice and channel (if GPIOs
// have a difference of 16).
//
func PWM_set_gpio_level(pin GPIO_pin, level uint16) {
	assert(pin < NUM_BANK0_GPIOS)
	PWM_set_chan_level(PWM_gpio_to_slice_num(pin), PWM_gpio_to_channel(pin), level)
}

// Get PWM counter
//
func PWM_get_counter(slice_num uint32) uint16 {
	assert(slice_num < NUM_PWM_SLICES)
	return uint16(pwm_groups[slice_num].ctr.Get())
}

// Set PWM counter
//
func PWM_set_counter(slice_num uint32, c uint16) {
	assert(slice_num < NUM_PWM_SLICES)
	pwm_groups[slice_num].ctr.Set(uint32(c))
}

// Advance PWM count and wait until advanced
//
func PWM_advance_count(slice_num uint32) {
	assert(slice_num < NUM_PWM_SLICES)
	pwm_groups[slice_num].csr.SetBits(rp.PWM_CH0_CSR_PH_ADV_Msk)
	for {
		if !pwm_groups[slice_num].csr.HasBits(rp.PWM_CH0_CSR_PH_ADV_Msk) {
			break
		}
	}
}

// Retard PWM count and wait
//
func PWM_retard_count(slice_num uint32) {
	assert(slice_num < NUM_PWM_SLICES)
	pwm_groups[slice_num].csr.SetBits(rp.PWM_CH0_CSR_PH_RET_Msk)
	for {
		if !pwm_groups[slice_num].csr.HasBits(rp.PWM_CH0_CSR_PH_RET_Msk) {
			break
		}
	}
}

// Set PWM clock divider using an 8:4 fractional value
//
/*
func PWM_set_clkdiv_int_frac(slice_num uint32,integer,fract uint8) {
    check_slice_num_param(slice_num);
	assert(integer >= 1)
	assert(fract < 16)
    pwm_hw->slice[slice_num].div = (((uint)integer) << PWM_CH0_DIV_INT_LSB) | (((uint)fract) << PWM_CH0_DIV_FRAC_LSB);
}
*/
