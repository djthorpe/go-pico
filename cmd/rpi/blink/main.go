package main

import (
	"fmt"
	"os"
	"time"

	// Modules
	gpio "github.com/djthorpe/go-pico/pkg/gpio"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

var (
	LEDPin     = Pin(22) // GPIO22
	GPIOConfig = gpio.Config{
		Out: []Pin{LEDPin},
	}
)

func main() {
	// Create GPIO
	gpio, err := GPIOConfig.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("loaded", gpio)

	// Blink lights
	for {
		gpio.High(LEDPin)
		time.Sleep(time.Millisecond * 800)
		gpio.Low(LEDPin)
		time.Sleep(time.Millisecond * 200)
	}
}
