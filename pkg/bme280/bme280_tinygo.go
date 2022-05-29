//go:build tinygo

package bme280

import (
	i2c "github.com/djthorpe/go-pico/pkg/i2c"
	spi "github.com/djthorpe/go-pico/pkg/spi"
)

func NewI2C(cfg I2CConfig) (*device, error) {
	return i2c.Config{Bus: cfg.Bus, Speed: cfg.Speed | DEFAULT_I2C_SPEED}.New()
}

func NewSPI(cfg SPIConfig) (*device, error) {
	return spi.Config{Bus: cfg.Bus, Speed: cfg.Speed | DEFAULT_SPI_SPEED}.New()
}
