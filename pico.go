package pico

type I2C interface {
	// Read
	ReadRegister_Uint8(slave, register uint8) (uint8, error)
	ReadRegister(slave, register uint8, data []uint8) error

	// Write
	WriteRegister_Uint8(slave, register, data uint8) error
}

type SPI interface {
	// Transfer data on SPI bus
	Transfer(w, r []uint8) error
}

type UART interface {
	// Print
	Print(args ...interface{})
	Println(args ...interface{})
	Printf(v string, args ...interface{})
}

type BME280 interface {
	// One-shot Measurement
	Sample() error
}
