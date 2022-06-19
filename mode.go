package pico

//////////////////////////////////////////////////////////////////////////////
// TYPES

type Mode uint8
type State uint8

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ModeOutput Mode = iota
	ModeInput
	ModeInputPulldown
	ModeInputPullup
	ModeUART
	ModePWM
	ModeI2C
	ModeSPI
	ModeOff
)

const (
	StateLow State = (1 << iota)
	StateHigh
	StateFall
	StateRise
	StateNone State = 0
	StateMax        = StateRise
)
