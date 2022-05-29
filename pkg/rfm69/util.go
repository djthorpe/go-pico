package rfm69

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func to_uint8_bool(value uint8) bool {
	return (value != 0x00)
}

func to_bool_uint8(value bool) uint8 {
	if value {
		return 0x01
	} else {
		return 0x00
	}
}

func to_bitratehz(v uint16) uint {
	return uint(float32(RFM_FXOSC_MHZ*1e6) / float32(v))
}
