//go:build debug

package pico

import (
	"fmt"
)

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v Pin) String() string {
	return fmt.Sprint("GP", uint(v))
}
