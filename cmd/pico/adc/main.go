package main

import (
	"fmt"
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

var (
	BOOTSEL     = Pin(23)
	BATTERY_ADC = Pin(29).ADC()
)

func main() {
	fmt.Println("=>adc")

	fmt.Println(BATTERY_ADC)

	for {
		fmt.Println(BATTERY_ADC, BATTERY_ADC.Get())
		time.Sleep(time.Second)
	}
}
