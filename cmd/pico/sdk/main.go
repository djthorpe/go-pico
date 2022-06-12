package main

import (
	// Namespace imports
	"fmt"
	"time"

	. "github.com/djthorpe/go-pico/pkg/pico"
)

var (
	PIN_LED     = Pin(25)
	PIN_BOOTSEL = Pin(23)
)

func main() {
	PIN_LED.SetMode(ModePWM)
	PIN_BOOTSEL.SetMode(ModeInputPullup)

	fmt.Println(PIN_LED, PIN_LED.PWM())
	fmt.Println(PIN_BOOTSEL)

	PIN_LED.PWM().Set(0xFFFF)
	PIN_LED.PWM().SetEnabled(true)

	for {
		fmt.Println("PWM=", PIN_LED.PWM().Get())
		time.Sleep(time.Second)
		PIN_LED.PWM().Inc()
	}
}
