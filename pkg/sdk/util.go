//go:build rp2040

package sdk

//go:inline
func bool_to_bit(v bool) uint32 {
	if v {
		return 1
	}
	return 0
}
