package pico

import (
	"io"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Pin uint          // GPIO Logical pin number
type EventField uint32 // Populated data fields in an event

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// UART communication
type UART interface {
	// Print
	Print(args ...interface{})
	Println(args ...interface{})
	Printf(v string, args ...interface{})
}

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
	EventSource

	High(...Pin)
	Low(...Pin)
}

// BME280 temperature, pressure and humidity sensor
type BME280 interface {
	io.Closer
	EventSource

	// One-shot Measurement, emitting an event on successful read
	Sample() error
}

// Mark an instance as a source of events
type EventSource interface{}

// Event
type Event interface {
	Source() EventSource // Emitter of event
	Is(EventField) bool  // Fields
	Bool() bool          // Pure boolean number
	Temperature() int32  // Temperature in milli-Celcius
	Pressure() int32     // Pressure in milli-Pascals
	Humidity() int32     // Relative Humidity in centi-Percent
	Altitude() int32     // Altitude in centi-Mmetres
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	Bool EventField = (1 << iota)
	Temperature
	Pressure
	Humidity
	Altitude
	Max  EventField = Altitude
	None EventField = 0
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f EventField) String() string {
	if f == None {
		return f.flag()
	}
	str := ""
	for v := EventField(1); v <= Max; v <<= 1 {
		if f&v == v {
			str += "|" + v.flag()
		}
	}
	return strings.TrimPrefix(str, "|")
}

func (f EventField) flag() string {
	switch f {
	case None:
		return "None"
	case Bool:
		return "Bool"
	case Temperature:
		return "Temperature"
	case Pressure:
		return "Pressure"
	case Humidity:
		return "Humidity"
	case Altitude:
		return "Altitude"
	default:
		return "??EventType"
	}
}
