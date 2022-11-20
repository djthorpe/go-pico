package sdk

import (
	"unsafe"

	// Module imports
	rp "device/rp"
	volatile "runtime/volatile"
)

// SDK documentation
// https://github.com/raspberrypi/pico-sdk/blob/master/src/rp2_common/hardware_i2c

//////////////////////////////////////////////////////////////////////////////
// TYPES

type i2c_t struct {
	con                volatile.Register32 // 0x0
	tar                volatile.Register32 // 0x4
	sar                volatile.Register32 // 0x8
	_                  [4]byte
	data_cmd           volatile.Register32 // 0x10
	ss_scl_hcnt        volatile.Register32 // 0x14
	ss_scl_lcnt        volatile.Register32 // 0x18
	fs_scl_hcnt        volatile.Register32 // 0x1C
	fs_scl_lcnt        volatile.Register32 // 0x20
	_                  [8]byte
	intr_stat          volatile.Register32 // 0x2C
	intr_mask          volatile.Register32 // 0x30
	raw_intr_stat      volatile.Register32 // 0x34
	rx_tl              volatile.Register32 // 0x38
	tx_tl              volatile.Register32 // 0x3C
	CLR_INTR           volatile.Register32 // 0x40
	CLR_RX_UNDER       volatile.Register32 // 0x44
	CLR_RX_OVER        volatile.Register32 // 0x48
	CLR_TX_OVER        volatile.Register32 // 0x4C
	CLR_RD_REQ         volatile.Register32 // 0x50
	CLR_TX_ABRT        volatile.Register32 // 0x54
	CLR_RX_DONE        volatile.Register32 // 0x58
	CLR_ACTIVITY       volatile.Register32 // 0x5C
	CLR_STOP_DET       volatile.Register32 // 0x60
	CLR_START_DET      volatile.Register32 // 0x64
	CLR_GEN_CALL       volatile.Register32 // 0x68
	enable             volatile.Register32 // 0x6C
	status             volatile.Register32 // 0x70
	TXFLR              volatile.Register32 // 0x74
	RXFLR              volatile.Register32 // 0x78
	SDA_HOLD           volatile.Register32 // 0x7C
	TX_ABRT_SOURCE     volatile.Register32 // 0x80
	SLV_DATA_NACK_ONLY volatile.Register32 // 0x84
	dma_cr             volatile.Register32 // 0x88
	DMA_TDLR           volatile.Register32 // 0x8C
	DMA_RDLR           volatile.Register32 // 0x90
	SDA_SETUP          volatile.Register32 // 0x94
	ACK_GENERAL_CALL   volatile.Register32 // 0x98
	ENABLE_STATUS      volatile.Register32 // 0x9C
	FS_SPKLEN          volatile.Register32 // 0xA0
	_                  [4]byte
	CLR_RESTART_DET    volatile.Register32 // 0xA8
	_                  [72]byte
	COMP_PARAM_1       volatile.Register32 // 0xF4
	COMP_VERSION       volatile.Register32 // 0xF8
	COMP_TYPE          volatile.Register32 // 0xFC
}

type i2c_groups_t [NUM_I2CS]i2c_t

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	i2c_groups = [NUM_I2CS](*i2c_t){(*i2c_t)(unsafe.Pointer(rp.I2C0)), (*i2c_t)(unsafe.Pointer(rp.I2C1))}
)

//////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Initialise the I2C HW block
//
func I2C_init(inst, baudrate uint32) uint32 {
	assert(inst < NUM_I2CS)
	I2C_reset(inst)
	I2C_unreset(inst)

	// Disable before config
	i2c_groups[inst].enable.Set(0)

	// Configure as a fast-mode master with RepStart support, 7-bit addresses
	i2c_groups[inst].con.Set(
		rp.I2C0_IC_CON_SPEED_FAST<<rp.I2C0_IC_CON_SPEED_Pos |
			rp.I2C0_IC_CON_MASTER_MODE_ENABLED<<rp.I2C0_IC_CON_MASTER_MODE_Pos |
			rp.I2C0_IC_CON_IC_SLAVE_DISABLE_SLAVE_DISABLED<<rp.I2C0_IC_CON_IC_SLAVE_DISABLE_Pos |
			rp.I2C0_IC_CON_IC_RESTART_EN_ENABLED<<rp.I2C0_IC_CON_IC_RESTART_EN_Pos |
			rp.I2C0_IC_CON_TX_EMPTY_CTRL_ENABLED<<rp.I2C0_IC_CON_TX_EMPTY_CTRL_Pos)

	// Set FIFO watermarks to 1 to make things simpler. This is encoded by a register value of 0.
	i2c_groups[inst].tx_tl.Set(0)
	i2c_groups[inst].rx_tl.Set(0)

	// Always enable the DREQ signalling -- harmless if DMA isn't listening
	i2c_groups[inst].dma_cr.Set(
		rp.I2C0_IC_DMA_CR_TDMAE_ENABLED<<rp.I2C0_IC_DMA_CR_TDMAE_Pos |
			rp.I2C0_IC_DMA_CR_RDMAE_ENABLED<<rp.I2C0_IC_DMA_CR_RDMAE_Pos)

	// Re-sets i2c->hw->enable upon returning
	//TODO
	// return I2C_set_baudrate(inst, baudrate)
	return 0
}

// Disable the I2C HW block
//
func I2C_deinit(inst uint32) {
	assert(inst < NUM_I2CS)
	I2C_reset(inst)
}

// Reset I2C block
//
func I2C_reset(inst uint32) {
	assert(inst < NUM_I2CS)
	switch inst {
	case 0:
		reset_block(rp.RESETS_RESET_I2C0)
	case 1:
		reset_block(rp.RESETS_RESET_I2C1)
	}
}

// Unreset I2C block
//
func I2C_unreset(inst uint32) {
	assert(inst < NUM_I2CS)
	switch inst {
	case 0:
		unreset_block_wait(rp.RESETS_RESET_I2C0)
	case 1:
		unreset_block_wait(rp.RESETS_RESET_I2C1)
	}
}

// Addresses of the form 000 0xxx or 111 1xxx are reserved. No slave should
// have these addresses.
//
func I2C_reserved_addr(addr uint8) bool {
	return (addr&0x78) == 0 || (addr&0x78) == 0x78
}

/*
// Determine the I2C instance that is attached to the specified GPIO
// pins for SDA and SCL
//
//go:inline
func I2C_gpio_to_inst(sda, scl GPIO_pin) (uint32, bool) {
	// TODO
	return 0, false
}

// Set I2C baudrate
//
func I2C_set_baudrate(inst, baudrate uint32) uint32 {
	// TODO
}

// Set I2C port to slave mode
//
func I2C_set_slave_mode(inst uint32, slave bool, addr uint8) {

}

// Attempt to write specified number of bytes to address, blocking until the
// specified absolute time is reached.
//
func i2c_write_blocking_until(inst uint32, addr uint8, w []byte, nostop bool, until time.Duration) {
	// TODO
}

// Attempt to read specified number of bytes from address, blocking until the
// specified absolute time is reached.
func i2c_read_blocking_until(inst uint32, addr uint8, r []byte, nostop bool, until time.Duration) {
	// TODO
}
*/
