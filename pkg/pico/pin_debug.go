//go:build rp2040 && !production

package pico

import (
	"fmt"
)

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v Pin) String() string {
	return fmt.Sprint("GP", uint(v))
}

func (v Mode) String() string {
	switch v {
	case ModeOutput:
		return "ModeOutput"
	case ModeInput:
		return "ModeInput"
	case ModeInputPulldown:
		return "ModeInputPulldown"
	case ModeInputPullup:
		return "ModeInputPullup"
	case ModeUART:
		return "ModeUART"
	case ModePWM:
		return "ModePWM"
	case ModeI2C:
		return "ModeI2C"
	case ModeSPI:
		return "ModeSPI"
	case ModeOff:
		return "ModeOff"
	default:
		return fmt.Sprintf("Mode(0x%02X)", uint(v))
	}
}
