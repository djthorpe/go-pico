package main

import (
	"machine"
	"time"

	"github.com/djthorpe/go-pico/pkg/uart"
)

func main() {
	stdout, err := uart.Config{
		BaudRate: 9600,
		DataBits: 8,
		StopBits: 1,
	}.New()
	if err != nil {
		panic(err)
	}

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for i := 0; i < 100; i++ {
		led.Low()
		stdout.Printf("Hello, World %v\n", i)
		time.Sleep(time.Millisecond * 900)
		led.High()
		time.Sleep(time.Millisecond * 100)
	}
}
