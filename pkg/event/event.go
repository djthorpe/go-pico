package pico

import (
	"fmt"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	EventSource
	flags         EventField
	boolValue     bool  // Pure boolean value
	tempValue     int32 // Temperature in milli-degrees Celsius
	humidityValue int32 // Relative Humidity in centi-percent
	pressureValue int32 // Pressure in milli-pascals
	altitudeValue int32 // Altitude in centi-metres
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func New(source EventSource) Event {
	return &event{EventSource: source}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (e *event) Source() EventSource {
	return e.EventSource
}

func (e *event) Is(f EventField) bool {
	return e.flags&f == f
}

func (e *event) SetBool(v bool) {
	e.flags |= Bool
	e.boolValue = v
}

func (e *event) Bool() bool {
	return e.boolValue
}

func (e *event) SetTemperature(v int32) {
	e.flags |= Temperature
	e.tempValue = v
}

func (e *event) Temperature() int32 {
	if e.Is(Temperature) {
		return e.tempValue
	} else {
		return 0
	}
}

func (e *event) Pressure() int32 {
	if e.Is(Pressure) {
		return e.pressureValue
	} else {
		return 0
	}
}

func (e *event) SetPressure(v int32) {
	e.flags |= Pressure
	e.pressureValue = v
}

func (e *event) Humidity() int32 {
	if e.Is(Humidity) {
		return e.humidityValue
	} else {
		return 0
	}
}

func (e *event) SetHumidity(v int32) {
	e.flags |= Humidity
	e.humidityValue = v
}

func (e *event) Altitude() int32 {
	if e.Is(Altitude) {
		return e.altitudeValue
	} else {
		return 0
	}
}

func (e *event) SetAltitude(v int32) {
	e.flags |= Altitude
	e.altitudeValue = v
}

func (e *event) SendOn(C chan<- Event) error {
	select {
	case C <- e:
		return nil
	default:
		return ErrTimeout
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e *event) String() string {
	str := "<event"
	if e.Is(Bool) {
		str += fmt.Sprintf(" bool=%v", e.Bool())
	}
	if e.Is(Temperature) {
		str += fmt.Sprintf(" temp=%.1fC", float32(e.Temperature())/1000.0)
	}
	if e.Is(Pressure) {
		str += fmt.Sprintf(" pressure=%.1fhPa", float32(e.Pressure())/100.0)
	}
	if e.Is(Humidity) {
		str += fmt.Sprintf(" humidity=%.0f%%", float32(e.Humidity())/100.0)
	}
	if e.Is(Altitude) {
		str += fmt.Sprintf(" altitude=%.1fm", float32(e.Altitude())/100.0)
	}
	return str + ">"
}
