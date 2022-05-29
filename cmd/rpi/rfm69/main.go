package main

import (
	"fmt"
	"os"

	// Modules
	rfm69 "github.com/djthorpe/go-pico/pkg/rfm69"
	// Namespace imports
	//. "github.com/djthorpe/go-pico"
)

// Device configuration
var (
	RFM69Config = rfm69.Config{Bus: 0, Slave: 1}
)

// Main
func main() {
	fmt.Println("=> rfm69")

	// Create BME280 device
	radio, err := RFM69Config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// Print device info
	fmt.Println(radio)
}
