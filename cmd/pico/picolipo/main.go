package main

import (
	"time"

	// Modules
	picolipo "github.com/djthorpe/go-pico/pkg/picolipo"
	uart "github.com/djthorpe/go-pico/pkg/uart"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

var (
	stdout = uart.New(uart.Config{BaudRate: 115200, DataBits: 8, StopBits: 1})
)

// Main
func main() {
	ch := make(chan Event)
	picolipo := picolipo.New(picolipo.Config{Temperature: true}, ch)

	// Print out info
	stdout.Println(picolipo)

	// Sample in the background, once per 5 seconds
	go sample(picolipo, time.Millisecond*5000)

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

// Print out events
func receive(ch <-chan Event) {
	// Output events in the foreground
	for evt := range ch {
		stdout.Println(evt)
	}
}
