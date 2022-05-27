package main

import (
	"os"

	// Modules
	i2c "github.com/djthorpe/go-pico/pkg/i2c"
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

	// Create I2C device
	device, err := i2c.Config{
		Bus:       0,
		Frequency: 100 * 1000, // 100KhZ
	}.New()
	if err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}

	stdout.Println("=> i2cdetect")
	stdout.Println("   0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F")

	data := make([]uint8, 1)
	for slave := uint8(0x00); slave <= uint8(0x7F); slave++ {
		if (slave % 16) == 0 {
			stdout.Printf("%02X ", slave)
		}

		// Perform a 1-byte dummy read from the probe address. If a slave
		// acknowledges this address, the function returns the number of bytes
		// transferred.

		// Skip over any reserved addresses
		if reserved_addr(slave) {
			stdout.Print("R")
		} else if err := device.Tx(slave, nil, data); err != nil {
			stdout.Print(err)
		} else {
			stdout.Print("@")
		}

		if slave%16 == 15 {
			stdout.Print("\n")
		} else {
			stdout.Print("  ")
		}
	}

	stdout.Println("<= i2cdetect")
}

// I2C reserves some addresses for special purposes. We exclude these from the scan.
// These are any addresses of the form 000 0xxx or 111 1xxx
func reserved_addr(addr uint8) bool {
	return addr == 0x78
	//return (addr&0x78) == 0 || (addr&0x78) == 0x78
}
