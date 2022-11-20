//go:build !rpi

package pico

import (
	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/sdk"
)

//////////////////////////////////////////////////////////////////////////////
// TYPES

type SPI struct {
	Num  uint32
	RX   Pin
	TX   Pin
	SCK  Pin
	CS   Pin
	Baud uint32
}

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SPI_DEFAULT_BAUD_RATE = 1000000
)

//////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func _NewSPI(config SPI) *SPI {
	// Set default baud rate
	if config.Baud == 0 {
		config.Baud = SPI_DEFAULT_BAUD_RATE
	}

	// Initialise SPI
	config.Baud = SPI_init(config.Num, config.Baud)

	// Return success
	return &config
}
