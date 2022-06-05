//go:build tinygo

package adc

import (
	"fmt"
	"machine"
	"sync"

	// Package imports
	event "github.com/djthorpe/go-pico/pkg/event"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Channel   uint   // ADC Channel (0..3) with 3 being temperature sensor
	Reference uint32 // analog reference voltage (AREF) in millivolts
	Samples   uint32 // number of samples for a single conversion (e.g., 4, 8, 16, 32)
}

type device struct {
	adc machine.ADC
	pin Pin
	ch  uint
	e   chan<- Event
	wg  sync.WaitGroup
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	lock sync.Mutex
	once sync.Once
)

const (
	DEFAULT_RESOLUTION = 12 // 12 bits
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create new ADC device, with channel. Panics on error
func New(cfg Config, ch chan<- Event) *device {
	if d, err := cfg.New(ch); err != nil {
		panic(err)
	} else {
		return d
	}
}

func (cfg Config) New(e chan<- Event) (*device, error) {
	this := new(device)

	// Initialize ADC
	once.Do(func() {
		machine.InitADC()
	})

	// Channel
	if e == nil {
		return nil, ErrBadParameter.With("adc:", cfg.Channel)
	} else {
		this.e = e
	}

	// Pico has four ADC channels 0-4
	ch := machine.ADCChannel(cfg.Channel)
	if ch == machine.ADC_TEMP_SENSOR {
		// There is no pin associated with temperature sensor
		if err := ch.Configure(machine.ADCConfig{
			Resolution: DEFAULT_RESOLUTION,
			Reference:  cfg.Reference,
			Samples:    cfg.Samples,
		}); err != nil {
			return nil, ErrBadParameter.With("adc:", cfg.Channel, err)
		}
	} else {
		if pin, err := ch.Pin(); err != nil {
			return nil, ErrBadParameter.With("adc:", cfg.Channel, err)
		} else {
			this.adc = machine.ADC{Pin: pin}
			this.pin = Pin(pin)
		}

		// Configure ADC
		if err := this.adc.Configure(machine.ADCConfig{
			Resolution: DEFAULT_RESOLUTION,
			Reference:  cfg.Reference,
			Samples:    cfg.Samples,
		}); err != nil {
			return nil, ErrBadParameter.With("adc:", cfg.Channel, err)
		}
	}

	// Set channel
	this.ch = uint(ch)

	// Return success
	return this, nil
}

func (d *device) Close() error {
	var result error

	// Wait for any background tasks to complete
	d.wg.Wait()

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d *device) String() string {
	str := "<adc"
	str += fmt.Sprint(" ch=", d.ch)
	if d.pin != 0 {
		str += fmt.Sprint(" pin=", d.pin)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the logical pin number for the ADC
func (d *device) Pin() Pin {
	return d.pin
}

// Initiate reading a sample
func (d *device) Sample() error {
	e := event.New(d)
	lock.Lock()
	if d.adc.Pin != 0 {
		e.Set(Sample, UnitNone, d.adc.Get())
	} else {
		e.Set(Temperature, Milli|Celcius, machine.ADCChannel(d.ch).ReadTemperature())
	}
	lock.Unlock()
	return e.Emit(d.e)
}
