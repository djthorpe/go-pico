//go:build !rpi

package pico

import (
	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/sdk"
)

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	_GPIO *gpio
)

//////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func init() {
	// Initialise GPIO
	_GPIO = _NewGPIO()

	// Initialise PWM
	for slice_num := uint32(0); slice_num < NUM_PWM_SLICES; slice_num++ {
		_NewPWM(slice_num)
	}
}
