//go:build !production

package pico

import (
	"fmt"
	"strings"
)

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

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

func (v State) String() string {
	if v == StateNone {
		return v._String()
	}
	str := ""
	for f := StateNone; f <= StateMax; f <<= 1 {
		if v&f != 0 {
			str += "|" + f._String()
		}
	}
	return strings.Trim(str, "|")
}

func (v State) _String() string {
	switch v {
	case StateLow:
		return "StateLow"
	case StateHigh:
		return "StateHigh"
	case StateFall:
		return "StateFall"
	case StateRise:
		return "StateRise"
	case StateNone:
		return "StateNone"
	default:
		return fmt.Sprintf("State(0x%02X)", uint(v))
	}
}
