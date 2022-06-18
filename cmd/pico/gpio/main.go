package main

import (
	"fmt"
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

var (
	BOOTSEL = Pin(23)
)

func main() {
	fmt.Println("=>gpio")
	BOOTSEL.SetMode(ModeInputPulldown)
	BOOTSEL.SetInterrupt(on_gpio_change)

	for {
		fmt.Println(BOOTSEL, BOOTSEL.Mode(), BOOTSEL.Get())
		time.Sleep(time.Second)
	}
}

func on_gpio_change(pin Pin) {
	fmt.Println(pin, "on_gpio_change")
}
