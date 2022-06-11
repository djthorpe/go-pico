//go:build rp2040 && !production

package sdk

import "fmt"

//////////////////////////////////////////////////////////////////////////////
// ASSERT

//go:inline
func assert(cond bool) {
	if !cond {
		panic("assertation failed")
	}
}

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v GPIO_pin) String() string {
	return fmt.Sprint("GP", uint(v))
}

func (v GPIO_function) String() string {
	switch v {
	case GPIO_FUNC_XIP:
		return "GPIO_FUNC_XIP"
	case GPIO_FUNC_SPI:
		return "GPIO_FUNC_SPI"
	case GPIO_FUNC_UART:
		return "GPIO_FUNC_UART"
	case GPIO_FUNC_I2C:
		return "GPIO_FUNC_I2C"
	case GPIO_FUNC_PWM:
		return "GPIO_FUNC_PWM"
	case GPIO_FUNC_SIO:
		return "GPIO_FUNC_SIO"
	case GPIO_FUNC_PIO0:
		return "GPIO_FUNC_PIO0"
	case GPIO_FUNC_PIO1:
		return "GPIO_FUNC_PIO1"
	case GPIO_FUNC_GPCK:
		return "GPIO_FUNC_GPCK"
	case GPIO_FUNC_USB:
		return "GPIO_FUNC_USB"
	case GPIO_FUNC_NULL:
		return "GPIO_FUNC_NULL"
	default:
		return fmt.Sprintf("GPIO_function(0x%02X)", uint(v))
	}
}

func (v GPIO_override) String() string {
	switch v {
	case GPIO_OVERRIDE_NORMAL:
		return "GPIO_OVERRIDE_NORMAL"
	case GPIO_OVERRIDE_INVERT:
		return "GPIO_OVERRIDE_INVERT"
	case GPIO_OVERRIDE_LOW:
		return "GPIO_OVERRIDE_LOW"
	case GPIO_OVERRIDE_HIGH:
		return "GPIO_OVERRIDE_HIGH"
	default:
		return fmt.Sprintf("GPIO_override(0x%02X)", uint(v))
	}
}

func (v GPIO_slew_rate) String() string {
	switch v {
	case GPIO_SLEW_RATE_SLOW:
		return "GPIO_SLEW_RATE_SLOW"
	case GPIO_SLEW_RATE_FAST:
		return "GPIO_SLEW_RATE_FAST"
	default:
		return fmt.Sprintf("GPIO_slew_rate(0x%02X)", uint(v))
	}
}

func (v GPIO_drive_strength) String() string {
	switch v {
	case GPIO_DRIVE_STRENGTH_2MA:
		return "GPIO_DRIVE_STRENGTH_2MA"
	case GPIO_DRIVE_STRENGTH_4MA:
		return "GPIO_DRIVE_STRENGTH_4MA"
	case GPIO_DRIVE_STRENGTH_8MA:
		return "GPIO_DRIVE_STRENGTH_8MA"
	case GPIO_DRIVE_STRENGTH_12MA:
		return "GPIO_DRIVE_STRENGTH_12MA"
	default:
		return fmt.Sprintf("GPIO_drive_strength(0x%02X)", uint(v))
	}
}
