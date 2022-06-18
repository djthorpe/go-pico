//go:build rp2040

package sdk

import (
	// Module imports
	rp "device/rp"
	volatile "runtime/volatile"
	"unsafe"
)

// SDK documentation
// https://github.com/raspberrypi/pico-sdk/blob/master/src/rp2_common/hardware_adc

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

type adc_t struct {
	cs     volatile.Register32 // 0x0
	result volatile.Register32 // 0x4
	fcs    volatile.Register32 // 0x8
	fifo   volatile.Register32 // 0xC
	div    volatile.Register32 // 0x10
	intr   volatile.Register32 // 0x14
	inte   volatile.Register32 // 0x18
	intf   volatile.Register32 // 0x1C
	ints   volatile.Register32 // 0x20
}

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ADC_BANK0_GPIOS_MIN = 26
	ADC_BANK0_GPIOS_MAX = 29
	ADC_NUM_CHANNELS    = 4
)

var (
	adc = (*adc_t)(unsafe.Pointer(rp.ADC))
)

//////////////////////////////////////////////////////////////////////////////
// METHODS

//  Initialise the ADC
//
func ADC_init() {
	// ADC is in an unknown state. We should start by resetting it
	reset_block(rp.RESETS_RESET_ADC)
	unreset_block_wait(rp.RESETS_RESET_ADC)

	// Now turn it back on. Staging of clock etc is handled internally
	adc.cs.SetBits(rp.ADC_CS_EN)

	// Internal staging completes in a few cycles, but poll to be sure
	for {
		if ADC_is_ready() {
			break
		}
	}
}

// Return ready status for ADC
//
//go:inline
func ADC_is_ready() bool {
	return adc.cs.HasBits(rp.ADC_CS_READY)
}

// Initialise the gpio for use as an ADC pin
//
//go:inline
func ADC_gpio_init(pin GPIO_pin) {
	assert(pin >= ADC_BANK0_GPIOS_MIN && pin <= ADC_BANK0_GPIOS_MAX)
	// Select NULL function to make output driver hi-Z
	GPIO_set_function(pin, GPIO_FUNC_NULL)
	// Also disable digital pulls and digital receiver
	GPIO_disable_pulls(pin)
	GPIO_set_input_enabled(pin, false)
}

// ADC input select
//
// Select an ADC input. 0...3 are GPIOs 26...29 respectively.
// Input 4 is the onboard temperature sensor.
//
//go:inline
func ADC_select_input(input uint32) {
	assert(input <= ADC_NUM_CHANNELS)
	v := input << rp.ADC_CS_AINSEL_Pos
	m := uint32(rp.ADC_CS_AINSEL_Msk)
	adc.cs.ReplaceBits(v, m, 0)
}

// Get the currently selected ADC input channel
//
//go:inline
func ADC_get_selected_input() uint32 {
	return (adc.cs.Get() & rp.ADC_CS_AINSEL_Msk) >> rp.ADC_CS_AINSEL_Pos
}

// Round Robin sampling selector
//
// This function sets which inputs are to be run through in round robin mode.
// Value between 0 and 0x1f (bit 0 to bit 4 for GPIO 26 to 29 and temperature sensor
// input respectively) Write a value of 0 to disable round robin sampling.
//
//go:inline
func ADC_set_round_robin(input_mask uint32) {
	assert(input_mask < (1 << NUM_ADC_CHANNELS))
	v := input_mask << rp.ADC_CS_RROBIN_Pos
	m := uint32(rp.ADC_CS_RROBIN_Msk)
	adc.cs.ReplaceBits(v, m, 0)
}

// Enable the onboard temperature sensor
//
//go:inline
func ADC_set_temp_sensor_enabled(enable bool) {
	if enable {
		adc.cs.SetBits(rp.ADC_CS_TS_EN)
	} else {
		adc.cs.ClearBits(rp.ADC_CS_TS_EN)
	}
}

