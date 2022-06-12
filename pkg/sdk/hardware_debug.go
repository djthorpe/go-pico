//go:build rp2040 && !production

package sdk

//////////////////////////////////////////////////////////////////////////////
// ASSERT

//go:inline
func assert(cond bool) {
	if !cond {
		panic("assertation failed")
	}
}
