package sdk

import (
	"fmt"
	"unsafe"

	// Module imports
	rp "device/rp"
	interrupt "runtime/interrupt"
	volatile "runtime/volatile"
)

// SDK documentation
// https://github.com/raspberrypi/pico-sdk/blob/master/src/rp2_common/hardware_gpio

//////////////////////////////////////////////////////////////////////////////
// TYPES

type GPIO_function uint32
type GPIO_irq_level uint8
type GPIO_override uint8
type GPIO_slew_rate uint8
type GPIO_drive_strength uint8
type GPIO_pin uint8

type gpio_pads_bank0_t struct {
	voltage_select volatile.Register32
	gpio           [30]volatile.Register32
}

type gpio_io_t struct {
	status volatile.Register32
	ctrl   volatile.Register32
}

type gpio_irqctrl_t struct {
	inte [4]volatile.Register32 // enable
	intf [4]volatile.Register32 // force
	ints [4]volatile.Register32 // status
}

type gpio_bank0_t struct {
	gpio               [30]gpio_io_t
	intr               [4]volatile.Register32
	proc0IRQctrl       gpio_irqctrl_t
	proc1IRQctrl       gpio_irqctrl_t
	dormantWakeIRQctrl gpio_irqctrl_t
}

type gpio_irq_callback_t func(interrupt.Interrupt)

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	GPIO_FUNC_XIP  GPIO_function = 0
	GPIO_FUNC_SPI  GPIO_function = 1
	GPIO_FUNC_UART GPIO_function = 2
	GPIO_FUNC_I2C  GPIO_function = 3
	GPIO_FUNC_PWM  GPIO_function = 4
	GPIO_FUNC_SIO  GPIO_function = 5
	GPIO_FUNC_PIO0 GPIO_function = 6
	GPIO_FUNC_PIO1 GPIO_function = 7
	GPIO_FUNC_GPCK GPIO_function = 8
	GPIO_FUNC_USB  GPIO_function = 9
	GPIO_FUNC_NULL GPIO_function = 0x1f
)

const (
	GPIO_OVERRIDE_NORMAL GPIO_override = 0 // peripheral signal selected via gpio_set_function
	GPIO_OVERRIDE_INVERT GPIO_override = 1 // invert peripheral signal selected via gpio_set_function
	GPIO_OVERRIDE_LOW    GPIO_override = 2 // drive low/disable output
	GPIO_OVERRIDE_HIGH   GPIO_override = 3 // drive high/enable output
)

const (
	GPIO_SLEW_RATE_SLOW GPIO_slew_rate = 0 // Slew rate limiting enabled
	GPIO_SLEW_RATE_FAST GPIO_slew_rate = 1 // Slew rate limiting disabled
)

const (
	GPIO_DRIVE_STRENGTH_2MA  GPIO_drive_strength = 0 ///< 2 mA nominal drive strength
	GPIO_DRIVE_STRENGTH_4MA  GPIO_drive_strength = 1 ///< 4 mA nominal drive strength
	GPIO_DRIVE_STRENGTH_8MA  GPIO_drive_strength = 2 ///< 8 mA nominal drive strength
	GPIO_DRIVE_STRENGTH_12MA GPIO_drive_strength = 3 ///< 12 mA nominal drive strength
)

const (
	GPIO_IRQ_LEVEL_NONE GPIO_irq_level = 0
	GPIO_IRQ_LEVEL_LOW  GPIO_irq_level = 1
	GPIO_IRQ_LEVEL_HIGH GPIO_irq_level = 2
	GPIO_IRQ_EDGE_FALL  GPIO_irq_level = 4
	GPIO_IRQ_EDGE_RISE  GPIO_irq_level = 8
	GPIO_IRQ_LEVEL_MAX  GPIO_irq_level = GPIO_IRQ_EDGE_RISE
)

const (
	GPIO_DIR_IN  = 0
	GPIO_DIR_OUT = 1
)

