package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {
	fmt.Println("Loaded blink\r\n")

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		led.Low()
		time.Sleep(time.Millisecond * 200)
		led.High()
		time.Sleep(time.Millisecond * 1000)
	}
}
