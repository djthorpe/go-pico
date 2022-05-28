package bme280

import (
	"fmt"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type EventType uint

type Event struct {
	Type        EventType
	Temperature int32
	Humidity    int32
	Pressure    int32
	Altitude    int32
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	Temperature EventType = (1 << iota)
	Pressure
	Humidity
	Altitude
	Max  EventType = Altitude
	None EventType = 0
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (e Event) Is(t EventType) bool {
	return e.Type&t == t
}

func (e Event) TempCelcius() float32 {
	if e.Is(Temperature) {
		return float32(e.Temperature) / 1000.0
	} else {
		return 0.0
	}
}

func (e Event) PressurePascals() float32 {
	if e.Is(Pressure) {
		return float32(e.Pressure) / 1000.0
	} else {
		return 0.0
	}
}

func (e Event) HumidityPercent() float32 {
	if e.Is(Humidity) {
		return float32(e.Humidity) / 100.0
	} else {
		return 0.0
	}
}

func (e Event) AltitudeMetres() float32 {
	if e.Is(Altitude) {
		return float32(e.Altitude) / 100.0
	} else {
		return 0.0
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e Event) String() string {
	str := "<event"
	str += fmt.Sprint(" type=", e.Type)
	if e.Is(Temperature) {
		str += fmt.Sprintf(" temp=%.1fC", e.TempCelcius())
	}
	if e.Is(Pressure) {
		str += fmt.Sprintf(" pressure=%.1fhPa", e.PressurePascals()/100.0)
	}
	if e.Is(Humidity) {
		str += fmt.Sprintf(" humidity=%.0f%%", e.HumidityPercent())
	}
	if e.Is(Altitude) {
		str += fmt.Sprintf(" altitude=%.1fm", e.AltitudeMetres())
	}
	return str + ">"
}

func (t EventType) String() string {
	if t == None {
		return t.flag()
	}
	str := ""
	for v := EventType(1); v <= Max; v <<= 1 {
		if t&v == v {
			str += "|" + v.flag()
		}
	}
	return strings.TrimPrefix(str, "|")
}

func (t EventType) flag() string {
	switch t {
	case None:
		return "None"
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
