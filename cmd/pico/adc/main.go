package main

import (
	"fmt"
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

var (
	BOOTSEL        = Pin(23)
	CHARGING       = Pin(24)
	BATTERY_ADC    = Pin(29).ADC()
	BATTERY_FACTOR = 3 * 3.3 / (1 << 12)
	BATTERY_FULL   = 4.2 // Reference voltages for a full/empty battery, in volts
	BATTERY_EMPTY  = 2.8 // the values could vary by battery size/manufacturer
)

func main() {
	fmt.Println("=>adc")

	fmt.Println(BATTERY_ADC)

	for {
		battery_v := float32(BATTERY_ADC.Get()) * float32(BATTERY_FACTOR)
		battery_pct := 100 * (battery_v - float32(BATTERY_EMPTY)) / float32(BATTERY_FULL-BATTERY_EMPTY)
		if battery_pct > 100 {
			battery_pct = 100
		}
		if battery_pct < 0 {
			battery_pct = 0
		}
		fmt.Printf("%v %.1f%% %v\n", BATTERY_ADC, battery_pct, CHARGING.Get())
		time.Sleep(time.Second)
	}
}
