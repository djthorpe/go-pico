//go:build debug

package pico

import "fmt"

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v *SPI) String() string {
	str := "<spi"
	str += fmt.Sprint(" num=", v.Num)
	if v.Baud > 0 {
		str += fmt.Sprint(" baud=", v.Baud)
	}
	str += fmt.Sprint(" rx=", v.RX)
	str += fmt.Sprint(" tx=", v.TX)
	str += fmt.Sprint(" sck=", v.SCK)
	str += fmt.Sprint(" cs=", v.CS)
	return str + ">"
}
