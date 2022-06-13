//go:build !production

package pico

import (
	"fmt"
)

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v Pin) String() string {
	return fmt.Sprint("GP", uint(v))
}
