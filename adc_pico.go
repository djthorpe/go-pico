//go:build pico

package pico

import ( // Module imports
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

const (
	// Temperature channel is the last one
	ADC_TEMP_SENSOR = NUM_ADC_CHANNELS - 1
)

var (
	adc = [NUM_ADC_CHANNELS]*ADC{}
)

//////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func _NewADC(ch uint32) *ADC {
	if ch > ADC_TEMP_SENSOR {
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
// METHODS

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

// Return temperature ReadTemperature does a one-shot sample of the internal temperature sensor and returns a milli-celsius reading.
// Only works on the  channel. aka AINSEL=4. Other channels will return 0
func (a *ADC) GetTemperature() float32 {
	if a.ch != ADC_TEMP_SENSOR {
		return 0
	}
	if a.temp == false {
		ADC_set_temp_sensor_enabled(true)
		a.temp = true
	}
	// Section 4.9.5
	return 27 - (a.GetVoltage(3.3)-0.706)/0.001721
}
