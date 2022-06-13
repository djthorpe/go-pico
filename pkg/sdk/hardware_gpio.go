//go:build rp2040

package sdk

import (
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
	ints [4]volatile.Register32 // force
	intf [4]volatile.Register32 // status
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
)

//////////////////////////////////////////////////////////////////////////////
// METHODS

//  Initialise a GPIO for (enabled I/O and set func to GPIO_FUNC_SIO)
//
func GPIO_init(pin GPIO_pin) {
	assert(pin < NUM_BANK0_GPIOS)
	rp.SIO.GPIO_OE_CLR.Set(uint32(1) << pin)
	rp.SIO.GPIO_OUT_CLR.Set(uint32(1) << pin)
	GPIO_set_function(pin, GPIO_FUNC_SIO)
}

// Resets a GPIO back to the NULL function, i.e. disables it.
//
func GPIO_deinit(pin GPIO_pin) {
	assert(pin < NUM_BANK0_GPIOS)
	GPIO_set_function(pin, GPIO_FUNC_NULL)
}

// Initialise multiple GPIOs (enabled I/O and set func to GPIO_FUNC_SIO)
//
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
//
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
//
func GPIO_get_function(pin GPIO_pin) GPIO_function {
	assert(pin < NUM_BANK0_GPIOS)
	return GPIO_function((gpio_io_bank0.gpio[pin].ctrl.Get() & rp.IO_BANK0_GPIO0_CTRL_FUNCSEL_Msk) >> rp.IO_BANK0_GPIO0_CTRL_FUNCSEL_Pos)
}

// Select up and down pulls on specific GPIO
//
func GPIO_set_pulls(pin GPIO_pin, up, down bool) {
	assert(pin < NUM_BANK0_GPIOS)
	gpio_pads_bank0.gpio[pin].ReplaceBits(bool_to_bit(up)<<rp.PADS_BANK0_GPIO0_PUE_Pos|bool_to_bit(down)<<rp.PADS_BANK0_GPIO0_PDE_Pos, rp.PADS_BANK0_GPIO0_PUE_Msk|rp.PADS_BANK0_GPIO0_PDE_Msk, 0)
}

// Set specified GPIO to be pulled up
//
func GPIO_pull_up(pin GPIO_pin) {
	GPIO_set_pulls(pin, true, false)
}

// Set specified GPIO to be pulled down
//
func GPIO_pull_down(pin GPIO_pin) {
	GPIO_set_pulls(pin, false, true)
}

// Set specified GPIO to be floating
//
func GPIO_disable_pulls(pin GPIO_pin) {
	GPIO_set_pulls(pin, false, false)
}

// Determine if the specified GPIO is pulled up
//
func GPIO_is_pulled_up(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return gpio_pads_bank0.gpio[pin].HasBits(rp.PADS_BANK0_GPIO0_PUE_Msk)
}

// Determine if the specified GPIO is pulled down
//
func GPIO_is_pulled_down(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return gpio_pads_bank0.gpio[pin].HasBits(rp.PADS_BANK0_GPIO0_PDE_Msk)
}

//  Set GPIO IRQ override
//
func GPIO_set_irqover(pin GPIO_pin, value GPIO_override) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(value <= GPIO_OVERRIDE_HIGH)
	gpio_io_bank0.gpio[pin].ctrl.ReplaceBits(uint32(value)<<rp.IO_BANK0_GPIO0_CTRL_IRQOVER_Pos, rp.IO_BANK0_GPIO0_CTRL_IRQOVER_Msk, 0)
}

// Set GPIO output override
//
func GPIO_set_outover(pin GPIO_pin, value GPIO_override) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(value <= GPIO_OVERRIDE_HIGH)
	gpio_io_bank0.gpio[pin].ctrl.ReplaceBits(uint32(value)<<rp.IO_BANK0_GPIO0_CTRL_OUTOVER_Pos, rp.IO_BANK0_GPIO0_CTRL_OUTOVER_Msk, 0)
}

