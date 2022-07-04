//go:build !rpi

package pico

import (
	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/sdk"
)

//////////////////////////////////////////////////////////////////////////////
// TYPES

type ADC struct {
	ch   uint32
	temp bool
}

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	adc = [NUM_ADC_CHANNELS]*ADC{}
)

//////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func _NewADC(ch uint32) *ADC {
	if ch > ADC_temperature_input() {
		return nil
	} else if adc := adc[ch]; adc != nil {
		return adc
	}

	// Initialise a new ADC
	adc[ch] = &ADC{
		ch: ch,
	}

	// Return the ADC
	return adc[ch]
}

//////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Get returns the raw ADC value, which is the first 12 bits
//
func (a *ADC) Get() uint16 {
	ADC_select_input(a.ch)
	for {
		if ADC_is_ready() {
			break
		}
	}
	return ADC_read()
}

// Return voltage given the value of the reference voltage
func (a *ADC) GetVoltage(vref float32) float32 {
	return float32(a.Get()) * vref / float32(1<<12)
}

// Return temperature ReadTemperature does a one-shot sample of the internal
// temperature sensor and returns a celsius reading.
//
// Only works on the  channel. Other channels will return 0
//
func (a *ADC) GetTemperature() float32 {
	if a.ch != ADC_temperature_input() {
		return 0
	}
	if a.temp == false {
		ADC_set_temp_sensor_enabled(true)
		a.temp = true
	}
	// Section 4.9.5
	return 27 - (a.GetVoltage(3.3)-0.706)/0.001721
}
