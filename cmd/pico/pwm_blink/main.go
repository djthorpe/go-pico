package main

import (
	. "github.com/djthorpe/go-pico/pkg/pico"
)

var (
	PIN_LED = Pin(25)
)

var (
	fade     uint16
	going_up bool
)

func main() {
	PIN_LED.SetMode(ModePWM)
	PIN_LED.PWM().SetEnabled(true)
	PIN_LED.PWM().SetInterrupt(on_pwm_wrap)

	select {}
}

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
	pwm.Set(PIN_LED, fade*fade)
}
