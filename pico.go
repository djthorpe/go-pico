package pico

type I2C interface {
	// Read
	ReadRegister_Uint8(slave, register uint8) (uint8, error)
	ReadRegister_Uint16(slave, register uint8) (uint16, error)
	ReadRegister(slave, register uint8, data []uint8) error

	// Write
	WriteRegister_Uint8(slave, register, data uint8) error
}