var (
	gpio_pads_bank0   = (*gpio_pads_bank0_t)(unsafe.Pointer(rp.PADS_BANK0))
	gpio_io_bank0     = (*gpio_bank0_t)(unsafe.Pointer(rp.IO_BANK0))
	gpio_irq_callback = [NUM_CORES]gpio_irq_callback_t{}
	gpio_raw_irq_mask = [NUM_CORES]uint32{}
)

//////////////////////////////////////////////////////////////////////////////
// METHODS

// Initialise a GPIO for (enabled I/O and set func to GPIO_FUNC_SIO)
func GPIO_init(pin GPIO_pin) {
	assert(pin < NUM_BANK0_GPIOS)
	rp.SIO.GPIO_OE_CLR.Set(uint32(1) << pin)
	rp.SIO.GPIO_OUT_CLR.Set(uint32(1) << pin)
	GPIO_set_function(pin, GPIO_FUNC_SIO)
}

// Resets a GPIO back to the NULL function, i.e. disables it.
func GPIO_deinit(pin GPIO_pin) {
	assert(pin < NUM_BANK0_GPIOS)
	GPIO_set_function(pin, GPIO_FUNC_NULL)
}

// Initialise multiple GPIOs (enabled I/O and set func to GPIO_FUNC_SIO)
func GPIO_init_mask(mask uint32) {
	for pin := GPIO_pin(0); pin < NUM_BANK0_GPIOS; pin++ {
		if mask&1 != 0 {
			GPIO_init(pin)
		}
		mask >>= 1
	}
}

// Select function for this GPIO, and ensure input/output are enabled at the pad.
// This also clears the input/output/irq override bits
func GPIO_set_function(pin GPIO_pin, fn GPIO_function) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(fn <= GPIO_FUNC_NULL)

	// Set input enable, clear output disable
	gpio_pads_bank0.gpio[pin].ReplaceBits(rp.PADS_BANK0_GPIO0_IE, rp.PADS_BANK0_GPIO0_IE_Msk|rp.PADS_BANK0_GPIO0_OD_Msk, 0)

	// Zero all fields apart from fsel; we want this IO to do what the peripheral tells it.
	// This doesn't affect e.g. pullup/pulldown, as these are in pad controls.
	gpio_io_bank0.gpio[pin].ctrl.Set(uint32(fn) << rp.IO_BANK0_GPIO0_CTRL_FUNCSEL_Pos)
}

// Return current function for this GPIO
func GPIO_get_function(pin GPIO_pin) GPIO_function {
	assert(pin < NUM_BANK0_GPIOS)
	return GPIO_function((gpio_io_bank0.gpio[pin].ctrl.Get() & rp.IO_BANK0_GPIO0_CTRL_FUNCSEL_Msk) >> rp.IO_BANK0_GPIO0_CTRL_FUNCSEL_Pos)
}

// Select up and down pulls on specific GPIO
func GPIO_set_pulls(pin GPIO_pin, up, down bool) {
	assert(pin < NUM_BANK0_GPIOS)
	gpio_pads_bank0.gpio[pin].ReplaceBits(bool_to_bit(up)<<rp.PADS_BANK0_GPIO0_PUE_Pos|bool_to_bit(down)<<rp.PADS_BANK0_GPIO0_PDE_Pos, rp.PADS_BANK0_GPIO0_PUE_Msk|rp.PADS_BANK0_GPIO0_PDE_Msk, 0)
}

// Set specified GPIO to be pulled up
func GPIO_pull_up(pin GPIO_pin) {
	GPIO_set_pulls(pin, true, false)
}

// Set specified GPIO to be pulled down
func GPIO_pull_down(pin GPIO_pin) {
	GPIO_set_pulls(pin, false, true)
}

// Set specified GPIO to be floating
func GPIO_disable_pulls(pin GPIO_pin) {
	GPIO_set_pulls(pin, false, false)
}

// Determine if the specified GPIO is pulled up
func GPIO_is_pulled_up(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return gpio_pads_bank0.gpio[pin].HasBits(rp.PADS_BANK0_GPIO0_PUE_Msk)
}

