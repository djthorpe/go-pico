//go:build !debug

package pico

//////////////////////////////////////////////////////////////////////////////
// ASSERT

//go:inline
func assert(cond bool, err error) error {
	if !cond {
		return err
	}
	return nil
}
