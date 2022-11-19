package main

import (
	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

// Define the pins used
var (
	LED    = Pin(25)
	BUTTON = Pin(23) // BOOTSEL pin on the Pico Lipo
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
func on_button(Pin) {
	LED.Set(!LED.Get())
}
