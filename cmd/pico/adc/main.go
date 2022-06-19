package main

import (
	"fmt"
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

var (
	BOOTSEL       = Pin(23)
	CHARGING      = Pin(24)
	BATTERY_ADC   = Pin(29).ADC()
	BATTERY_FULL  = 4.2 // Reference voltages for a full/empty battery, in volts
	BATTERY_EMPTY = 2.8 // the values could vary by battery size/manufacturer
)

func main() {
	fmt.Println("=>adc")

	fmt.Println(BATTERY_ADC)

	for {
		battery_s := BATTERY_ADC.Get()
		battery_v := float32(battery_s) * float32(3.3) / (float32(1<<12) - 1.0)
		fmt.Println("battery_s=", battery_s, "battery_v=", battery_v)
		battery_pct := (battery_v - float32(BATTERY_EMPTY)) / (float32(BATTERY_FULL) - float32(BATTERY_EMPTY))
		fmt.Printf("%v %.1f%% %v\n", BATTERY_ADC, battery_pct, CHARGING.Get())
		time.Sleep(time.Second)
	}
}