// Determine if the specified GPIO is pulled down
func GPIO_is_pulled_down(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return gpio_pads_bank0.gpio[pin].HasBits(rp.PADS_BANK0_GPIO0_PDE_Msk)
}

// Set GPIO IRQ override
func GPIO_set_irqover(pin GPIO_pin, value GPIO_override) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(value <= GPIO_OVERRIDE_HIGH)
	gpio_io_bank0.gpio[pin].ctrl.ReplaceBits(uint32(value)<<rp.IO_BANK0_GPIO0_CTRL_IRQOVER_Pos, rp.IO_BANK0_GPIO0_CTRL_IRQOVER_Msk, 0)
}

// Set GPIO output override
func GPIO_set_outover(pin GPIO_pin, value GPIO_override) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(value <= GPIO_OVERRIDE_HIGH)
	gpio_io_bank0.gpio[pin].ctrl.ReplaceBits(uint32(value)<<rp.IO_BANK0_GPIO0_CTRL_OUTOVER_Pos, rp.IO_BANK0_GPIO0_CTRL_OUTOVER_Msk, 0)
}

// Set GPIO input override
func GPIO_set_inover(pin GPIO_pin, value GPIO_override) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(value <= GPIO_OVERRIDE_HIGH)
	gpio_io_bank0.gpio[pin].ctrl.ReplaceBits(uint32(value)<<rp.IO_BANK0_GPIO0_CTRL_INOVER_Pos, rp.IO_BANK0_GPIO0_CTRL_INOVER_Msk, 0)
}

// Set GPIO output enable override
func GPIO_set_oeover(pin GPIO_pin, value GPIO_override) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(value <= GPIO_OVERRIDE_HIGH)
	gpio_io_bank0.gpio[pin].ctrl.ReplaceBits(uint32(value)<<rp.IO_BANK0_GPIO0_CTRL_OEOVER_Pos, rp.IO_BANK0_GPIO0_CTRL_OEOVER_Msk, 0)
}

// Enable GPIO input
func GPIO_set_input_enabled(pin GPIO_pin, enabled bool) {
	assert(pin < NUM_BANK0_GPIOS)
	gpio_pads_bank0.gpio[pin].ReplaceBits(bool_to_bit(enabled)<<rp.PADS_BANK0_GPIO0_IE_Pos, rp.PADS_BANK0_GPIO0_IE_Msk, 0)
}

// Set or clear GPIO output enabled
func GPIO_set_output_enabled(pin GPIO_pin, enabled bool) {
	assert(pin < NUM_BANK0_GPIOS)
	if enabled {
		rp.SIO.GPIO_OE_SET.Set(1 << pin)
	} else {
		rp.SIO.GPIO_OE_CLR.Set(1 << pin)
	}
}

// Get GPIO output enabled state
func GPIO_get_output_enabled(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return rp.SIO.GPIO_OE.HasBits(1 << pin)
}

// Enable/disable GPIO input hysteresis (Schmitt trigger)
func GPIO_set_input_hysteresis_enabled(pin GPIO_pin, enabled bool) {
	assert(pin < NUM_BANK0_GPIOS)
	gpio_pads_bank0.gpio[pin].ReplaceBits(bool_to_bit(enabled)<<rp.PADS_BANK0_GPIO0_SCHMITT_Pos, rp.PADS_BANK0_GPIO0_SCHMITT_Msk, 0)
}

// Determine whether input hysteresis is enabled on a specified GPIO
func GPIO_is_input_hysteresis_enabled(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return gpio_pads_bank0.gpio[pin].HasBits(rp.PADS_BANK0_GPIO0_SCHMITT_Msk)
}

// Set slew rate for a specified GPIO
func GPIO_set_slew_rate(pin GPIO_pin, slew GPIO_slew_rate) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(slew <= GPIO_SLEW_RATE_FAST)
	gpio_pads_bank0.gpio[pin].ReplaceBits(uint32(slew)<<rp.PADS_BANK0_GPIO0_SLEWFAST_Pos, rp.PADS_BANK0_GPIO0_SLEWFAST_Msk, 0)
}

