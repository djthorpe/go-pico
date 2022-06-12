//go:build rp2040

package sdk

import (
	// Module imports
	rp "device/rp"
)

// SDK documentation
// https://github.com/raspberrypi/pico-sdk/tree/master/src/rp2_common/hardware_spi

//////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	SPI_cpha_t  uint32
	SPI_cpol_t  uint32
	SPI_order_t uint32
)

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SPI_CPHA_0 SPI_cpha_t = 0
	SPI_CPHA_1 SPI_cpha_t = 1
)

const (
	SPI_CPOL_0 SPI_cpol_t = 0
	SPI_CPOL_1 SPI_cpol_t = 1
)

const (
	SPI_LSB_FIRST SPI_order_t = 0
	SPI_MSB_FIRST SPI_order_t = 1
)

const (
	_SPI_FIFO_DEPTH = 8
)

var (
	spi_groups = [NUM_SPIS]*rp.SPI0_Type{rp.SPI0, rp.SPI1}
)

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Reset SPI
//
func SPI_reset(spi uint32) {
	assert(spi < NUM_SPIS)
	switch spi {
	case 0:
		reset_block(rp.RESETS_RESET_SPI0_Msk)
	case 1:
		reset_block(rp.RESETS_RESET_SPI1_Msk)
	}

}

// Unreset SPI
//
func SPI_unreset(spi uint32) {
	assert(spi < NUM_SPIS)
	switch spi {
	case 0:
		unreset_block(rp.RESETS_RESET_SPI0_Msk)
	case 1:
		unreset_block(rp.RESETS_RESET_SPI1_Msk)
	}
}

// Initialise SPI instances
//
func SPI_init(spi, baudrate uint32) {
	assert(spi < NUM_SPIS)
	SPI_reset(spi)
	SPI_unreset(spi)
	// TODO SPI_set_baudrate(spi, baudrate)

	// Always enable DREQ signals -- harmless if DMA is not listening
	SPI_set_format(spi, 8, SPI_CPOL_0, SPI_CPHA_0, SPI_MSB_FIRST)
	spi_groups[spi].SSPDMACR.SetBits(rp.SPI0_SSPDMACR_TXDMAE | rp.SPI0_SSPDMACR_RXDMAE)
	SPI_set_format(spi, 8, SPI_CPOL_0, SPI_CPHA_0, SPI_MSB_FIRST)

	// Finally enable the SPI
	spi_groups[spi].SSPCR1.SetBits(rp.SPI0_SSPCR1_SSE)
}

// Deinitialise SPI instances
//
func SPI_deinit(spi uint32) {
	assert(spi < NUM_SPIS)
	spi_groups[spi].SSPCR1.ClearBits(rp.SPI0_SSPCR1_SSE)
	spi_groups[spi].SSPDMACR.ClearBits(rp.SPI0_SSPDMACR_TXDMAE | rp.SPI0_SSPDMACR_RXDMAE)
	SPI_reset(spi)
}

/*
// Set SPI baudrate
//
func SPI_set_baudrate(spi, baudrate uint32) uint32 {
	assert(spi < NUM_SPIS)
	assert(baudrate <= clock_get_hz(CLOCK_PERI))

	uint32 prescale, postdiv

    // Find smallest prescale value which puts output frequency in range of
    // post-divide. Prescale is an even number from 2 to 254 inclusive.
    for (prescale = 2; prescale <= 254; prescale += 2) {
        if (freq_in < (prescale + 2) * 256 * (uint64_t) baudrate)
            break;
    }
    invalid_params_if(SPI, prescale > 254); // Frequency too low

    // Find largest post-divide which makes output <= baudrate. Post-divide is
    // an integer in the range 1 to 256 inclusive.
    for (postdiv = 256; postdiv > 1; --postdiv) {
        if (freq_in / (prescale * (postdiv - 1)) > baudrate)
            break;
    }

    spi_get_hw(spi)->cpsr = prescale;
    hw_write_masked(&spi_get_hw(spi)->cr0, (postdiv - 1) << SPI_SSPCR0_SCR_LSB, SPI_SSPCR0_SCR_BITS);

    // Return the frequency we were able to achieve
    return freq_in / (prescale * postdiv);
}
*/

