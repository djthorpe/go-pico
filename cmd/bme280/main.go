package main

import (
	"os"
	"time"

	// Modules
	bme280 "github.com/djthorpe/go-pico/pkg/bme280"
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

	stdout.Println("=> bme280")

	// Create BME280 device
	device, err := bme280.Config{}.New()
	if err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}

	// Reset device
	if err := device.Reset(); err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}

	// Set temperature oversampling
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

	// Read temperature every second
	for {
		// Force sampling
		if err := device.Sample(); err != nil {
			stdout.Println(err)
			os.Exit(-1)
		}
		time.Sleep(time.Millisecond * 100)
	}
}
