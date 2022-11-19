package main

import (
	"os"
	"time"

	// Modules
	bme280 "github.com/djthorpe/go-pico/pkg/bme280"
	spi "github.com/djthorpe/go-pico/pkg/spi"
	uart "github.com/djthorpe/go-pico/pkg/uart"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

// Main
func main() {
	stdout = uart.New({BaudRate: 115200, DataBits: 8, StopBits: 1})
	spi = spi.New({Bus: 0, Slave: 1})
	bme280 = bme280.New({ SPI: spi})

	// Create console
	stdout, err := UARTConfig.New()
	if err != nil {
		panic(err)
	}

	stdout.Println("=> bme280")

	// Create BME280 device
	device, err := BME280Config.New()
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

	// Receive in the background
	go receive(stdout, device.C())

	// Sample in the foreground, once per second
	sample(device, stdout, time.Millisecond*1000)
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

func receive(stdout UART, ch <-chan bme280.Event) {
	// Output events in the foreground
	for evt := range ch {
		stdout.Println(evt)
	}
}
