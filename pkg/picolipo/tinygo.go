//go:build tinygo

package adc

import (
	"fmt"
	"sync"

	// Package imports
	adc "github.com/djthorpe/go-pico/pkg/adc"
	event "github.com/djthorpe/go-pico/pkg/event"
	gpio "github.com/djthorpe/go-pico/pkg/gpio"
	multierror "github.com/hashicorp/go-multierror"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Temperature bool // Set to true to read temperature
}

type device struct {
	temp, battery ADC
	gpio          GPIO
	in            chan Event
	out           chan<- Event
	volts         uint16
	celcius       uint32
	wg            sync.WaitGroup
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DEFAULT_SAMPLES = 10
	PIN_BOOTSEL     = Pin(23)
	PIN_CHARGING    = Pin(24)
	BATTERY_FACTOR  = 3 * 3.3 / 65535
	BATTERY_FULL    = 4.2 //  reference voltages for a full/empty battery, in volts
	BATTERY_EMPTY   = 2.8 // the values could vary by battery size/manufacturer so you might need to adjust them
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

func (cfg Config) New(out chan<- Event) (*device, error) {
	this := new(device)

	// Channel
	if out == nil {
		return nil, ErrBadParameter.With("picolipo")
	} else {
		this.out = out
		this.in = make(chan Event)
	}

	// Watch the charging pin for plugging/unplugging
	g := gpio.Config{
		In:    []Pin{PIN_CHARGING},
		Watch: []Pin{PIN_CHARGING},
	}
	if gpio, err := g.New(this.in); err != nil {
		return nil, err
	} else {
		this.gpio = gpio
	}

	// Temperature
	if cfg.Temperature {
		c := adc.Config{
			Channel: 4,
			Samples: DEFAULT_SAMPLES,
		}
		if adc, err := c.New(this.in); err != nil {
			return nil, err
		} else {
			this.temp = adc
		}
	}

	// Battery
	d := adc.Config{
		Channel: 3,
		Samples: DEFAULT_SAMPLES,
	}
	if adc, err := d.New(this.in); err != nil {
		return nil, err
	} else {
		this.battery = adc
	}

	// Receive events from the "in" channel and
	// post them on the "out" channel
	this.wg.Add(1)
	go this.receive()

	// Return success
	return this, nil
}

func (d *device) Close() error {
	var result error

	// Close temp and battery ADC channels
	if d.temp != nil {
		if err := d.temp.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if d.battery != nil {
		if err := d.battery.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if d.gpio != nil {
		if err := d.gpio.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Close in channel, and wait for it to close
	close(d.in)
	d.wg.Wait()

	// Release resources
	d.temp = nil
	d.battery = nil
	d.gpio = nil
	d.in = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d *device) String() string {
	str := "<picolipo"
	if d.Charging() {
		str += " charging"
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Initiate reading samples for battery and temperature
func (d *device) Sample() error {
	var result error
	if d.temp != nil {
		if err := d.temp.Sample(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if d.battery != nil {
		if err := d.battery.Sample(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

// Return charging state
func (d *device) Charging() bool {
	return d.gpio.Get(PIN_CHARGING)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (d *device) receive() {
	for evt := range d.in {
		switch evt.Source() {
		case d.temp:
			if evt.Is(Temperature) {
				// Emit changed temperature value
				v, u := evt.Value(Temperature)
				if v_ := v.(uint32); v_ != d.celcius {
					d.celcius = v_
					event.New(d).Set(Temperature, u, v_).Emit(d.out)
				}
			}
		case d.battery:
			if evt.Is(Sample) {
				// Emit changed battery value
				v, _ := evt.Value(Sample)
				if v_ := v.(uint16); v_ != d.volts {
					d.volts = v_
					vf := float32(v_) * BATTERY_FACTOR
					vp := uint16(100.0 * ((vf - BATTERY_EMPTY) / (BATTERY_FULL - BATTERY_EMPTY)))
					if vp > 100 {
						vp = 100
					}
					if d.Charging() {
						event.New(d).Set(Battery, Volt, vf).Set(Charging, Percent, vp).Emit(d.out)
					} else {
						event.New(d).Set(Battery, Volt, vf).Set(Discharging, Percent, vp).Emit(d.out)
					}
				}
			}
		case PIN_CHARGING:
			if evt.Is(Sample) {
				v, _ := evt.Value(Sample)
				if v_, ok := v.(bool); ok {
					if v_ {
						event.New(d).Set(Charging, UnitNone, v).Emit(d.out)
					} else {
						event.New(d).Set(Discharging, UnitNone, v).Emit(d.out)
					}
				}
			}
		}
	}
	fmt.Println("finished receiving events")
	d.wg.Done()
}
