//go:build rp2040

package sdk

import rp "device/rp"

const (
	NUM_BANK0_GPIOS = 30
	NUM_CORES       = 2
)

// Return the core number the call was made from
//
//go:inline
func get_core_num() uint32 {
	return rp.SIO.CPUID.Get()
}