// Set GPIO input override
//
func GPIO_set_inover(pin GPIO_pin, value GPIO_override) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(value <= GPIO_OVERRIDE_HIGH)
	gpio_io_bank0.gpio[pin].ctrl.ReplaceBits(uint32(value)<<rp.IO_BANK0_GPIO0_CTRL_INOVER_Pos, rp.IO_BANK0_GPIO0_CTRL_INOVER_Msk, 0)
}

// Set GPIO output enable override
//
func GPIO_set_oeover(pin GPIO_pin, value GPIO_override) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(value <= GPIO_OVERRIDE_HIGH)
	gpio_io_bank0.gpio[pin].ctrl.ReplaceBits(uint32(value)<<rp.IO_BANK0_GPIO0_CTRL_OEOVER_Pos, rp.IO_BANK0_GPIO0_CTRL_OEOVER_Msk, 0)
}

// Enable GPIO input
//
func GPIO_set_input_enabled(pin GPIO_pin, enabled bool) {
	assert(pin < NUM_BANK0_GPIOS)
	gpio_pads_bank0.gpio[pin].ReplaceBits(bool_to_bit(enabled)<<rp.PADS_BANK0_GPIO0_IE_Pos, rp.PADS_BANK0_GPIO0_IE_Msk, 0)
}

// Set or clear GPIO output enabled
//
func GPIO_set_output_enabled(pin GPIO_pin, enabled bool) {
	assert(pin < NUM_BANK0_GPIOS)
	if enabled {
		rp.SIO.GPIO_OE_SET.Set(1 << pin)
	} else {
		rp.SIO.GPIO_OE_CLR.Set(1 << pin)
	}
}

// Get GPIO output enabled state
//
func GPIO_get_output_enabled(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return rp.SIO.GPIO_OE.HasBits(1 << pin)
}

// Enable/disable GPIO input hysteresis (Schmitt trigger)
//
func GPIO_set_input_hysteresis_enabled(pin GPIO_pin, enabled bool) {
	assert(pin < NUM_BANK0_GPIOS)
	gpio_pads_bank0.gpio[pin].ReplaceBits(bool_to_bit(enabled)<<rp.PADS_BANK0_GPIO0_SCHMITT_Pos, rp.PADS_BANK0_GPIO0_SCHMITT_Msk, 0)
}

// Determine whether input hysteresis is enabled on a specified GPIO
//
func GPIO_is_input_hysteresis_enabled(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return gpio_pads_bank0.gpio[pin].HasBits(rp.PADS_BANK0_GPIO0_SCHMITT_Msk)
}

// Set slew rate for a specified GPIO
//
func GPIO_set_slew_rate(pin GPIO_pin, slew GPIO_slew_rate) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(slew <= GPIO_SLEW_RATE_FAST)
	gpio_pads_bank0.gpio[pin].ReplaceBits(uint32(slew)<<rp.PADS_BANK0_GPIO0_SLEWFAST_Pos, rp.PADS_BANK0_GPIO0_SLEWFAST_Msk, 0)
}

// Determine current slew rate for a specified GPIO
//
func GPIO_get_slew_rate(pin GPIO_pin) GPIO_slew_rate {
	assert(pin < NUM_BANK0_GPIOS)
	return GPIO_slew_rate((gpio_pads_bank0.gpio[pin].Get() & rp.PADS_BANK0_GPIO0_SLEWFAST) >> rp.PADS_BANK0_GPIO0_SLEWFAST_Pos)
}

// Set drive strength for a specified GPIO
//
func GPIO_set_drive_strength(pin GPIO_pin, drive GPIO_drive_strength) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(drive <= GPIO_DRIVE_STRENGTH_12MA)
	gpio_pads_bank0.gpio[pin].ReplaceBits(uint32(drive)<<rp.PADS_BANK0_GPIO0_DRIVE_Pos, rp.PADS_BANK0_GPIO0_DRIVE_Msk, 0)
}

