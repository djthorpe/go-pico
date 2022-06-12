//go:build rp2040 && !production

package pico

import "fmt"

//////////////////////////////////////////////////////////////////////////////
// ASSERT

//go:inline
func assert(cond bool, err error) error {
	if !cond {
		fmt.Println(err)
		return err
	}
	return nil
}
