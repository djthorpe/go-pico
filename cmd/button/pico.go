package main

import (
	. "github.com/djthorpe/go-pico"
)

// Define the pins used
var (
	LED    = Pin(25) // On-board LED
	BUTTON = Pin(23) // BOOTSEL button on the Pico Lipo
)

// Main function
func main() {
	LED.SetMode(ModeOutput)
	LED.Set(true)
	BUTTON.SetMode(ModeInput)
	BUTTON.SetInterrupt(on_button)

	// Wait forever
	select {}
}

// Called when the button is pressed or released
func on_button(p Pin, s State) {
	if p == BUTTON && s&StateRise != 0 {
		LED.Set(!LED.Get())
	}
}
