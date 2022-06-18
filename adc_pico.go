//go:build pico

package pico

import ( // Module imports
	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/sdk"
)

//////////////////////////////////////////////////////////////////////////////
// TYPES

type ADC struct {
	ch uint32
}

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ADC_CH_TEMPERATURE = ADC_NUM_CHANNELS
)

var (
	adc = [NUM_ADC_CHANNELS]*ADC{}
)

//////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func _NewADC(ch uint32) *ADC {
	if ch > ADC_CH_TEMPERATURE {
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
	return ADC_read()
}

// GetVoltage returns
/*

// getOnce returns a one-shot ADC sample reading from an ADC channel.
func (c ADCChannel) getOnce() uint16 {
	// Make it safe to sample multiple ADC channels in separate go routines.
	adcLock.Lock()
	rp.ADC.CS.ReplaceBits(uint32(c), 0b111, rp.ADC_CS_AINSEL_Pos)
	rp.ADC.CS.SetBits(rp.ADC_CS_START_ONCE)

	waitForReady()
	adcLock.Unlock()

	// rp2040 is a 12-bit ADC, scale raw reading to 16-bits.
	return uint16(rp.ADC.RESULT.Get()) << 4
}

// getVoltage does a one-shot sample and returns a millivolts reading.
// Integer portion is stored in the high 16 bits and fractional in the low 16 bits.
func (c ADCChannel) getVoltage() uint32 {
	return (adcAref << 16) / (1 << 12) * uint32(c.getOnce()>>4)
}

// ReadTemperature does a one-shot sample of the internal temperature sensor and returns a milli-celsius reading.
// Only works on the ADC_TEMP_SENSOR channel. aka AINSEL=4. Other channels will return 0
func (c ADCChannel) ReadTemperature() (millicelsius uint32) {
	if c != ADC_TEMP_SENSOR {
		return
	}
	// T = 27 - (ADC_voltage - 0.706)/0.001721
	return (27000<<16 - (c.getVoltage()-706<<16)*581) >> 16
}
*/
