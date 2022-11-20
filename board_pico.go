//go:build pico

package pico

// map_spi maps from a GPIO pin to a SPI device
var map_spi = map[Pin]SPI{
	Pin(0): SPI{Num: 0, RX: Pin(0), TX: Pin(3), SCK: Pin(2), CS: Pin(1)},
	Pin(8): SPI{Num: 1, RX: Pin(8), TX: Pin(11), SCK: Pin(10), CS: Pin(9)},
}

// map_adc maps from a GPIO pin to an ADC channel
var map_adc = map[Pin]ADC{
	Pin(26): ADC{Num: 0},
	Pin(27): ADC{Num: 1},
	Pin(28): ADC{Num: 2},
	Pin(29): ADC{Num: 3},
}
