//go:build rp2040

package sdk

import rp "device/rp"

const (
	NUM_CORES              = 2
	NUM_DMA_CHANNELS       = 12
	NUM_DMA_TIMERS         = 4
	NUM_IRQS               = 32
	NUM_PIOS               = 2
	NUM_PIO_STATE_MACHINES = 4
	NUM_PWM_SLICES         = 8
	NUM_SPIN_LOCKS         = 32
	NUM_UARTS              = 2
	NUM_I2CS               = 2
	NUM_SPIS               = 2
	NUM_TIMERS             = 4
	NUM_ADC_CHANNELS       = 5
	NUM_BANK0_GPIOS        = 30
	NUM_QSPI_GPIOS         = 6
	PIO_INSTRUCTION_COUNT  = 32
	XOSC_MHZ               = 12
)

// Return the core number the call was made from
//
//go:inline
func get_core_num() uint32 {
	return rp.SIO.CPUID.Get()
}
