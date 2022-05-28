package main

import (
	"os"
	"time"

	// Modules
	bme280 "github.com/djthorpe/go-pico/pkg/bme280"
	uart "github.com/djthorpe/go-pico/pkg/uart"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
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

	stdout.Println("=> bme280")

	// Create a channel to receive events in background
	ch := make(chan bme280.Event)

	// Create BME280 device
	device, err := bme280.SPIConfig{}.New(ch)
	if err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}

	// Reset device
	if err := device.Reset(); err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}

	// Set oversampling
	if err := device.SetTempOversample(bme280.BME280_OVERSAMPLE_16); err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}
	if err := device.SetHumidityOversample(bme280.BME280_OVERSAMPLE_16); err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}
	if err := device.SetPressureOversample(bme280.BME280_OVERSAMPLE_16); err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}

	// Print device info
	stdout.Println(device)

	// Sample in the background
	go sample(device, stdout, time.Millisecond*1000)

	// Receive in the foreground (blocking)
	receive(stdout, ch)
}

func sample(device BME280, stdout UART, frequency time.Duration) {
	// Read temperature every second
	for {
		// Force sampling
		if err := device.Sample(); err != nil {
			stdout.Println(err)
		}
		time.Sleep(frequency)
	}
}

func receive(stdout UART, ch chan bme280.Event) {
	// Output events in the foreground
	for {
		select {
		case evt := <-ch:
			stdout.Println(evt)
		}
	}
}