// Perform a single conversion
//
//go:inline
func ADC_read() uint16 {
	adc.cs.SetBits(rp.ADC_CS_START_ONCE)
	for {
		if adc.cs.HasBits(rp.ADC_CS_READY) {
			break
		}
	}
	return uint16(adc.result.Get() & rp.ADC_RESULT_RESULT_Msk)
}

// Enable or disable free-running sampling mode
//
//go:inline
func ADC_run(run bool) {
	if run {
		adc.cs.SetBits(rp.ADC_CS_START_MANY)
	} else {
		adc.cs.ClearBits(rp.ADC_CS_START_MANY)
	}
}

/*
// Set the ADC clock divisor
//
// Period of samples will be (1 + div) cycles on average. Note it takes 96 cycles to
// perform a conversion, so any period less than that will be clamped to 96.
//
//go:inline
func ADC_set_clkdiv(clkdiv float32) {
	// TODO
    invalid_params_if(ADC, clkdiv >= 1 << (ADC_DIV_INT_MSB - ADC_DIV_INT_LSB + 1));
	adc.div.Set(v)
    adc_hw->div = (uint32_t)(clkdiv * (float) (1 << ADC_DIV_INT_LSB));
}
*/

// Setup the ADC FIFO
//
// FIFO is 4 samples long, if a conversion is completed and the FIFO is full, the result is dropped
//
//go:inline
func ADC_fifo_setup(en, dreq_en bool, dreq_thresh uint16, err_in_fifo, byte_shift bool) {
	v := bool_to_bit(en) << rp.ADC_FCS_EN_Pos
	v |= bool_to_bit(dreq_en) << rp.ADC_FCS_DREQ_EN_Pos
	v |= uint32(dreq_thresh) << rp.ADC_FCS_THRESH_Pos
	v |= bool_to_bit(err_in_fifo) << rp.ADC_FCS_ERR_Pos
	v |= bool_to_bit(byte_shift) << rp.ADC_FCS_SHIFT_Pos
	m := uint32(rp.ADC_FCS_EN_Msk | rp.ADC_FCS_DREQ_EN_Msk | rp.ADC_FCS_THRESH_Msk | rp.ADC_FCS_ERR_Msk | rp.ADC_FCS_SHIFT_Msk)
	adc.fcs.ReplaceBits(v, m, 0)
}

// Check FIFO empty state
//
//go:inline
func ADC_fifo_is_empty() bool {
	return adc.fcs.HasBits(rp.ADC_FCS_EMPTY)
}

// Get number of entries in the ADC FIFO
//
// The ADC FIFO is 4 entries long. This function will return how many samples are
// currently present.
//
//go:inline
func ADC_fifo_get_level() uint8 {
	return uint8((adc.fcs.Get() & rp.ADC_FCS_LEVEL_Msk) >> rp.ADC_FCS_LEVEL_Pos)
}

// Get ADC result from FIFO
//
// Pops the latest result from the ADC FIFO.
//
//go:inline
func ADC_fifo_get() uint16 {
	return uint16(adc.fifo.Get() & rp.ADC_FIFO_VAL_Msk)
}

// Wait for the ADC FIFO to have data.
//
// Blocks until data is present in the FIFO
//
//go:inline
func ADC_fifo_get_blocking() uint16 {
	for {
		if ADC_fifo_is_empty() {
			break
		}
	}
	return ADC_fifo_get()
}

// Drain the ADC FIFO
//
// Will wait for any conversion to complete then drain the FIFO, discarding any results
//
//go:inline
func ADC_fifo_drain() {
	// Potentially there is still a conversion in progress
	// wait for this to complete before draining
	for {
		if ADC_is_ready() {
			break
		}
	}
	// Drain FIFO
	for {
		if ADC_fifo_is_empty() {
			break
		}
		ADC_fifo_get()
	}
}

// Enable/Disable ADC interrupts.
//
//go:inline
func ADC_irq_set_enabled(enabled bool) {
	adc.inte.Set(bool_to_bit(enabled))
}
