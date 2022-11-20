//go:build !production

package sdk

//////////////////////////////////////////////////////////////////////////////
// ASSERT

//go:inline
func assert(cond bool) {
	if !cond {
		panic("assertation failed")
	}
}
