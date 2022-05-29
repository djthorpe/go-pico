//go:build rpi

package bme280

import (
	i2c "github.com/djthorpe/go-pico/pkg/i2c"
	spi "github.com/djthorpe/go-pico/pkg/spi"
)

func NewI2C(cfg I2CConfig) (*device, error) {
	return i2c.Config{Bus: cfg.Bus}.New()
}

func NewSPI(cfg SPIConfig) (*device, error) {
	return spi.Config{Bus: cfg.Bus, Slave: cfg.Slave, Speed: cfg.Speed | DEFAULT_SPI_SPEED}.New()
}
