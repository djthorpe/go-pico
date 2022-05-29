// go:build linux

package spi_test

import (
	"testing"

	// Namespace import
	. "github.com/djthorpe/go-pico/pkg/spi"
)

func Test_SPI_001(t *testing.T) {
	for bus := uint(0); bus <= 2; bus++ {
		for slave := uint(0); slave <= 2; slave++ {
			spi, err := Config{Bus: bus, Slave: slave}.New()
			if err != nil {
				t.Log("Not found:", bus, ",", slave, ":", err)
			} else {
				t.Log(spi)
				if err := spi.Close(); err != nil {
					t.Error(err)
				}
			}
		}
	}
}
