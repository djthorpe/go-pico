package main

import (
	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

// Define the pins used
var (
	LED = Pin(25)
)

// Global variables, for the PWM state
var (
	fade     uint16
	going_up bool
)

// Main function
func main() {
	LED.SetMode(ModePWM)
	LED.PWM().SetEnabled(true)
	LED.PWM().SetInterrupt(on_pwm_wrap)
	// Wait forever
	select {}
}

// Called when the PWM counter wraps
func on_pwm_wrap(pwm *PWM) {
	if going_up {
		fade = fade + 1
		if fade > 255 {
			fade = 255
			going_up = false
		}
	} else {
		if fade == 0 {
			going_up = true
		} else {
			fade = fade - 1
		}
	}
	// Square the fade value to make the LED's brightness appear more linear
	// Note this range matches with the wrap value
	pwm.Set(LED, fade*fade)
}
