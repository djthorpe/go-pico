package main

import (
	"os"
	"time"

	// Modules
	gpio "github.com/djthorpe/go-pico/pkg/gpio"
	uart "github.com/djthorpe/go-pico/pkg/uart"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

var (
	UARTConfig = uart.Config{BaudRate: 115200, DataBits: 8, StopBits: 1}
	LEDPin     = Pin(25)
	GPIOConfig = gpio.Config{Out: []Pin{LEDPin}}
)

func main() {
	// Create console
	stdout, err := UARTConfig.New()
	if err != nil {
		panic(err)
	}

	// Create GPIO
	gpio, err := GPIOConfig.New()
	if err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}

	stdout.Println("loaded", gpio)

	// Blink lights
	for {
		gpio.High(LEDPin)
		time.Sleep(time.Millisecond * 800)
		gpio.Low(LEDPin)
		time.Sleep(time.Millisecond * 200)
	}
}