// Determine current slew rate for a specified GPIO
//
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
		gpio_set_irq_enabled(pin, events, enabled, gpio_io_bank0.proc0IRQctrl)
	case 1:
		gpio_set_irq_enabled(pin, events, enabled, gpio_io_bank0.proc1IRQctrl)
	default:
		assert(false)
	}
}

func GPIO_acknowledge_irq(pin GPIO_pin, events GPIO_irq_level) {
	assert(pin < NUM_BANK0_GPIOS)
	assert(events <= (GPIO_IRQ_LEVEL_MAX<<1)-1)
	target := (pin % 8) * 4
	gpio_io_bank0.intr[pin/8].Set(uint32(events) << target)
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
//
func GPIO_set_dir_out_masked(mask uint32) {
	rp.SIO.GPIO_OE_SET.Set(mask)
}

// Set a number of GPIOs to input
//
func GPIO_set_dir_in_masked(mask uint32) {
	rp.SIO.GPIO_OE_CLR.Set(mask)
}

// Set multiple GPIO directions
//
func GPIO_set_dir_masked(mask, value uint32) {
	rp.SIO.GPIO_OUT_XOR.Set((rp.SIO.GPIO_OE.Get() ^ value) & mask)
}

// Set direction of all pins simultaneously
//
func GPIO_set_dir_all_bits(values uint32) {
	rp.SIO.GPIO_OE.Set(values)
}

// Set a single GPIO direction
//
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
//
func GPIO_is_dir_out(pin GPIO_pin) bool {
	assert(pin < NUM_BANK0_GPIOS)
	return rp.SIO.GPIO_OE.HasBits(1 << pin)
}

// Get a specific GPIO direction
//
func GPIO_get_dir(pin GPIO_pin) uint8 {
	if GPIO_is_dir_out(pin) {
		return GPIO_DIR_OUT
	}
	return GPIO_DIR_IN
}

//////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - INTERRUPTS

func gpio_set_irq_enabled(pin GPIO_pin, events GPIO_irq_level, enabled bool, base gpio_irqctrl_t) {
	// Clear stale events which might cause immediate spurious handler entry
	GPIO_acknowledge_irq(pin, events)

	// Enable or disable interrupt
	target := (pin % 8) * 4
	base.inte[pin/8].ClearBits(uint32(events) << target)
	if enabled {
		base.inte[pin/8].SetBits(uint32(events) << target)
	}
}

/*
// Callbacks to be called for pins configured with SetInterrupt.
var (
	pinCallbacks [2]func(Pin)
	setInt       [2]bool
)

// SetInterrupt sets an interrupt to be executed when a particular pin changes
// state. The pin should already be configured as an input, including a pull up
// or down if no external pull is provided.
//
// This call will replace a previously set callback on this pin. You can pass a
// nil func to unset the pin change interrupt. If you do so, the change
// parameter is ignored and can be set to any value (such as 0).
func (p Pin) SetInterrupt(change PinChange, callback func(Pin)) error {
	if p > 31 || p < 0 {
		return ErrInvalidInputPin
	}
	core := CurrentCore()
	if callback == nil {
		// disable current interrupt
		p.setInterrupt(change, false)
		pinCallbacks[core] = nil
		return nil
	}

	if pinCallbacks[core] != nil {
		// Callback already configured. Should disable callback by passing a nil callback first.
		return ErrNoPinChangeChannel
	}
	p.setInterrupt(change, true)
	pinCallbacks[core] = callback

	if setInt[core] {
		// interrupt has already been set. Exit.
		println("core set")
		return nil
	}
	interrupt.New(rp.IRQ_IO_IRQ_BANK0, gpioHandleInterrupt).Enable()
	irqSet(rp.IRQ_IO_IRQ_BANK0, true)
	return nil
}

// gpioHandleInterrupt finds the corresponding pin for the interrupt.
// C SDK equivalent of gpio_irq_handler
func gpioHandleInterrupt(intr interrupt.Interrupt) {
	// panic("END") // if program is not ended here rp2040 will call interrupt again when finished, a vicious spin cycle.
	core := CurrentCore()
	callback := pinCallbacks[core]
	if callback != nil {
		// TODO fix gpio acquisition (see below)
		// For now all callbacks get pin 255 (nonexistent).
		callback(0xff)
	}
	var gpio Pin
	for gpio = 0; gpio < _NUMBANK0_GPIOS; gpio++ {
		// Acknowledge all GPIO interrupts for now
		// since we are yet unable to acquire interrupt status
		gpio.acknowledgeInterrupt(0xff) // TODO fix status get. For now we acknowledge all pending interrupts.
		// Commented code below from C SDK not working.
		// statreg := base.intS[gpio>>3]
		// change := getIntChange(gpio, statreg.Get())
		// if change != 0 {
		// 	gpio.acknowledgeInterrupt(change)
		// 	if callback != nil {
		// 		callback(gpio)
		// 		return
		// 	} else {
		// 		panic("unset callback in handler")
		// 	}
		// }
	}
}

// events returns the bit representation of the pin change for the rp2040.
func (change PinChange) events() uint32 {
	return uint32(change)
}

// Acquire interrupt data from a INT status register.
func getIntChange(p Pin, status uint32) PinChange {
	return PinChange(status>>(4*(p%8))) & 0xf
}

func GPIO_set_irq_enabled_with_callback(pin GPIO_pin, events uint32, enabled bool, callback gpio_irq_callback_t) {
	assert(pin < NUM_BANK0_GPIOS)

	core := get_core_num()
	// TODO: p.setInterrupt(change, false)

	// disable current interrupt
	gpio_irq_callback[core] = nil
	if callback == nil {
		return
	}

	// enable current interrupt
	p.setInterrupt(change, true)
	pinCallbacks[core] = callback

	if setInt[core] {
		// interrupt has already been set. Exit.
		println("core set")
		return nil
	}
	interrupt.New(rp.IRQ_IO_IRQ_BANK0, gpioHandleInterrupt).Enable()
	irqSet(rp.IRQ_IO_IRQ_BANK0, true)
	return nil
	// TODO
}

// Clears interrupt flag on a pin
func (p Pin) acknowledgeInterrupt(change PinChange) {
	ioBank0.intR[p>>3].Set(p.ioIntBit(change))
}

// Basic interrupt setting via ioBANK0 for GPIO interrupts.
func (p Pin) setInterrupt(change PinChange, enabled bool) {
	// Separate mask/force/status per-core, so check which core called, and
	// set the relevant IRQ controls.
	switch CurrentCore() {
	case 0:
		p.ctrlSetInterrupt(change, enabled, &ioBank0.proc0IRQctrl)
	case 1:
		p.ctrlSetInterrupt(change, enabled, &ioBank0.proc1IRQctrl)
	}
}

// ctrlSetInterrupt acknowledges any pending interrupt and enables or disables
// the interrupt for a given IRQ control bank (IOBANK, DormantIRQ, QSPI).
//
// pico-sdk calls this the _gpio_set_irq_enabled, not to be confused with
// gpio_set_irq_enabled (no leading underscore).
func (p Pin) ctrlSetInterrupt(change PinChange, enabled bool, base *irqCtrl) {
	p.acknowledgeInterrupt(change)
	enReg := &base.intE[p>>3]
	if enabled {
		enReg.SetBits(p.ioIntBit(change))
	} else {
		enReg.ClearBits(p.ioIntBit(change))
	}
}

// Enable or disable a specific interrupt on the executing core.
// num is the interrupt number which must be in [0,31].
func irqSet(num uint32, enabled bool) {
	if num >= _NUMIRQ {
		return
	}
	irqSetMask(1<<num, enabled)
}

func irqSetMask(mask uint32, enabled bool) {
	if enabled {
		// Clear pending before enable
		// (if IRQ is actually asserted, it will immediately re-pend)
		rp.PPB.NVIC_ICPR.Set(mask)
		rp.PPB.NVIC_ISER.Set(mask)
	} else {
		rp.PPB.NVIC_ICER.Set(mask)
	}
}

*/
