//go:build !production

package sdk

import (
	"fmt"
	"strings"
)

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

func (v GPIO_irq_level) _String() string {
	switch v {
	case GPIO_IRQ_LEVEL_NONE:
		return "GPIO_IRQ_LEVEL_NONE"
	case GPIO_IRQ_LEVEL_LOW:
		return "GPIO_IRQ_LEVEL_LOW"
	case GPIO_IRQ_LEVEL_HIGH:
		return "GPIO_IRQ_LEVEL_HIGH"
	case GPIO_IRQ_EDGE_FALL:
		return "GPIO_IRQ_EDGE_FALL"
	case GPIO_IRQ_EDGE_RISE:
		return "GPIO_IRQ_EDGE_RISE"
	default:
		return fmt.Sprintf("GPIO_irq_level(0x%02X)", uint(v))
	}
}

func (v GPIO_irq_level) String() string {
	if v == GPIO_IRQ_LEVEL_NONE {
		return v._String()
	}
	str := ""
	for f := GPIO_irq_level(1); f <= GPIO_IRQ_LEVEL_MAX; f <<= 1 {
		if v&f != 0 {
			str += f._String() + "|"
		}
	}
	return strings.Trim(str, "|")
}
