package main

import (
	"fmt"
	"time"

	// Modules

	"github.com/djthorpe/go-pico/pkg/picodisplay"
	picolipo "github.com/djthorpe/go-pico/pkg/picolipo"
	spi "github.com/djthorpe/go-pico/pkg/spi"
	st7789 "github.com/djthorpe/go-pico/pkg/st7789"
	uart "github.com/djthorpe/go-pico/pkg/uart"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

var (
	// Output
	stdout = uart.New(uart.Config{BaudRate: 115200, DataBits: 8, StopBits: 1})

	// Channel for events
	ch = make(chan Event)

	// Display
	display = st7789.New(st7789.Config{
		SPI: spi.New(spi.Config{}),
	}, ch)

	// Pico Lipo
	board = picolipo.New(picolipo.Config{
		Temperature: true,
	}, ch)

	// Pico Display
	picodisp = picodisplay.New(picodisplay.Config{}, ch)
)

// Main
func main() {
	// Print out info
	stdout.Println(board, display, picodisp)

	picodisp.SetLED(true, true, true)

	// Sample in the background, once per 5 seconds
	go sample(board, time.Millisecond*5000)

	// Switch backlight on and off
	go backlight(display, time.Millisecond*5000)

	// Receive events in the foreground
	receive(ch)
}

// Sample
func sample(device ADC, frequency time.Duration) {
	i := 0
	// Read sample every second
	for {
		if i == 0 {
			// Force sampling
			if err := device.Sample(); err != nil {
				stdout.Println(err)
			}
		}
		time.Sleep(time.Millisecond * 100)
		i = (i + 1) % 50
	}
}

// Backlight
func backlight(display Display, frequency time.Duration) {
	i := 0
	// Read sample every second
	for {
		if i == 0 {
			display.SetBacklight(!display.Backlight())
			fmt.Println(display)
		}
		time.Sleep(time.Millisecond * 100)
		i = (i + 1) % 50
	}
}

// Print out events
func receive(ch <-chan Event) {
	// Output events in the foreground
	for evt := range ch {
		stdout.Println(evt)
	}
}
