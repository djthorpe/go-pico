package pico

import "io"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Pin uint // GPIO Logical pin number

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// I2C interface
type I2C interface {
	io.Closer

	// Read
	ReadRegister(slave, register uint8, data []uint8) error
	ReadRegister_Uint8(slave, register uint8) (uint8, error)

	// Write
	WriteRegister_Uint8(slave, register, data uint8) error
}

// SPI interface
type SPI interface {
	io.Closer

	// Transfer data on SPI bus
	Transfer(w, r []uint8) error
}

// GPIO interface
type GPIO interface {
	io.Closer

	High(...Pin)
	Low(...Pin)
}

// UART communication
type UART interface {
	// Print
	Print(args ...interface{})
	Println(args ...interface{})
	Printf(v string, args ...interface{})
}

// BME280 temperature, pressure and humidity sensor
type BME280 interface {
	io.Closer

	// One-shot Measurement, emitting an event on successful read
	Sample() error
}
