package main

import (
	"os"

	// Modules
	rfm69 "github.com/djthorpe/go-pico/pkg/rfm69"
	uart "github.com/djthorpe/go-pico/pkg/uart"
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

	stdout.Println("=> rfm69")

	// Create GPIO
	device, err := rfm69.Config{}.New()
	if err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}

	stdout.Println(device)
}
