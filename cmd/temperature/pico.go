package main

import (
	"fmt"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

// Define the pins used
var (
	BUTTON = Pin(23) // BOOTSEL button on the Pico Lipo
	TEMP   = GPIO.Temperature()
)

// Main function
func main() {
	CH := make(chan uint16, 10)

	BUTTON.SetMode(ModeInput)
	BUTTON.SetInterrupt(func(p Pin, s State) {
		if s == StateFall {
			CH <- TEMP.Get()
		}
	})

	// Wait forever
	for evt := range CH {
		fmt.Println(evt)
	}
}

// Called when the button is pressed or released
