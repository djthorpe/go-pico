//go:build !production

package pico

import "fmt"

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v *ADC) String() string {
	str := "<adc"
	str += fmt.Sprint(" ch=", v.ch)
	return str + ">"
}
