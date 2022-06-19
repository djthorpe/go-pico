package main

import (
	"fmt"
	"time"

	// Module imports
	math "github.com/djthorpe/go-pico/pkg/math32"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

var (
	BOOTSEL       = Pin(23)
	CHARGING      = Pin(24)
	BATTERY_ADC   = Pin(29).ADC()
	BATTERY_VREF  = float32(3 * 3.3)
	BATTERY_FULL  = float32(4.2) // Reference voltages for a full/empty battery, in volts
	BATTERY_EMPTY = float32(2.8) // the values could vary by battery size/manufacturer
)

func main() {
	fmt.Println("=>adc")

	for {
		battery_v := BATTERY_ADC.GetVoltage(BATTERY_VREF)
		battery_pct := (battery_v - BATTERY_EMPTY) * 100.0 / (BATTERY_FULL - BATTERY_EMPTY)
		battery_pct = math.Max(0, math.Min(battery_pct, 100.0))
		fmt.Printf("%v v=%.2fV %.0f%% charging=%v\n", BATTERY_ADC, battery_v, battery_pct, CHARGING.Get())
		time.Sleep(time.Second)
	}
}
