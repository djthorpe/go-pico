package pico

//////////////////////////////////////////////////////////////////////////////
// TYPES

type Mode uint

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
