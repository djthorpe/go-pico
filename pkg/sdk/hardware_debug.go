//go:build debug

package sdk

//////////////////////////////////////////////////////////////////////////////
// ASSERT

//go:inline
func assert(cond bool) {
	if !cond {
		panic("assertation failed")
	}
}
