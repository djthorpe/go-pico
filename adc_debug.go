//go:build debug

package pico

import "fmt"

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v *ADC) String() string {
	str := "<adc"
	if v.Pin != 0 {
		str += fmt.Sprint(" pin=", v.Pin)
	}
	str += fmt.Sprint(" ch=", v.Num)
	return str + ">"
}
