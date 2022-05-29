package main

import (
	"os"
	"time"

	// Modules
	gpio "github.com/djthorpe/go-pico/pkg/gpio"
	uart "github.com/djthorpe/go-pico/pkg/uart"
)

const (
	RED_LED   = gpio.Pin(15)
	GREEN_LED = gpio.Pin(14)
)

func main() {
	// Create console
	stdout, err := uart.Config{
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
	}.New()
	if err != nil {
		panic(err)
	}

	stdout.Println("=> gpio")

	// Create GPIO
	device, err := gpio.Config{
		Out: []gpio.Pin{RED_LED, GREEN_LED},
	}.New()
	if err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}

	stdout.Println(device)

	// Blink lights
	for {
		device.High(RED_LED)
		device.Low(GREEN_LED)
		time.Sleep(time.Millisecond * 500)
		device.High(GREEN_LED)
		device.Low(RED_LED)
		time.Sleep(time.Millisecond * 500)
	}
}