// Determine current slew rate for a specified GPIO
func GPIO_get_slew_rate(pin GPIO_pin) GPIO_slew_rate {
	assert(pin < NUM_BANK0_GPIOS)
	return GPIO_slew_rate((gpio_pads_bank0.gpio[pin].Get() & rp.PADS_BANK0_GPIO0_SLEWFAST) >> rp.PADS_BANK0_GPIO0_SLEWFAST_Pos)
}

// Set drive strength for a specified GPIO
func GPIO_set_drive_strength(pin GPIO_pin, drive GPIO_drive_strength) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(drive <= GPIO_DRIVE_STRENGTH_12MA)
	gpio_pads_bank0.gpio[pin].ReplaceBits(uint32(drive)<<rp.PADS_BANK0_GPIO0_DRIVE_Pos, rp.PADS_BANK0_GPIO0_DRIVE_Msk, 0)
}

// Determine current slew rate for a specified GPIO
func GPIO_get_drive_strength(pin GPIO_pin) GPIO_drive_strength {
	assert(pin < NUM_BANK0_GPIOS)
	return GPIO_drive_strength((gpio_pads_bank0.gpio[pin].Get() & rp.PADS_BANK0_GPIO0_DRIVE_Msk) >> rp.PADS_BANK0_GPIO0_DRIVE_Msk)
}

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - INTERRUPT

func GPIO_set_irq_enabled(pin GPIO_pin, events GPIO_irq_level, enabled bool) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(events <= (GPIO_IRQ_LEVEL_MAX<<1)-1)

	switch get_core_num() {
	case 0:
		gpio_set_irq_enabled(pin, events, enabled, &gpio_io_bank0.proc0IRQctrl)
	case 1:
		gpio_set_irq_enabled(pin, events, enabled, &gpio_io_bank0.proc1IRQctrl)
	default:
		assert(false)
	}
}

func GPIO_acknowledge_irq(pin GPIO_pin, events GPIO_irq_level) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(events <= (GPIO_IRQ_LEVEL_MAX<<1)-1)
	gpio_acknowledge_irq(pin, events)
}

func GPIO_irq_status(pin GPIO_pin) GPIO_irq_level {
	assert(pin < NUM_BANK0_GPIOS)
	return gpio_irq_status(pin)
}

func GPIO_default_irq_handler(interrupt.Interrupt) {
	switch get_core_num() {
	case 0:
		fmt.Println("IRQ Handler 0")
		gpio_default_irq_handler(gpio_irq_callback[0], &gpio_io_bank0.proc0IRQctrl, gpio_raw_irq_mask[0])
	case 1:
		fmt.Println("IRQ Handler 1")
		gpio_default_irq_handler(gpio_irq_callback[1], &gpio_io_bank0.proc1IRQctrl, gpio_raw_irq_mask[1])
	default:
		assert(false)
	}
}

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - INPUT

// Get state of a single specified GPIO
//
//go:inline
func GPIO_get(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return rp.SIO.GPIO_IN.HasBits(1 << pin)
}

// Get state of all GPIO pins
//
//go:inline
func GPIO_get_all() uint32 {
	return rp.SIO.GPIO_IN.Get()
}

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - OUTPUT

// Drive high every GPIO appearing in mask
//
//go:inline
func GPIO_set_mask(mask uint32) {
	rp.SIO.GPIO_OUT_SET.Set(mask)
}

// Drive low every GPIO appearing in mask
//
//go:inline
func GPIO_clr_mask(mask uint32) {
	rp.SIO.GPIO_OUT_CLR.Set(mask)
}

// Toggle every GPIO appearing in mask
//
//go:inline
func GPIO_xor_mask(mask uint32) {
	rp.SIO.GPIO_OUT_XOR.Set(mask)
}