// Configure SPI
//
// Must be SPI_MSB_FIRST, no other values supported on the PL022
//
func SPI_set_format(spi uint32, data_bits uint8, cpol SPI_cpol_t, cpha SPI_cpha_t, order SPI_order_t) {
	assert(spi < NUM_SPIS)
	assert(data_bits >= 4 && data_bits <= 16)
	assert(cpol == SPI_CPOL_0 || cpol == SPI_CPOL_1)
	assert(cpha == SPI_CPHA_0 || cpha == SPI_CPHA_1)
	assert(order == SPI_MSB_FIRST)
	v := uint32(uint32(data_bits-1)<<rp.SPI0_SSPCR0_DSS_Pos | uint32(cpol)<<rp.SPI0_SSPCR0_SPO_Pos | uint32(cpha)<<rp.SPI0_SSPCR0_SPH_Pos)
	m := uint32(rp.SPI0_SSPCR0_DSS_Msk | rp.SPI0_SSPCR0_SPO_Msk | rp.SPI0_SSPCR0_SPH_Msk)
	spi_groups[spi].SSPCR0.ReplaceBits(v, m, 0)
}

// Set SPI master/slave
//
// By default, spi_init() sets master-mode
//
func SPI_set_slave(spi uint32, slave bool) {
	assert(spi < NUM_SPIS)
	if slave {
		spi_groups[spi].SSPCR1.SetBits(rp.SPI0_SSPCR1_MS)
	} else {
		spi_groups[spi].SSPCR1.ClearBits(rp.SPI0_SSPCR1_MS)
	}
}

// Check whether a write can be done on SPI device
//
// Although the controllers each have a 8 deep TX FIFO, the current HW
// implementation can only return 0 or 1 rather than the space available.
//
func SPI_is_writable(spi uint32) uint32 {
	assert(spi < NUM_SPIS)
	return (spi_groups[spi].SSPSR.Get() & rp.SPI0_SSPSR_TNF_Msk) >> rp.SPI0_SSPSR_TNF_Pos
}

// Check whether a read can be done on SPI device
//
// Although the controllers each have a 8 deep RX FIFO,
// the current HW implementation can only return 0 or 1
//
func SPI_is_readable(spi uint32) uint32 {
	assert(spi < NUM_SPIS)
	return (spi_groups[spi].SSPSR.Get() & rp.SPI0_SSPSR_RNE_Msk) >> rp.SPI0_SSPSR_RNE_Pos
}

// Write/Read to/from an SPI device
//
func SPI_write_read_blocking(spi uint32, w, r []uint8) {
	assert(spi < NUM_SPIS)
	assert(len(w) > 0)

	// Never have more transfers in flight than will fit into the RX FIFO,
	// else FIFO will overflow if this code is heavily interrupted
	rx_remaining := uint32(len(r))
	tx_remaining := uint32(len(w))
	for rx_remaining > 0 && tx_remaining > 0 {
		if tx_remaining > 0 && bits_to_bool(SPI_is_writable(spi)) && rx_remaining-tx_remaining < _SPI_FIFO_DEPTH {
			spi_groups[spi].SSPDR.Set(w[0])
			tx_remaining--
		}
		if rx_remaining > 0 && bits_to_bool(SPI_is_readable(spi)) {
			r[0] = spi_groups[spi].SSPDR.Get()
			rx_remaining--
		}
	}
}

/*

// Write to an SPI device, blocking
//
func SPI_write_blocking(spi_inst_t *spi, w []uint8) uint32

// Read from an SPI device
//
// Blocks until all data is transferred. No timeout, as SPI hardware always transfers at a known data rate.
// repeated_tx_data is output repeatedly on TX as data is read in from RX
func SPI_read_blocking(spi uint32, repeated_tx_data uint8, r []uint8)

// Write/Read half words to/from an SPI device
//
func spi_write16_read16_blocking(spi uint32, w, r []uint16)

// Write half words to an SPI device
//
func SPI_write16_blocking(spi uint32, w []uint16) uint32

// Read half words from an SPI device
//
// Blocks until all data is transferred. No timeout, as SPI hardware always transfers at
// a known data rate. repeated_tx_data is output repeatedly on TX as data is read in from RX
//
func SPI_read16_blocking(spi uint32, repeated_tx_data uint16, w []uint16) uint32
*/

/*
	SSPCR0       volatile.Register32 // 0x0
	SSPCR1       volatile.Register32 // 0x4
	SSPDR        volatile.Register32 // 0x8
	SSPSR        volatile.Register32 // 0xC
	SSPCPSR      volatile.Register32 // 0x10
	SSPIMSC      volatile.Register32 // 0x14
	SSPRIS       volatile.Register32 // 0x18
	SSPMIS       volatile.Register32 // 0x1C
	SSPICR       volatile.Register32 // 0x20
	SSPDMACR     volatile.Register32 // 0x24
*/
