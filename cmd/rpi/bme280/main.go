package main

import (
	"fmt"
	"os"
	"time"

	// Modules
	bme280 "github.com/djthorpe/go-pico/pkg/bme280"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

// Device configuration
var (
	BME280Config = bme280.SPIConfig{Bus: 0}
)

// Main
func main() {
	fmt.Println("=> bme280")

	// Create BME280 device
	device, err := BME280Config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// Reset device
	if err := device.Reset(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// Set oversampling
	if err := device.SetTempOversample(bme280.BME280_OVERSAMPLE_16); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := device.SetHumidityOversample(bme280.BME280_OVERSAMPLE_16); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := device.SetPressureOversample(bme280.BME280_OVERSAMPLE_16); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// Print device info
	fmt.Println(device)

	// Receive in the background
	go receive(device.C())

	// Sample in the foreground, once per second
	sample(device, time.Millisecond*1000)
}

func sample(device BME280, frequency time.Duration) {
	// Read temperature every second
	for {
		// Force sampling
		if err := device.Sample(); err != nil {
			fmt.Println(err)
		}
		time.Sleep(frequency)
	}
}

func receive(ch <-chan bme280.Event) {
	// Output events in the foreground
	for evt := range ch {
		fmt.Println(evt)
	}
}