// Drive a single GPIO high/low
//
//go:inline
func GPIO_put(pin GPIO_pin, value bool) {
	assert(pin < NUM_BANK0_GPIOS)
	if value {
		GPIO_set_mask(1 << pin)
	} else {
		GPIO_clr_mask(1 << pin)
	}
}

// Determine whether a GPIO is currently driven high or low
//
//go:inline
func GPIO_get_out_level(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return rp.SIO.GPIO_OUT.HasBits(1 << pin)
}

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - DIRECTION

// Set a number of GPIOs to output
func GPIO_set_dir_out_masked(mask uint32) {
	rp.SIO.GPIO_OE_SET.Set(mask)
}

// Set a number of GPIOs to input
func GPIO_set_dir_in_masked(mask uint32) {
	rp.SIO.GPIO_OE_CLR.Set(mask)
}

// Set multiple GPIO directions
func GPIO_set_dir_masked(mask, value uint32) {
	rp.SIO.GPIO_OUT_XOR.Set((rp.SIO.GPIO_OE.Get() ^ value) & mask)
}

// Set direction of all pins simultaneously
func GPIO_set_dir_all_bits(values uint32) {
	rp.SIO.GPIO_OE.Set(values)
}

// Set a single GPIO direction
func GPIO_set_dir(pin GPIO_pin, dir uint8) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(dir == GPIO_DIR_IN || dir == GPIO_DIR_OUT)
	switch dir {
	case GPIO_DIR_IN:
		GPIO_set_dir_in_masked(1 << pin)
	case GPIO_DIR_OUT:
		GPIO_set_dir_out_masked(1 << pin)
	default:
		assert(false)
	}
}

// Check if a specific GPIO direction is OUT
func GPIO_is_dir_out(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return rp.SIO.GPIO_OE.HasBits(1 << pin)
}

// Get a specific GPIO direction
func GPIO_get_dir(pin GPIO_pin) uint8 {
	if GPIO_is_dir_out(pin) {
		return GPIO_DIR_OUT
	}
	return GPIO_DIR_IN
}

//////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - INTERRUPTS

func gpio_set_irq_enabled(pin GPIO_pin, events GPIO_irq_level, enabled bool, base *gpio_irqctrl_t) {
	// Clear stale events which might cause immediate spurious handler entry
	gpio_acknowledge_irq(pin, events)

	target := (uint32(pin) % 8) << 2
	offset := uint32(pin) >> 3
	mask := uint32((GPIO_IRQ_LEVEL_MAX)<<1 - 1)

	// Disable interrupt
	base.inte[offset].ClearBits(mask << target)

	// Enable interrupt
	if enabled {
		base.inte[offset].SetBits(uint32(events) << target)
	}
}

//go:inline
func gpio_acknowledge_irq(pin GPIO_pin, events GPIO_irq_level) {
	target := (uint32(pin) % 8) << 2
	offset := uint32(pin) >> 3
	gpio_io_bank0.intr[offset].Set(uint32(events) << target)
}

//go:inline
func gpio_irq_status(pin GPIO_pin) GPIO_irq_level {
	offset := uint32(pin) >> 3
	gpio_io_bank0.intr[offset].Get()
	return GPIO_IRQ_LEVEL_NONE
}

//go:inline
func gpio_default_irq_handler(callback gpio_irq_callback_t, base *gpio_irqctrl_t, mask uint32) {
	fmt.Println("gpio_default_irq_handler")
	for pin := GPIO_pin(0); pin < NUM_BANK0_GPIOS; pin += 8 {
		events8 := base.ints[pin>>3]
		fmt.Printf("Pin=%v Events8=%08X\n", pin, events8)
		// note we assume events8 is 0 for non-existent GPIO
		//for i := pin; events8 && i < pin+8; i++ {
		//	events := events8 & 0x0F
		//	if events && (mask&(1<<i) == 0) {
		//		gpio_acknowledge_irq(i, events)
		//		if callback != nil {
		//			callback(i, events)
		//		}
		//	}
		//	events8 >>= 4
		//}
	}
}
