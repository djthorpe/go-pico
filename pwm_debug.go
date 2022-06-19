//go:build !production

package pico

import "fmt"

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v *PWM) String() string {
	str := "<pwm"
	str += fmt.Sprint(" slice_num=", v.slice_num)
	return str + ">"
}
