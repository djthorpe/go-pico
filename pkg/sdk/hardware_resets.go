//go:build rp2040

package sdk

import (
	// Module imports
	rp "device/rp"
)

// SDK documentation
// https://github.com/raspberrypi/pico-sdk/tree/master/src/rp2_common/hardware_resets
//
//
// Hardware Reset API

// Reset the specified HW blocks
//
func reset_block(bits uint32) {
	rp.RESETS.RESET.SetBits(bits)
}

// Bring specified HW blocks out of reset
//
func unreset_block(bits uint32) {
	rp.RESETS.RESET.ClearBits(bits)
}

// Bring specified HW blocks out of reset and wait for completion
//
func unreset_block_wait(bits uint32) {
	unreset_block(bits)
	for {
		if rp.RESETS.RESET_DONE.Get()&bits == bits {
			break
		}
	}
}
