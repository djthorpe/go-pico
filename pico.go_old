package pico

import (
	"io"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Pin uint          // GPIO Logical pin number
type EventField uint32 // Populated data fields in an event
type EventUnit uint32  // Populated data fields in an event

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
	Get(Pin) bool
	Set(Pin, bool)
}

// ADC interface
type ADC interface {
	io.Closer
	EventSource

	// One-shot measurement, emitting an event on successful read
	Sample() error
}

// Display interface
type Display interface {
	io.Closer
	EventSource

	// Get current backlight state
	Backlight() bool

	// Switch backlight on or off
	SetBacklight(bool)
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
	Source() EventSource                          // Emitter of event
	Is(EventField) bool                           // Field name
	Set(EventField, EventUnit, interface{}) Event // Set field, unit and value
	Emit(C chan<- Event) error                    // Emit event
	Value(EventField) (interface{}, EventUnit)    // Get untyped value or nil
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	Temperature EventField = (1 << iota)
	Pressure
	Humidity
	Altitude
	Battery
	Charging
	Discharging
	Button
	Sample
	FieldMax  EventField = Sample
	FieldNone EventField = 0
)

const (
	Centi EventUnit = (1 << iota)
	Milli
	Celcius
	Pascal
	Percent
	Metre
	Volt
	UnitMax  EventUnit = Volt
	UnitNone EventUnit = 0
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f EventField) String() string {
	if f == FieldNone {
		return f.flag()
	}
	str := ""
	for v := EventField(1); v <= FieldMax; v <<= 1 {
		if f&v == v {
			str += "|" + v.flag()
		}
	}
	return strings.TrimPrefix(str, "|")
}

func (u EventUnit) String() string {
	if u == UnitNone {
		return u.flag()
	}
	str := ""
	for v := EventUnit(1); v <= UnitMax; v <<= 1 {
		if u&v == v {
			str += v.flag()
		}
	}
	return str
}

func (f EventField) flag() string {
	switch f {
	case FieldNone:
		return "None"
	case Temperature:
		return "Temperature"
	case Pressure:
		return "Pressure"
	case Humidity:
		return "Humidity"
	case Altitude:
		return "Altitude"
	case Battery:
		return "Battery"
	case Charging:
		return "Charging"
	case Discharging:
		return "Discharging"
	case Button:
		return "Button"
	case Sample:
		return "Sample"
	default:
		return "??EventField"
	}
}

func (f EventUnit) flag() string {
	switch f {
	case UnitNone:
		return ""
	case Centi:
		return "Centi"
	case Milli:
		return "Milli"
	case Celcius:
		return "Celcius"
	case Pascal:
		return "Pascal"
	case Percent:
		return "Percent"
	case Metre:
		return "Metre"
	case Volt:
		return "Volt"
	default:
		return "??EventUnit"
	}
}
