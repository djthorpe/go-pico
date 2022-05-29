//go:build production

package rfm69

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - UINT8

func (d *device) readreg_uint8(r register) (uint8, error) {
	var buf [2]uint8
	if err := d.SPI.Transfer([]byte{uint8(r & RFM_REG_MAX), 0}, buf[:]); err != nil {
		return 0, err
	}
	// Return success
	return buf[1], nil
}

func (d *device) writereg_uint8(r register, data uint8) error {
	return d.SPI.Transfer([]byte{byte((r & RFM_REG_MAX) | RFM_REG_WRITE), data}, nil)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - UINT16

func (d *device) readreg_uint16(r register) (uint16, error) {
	var buf [3]uint8
	if err := d.SPI.Transfer([]byte{uint8(r & RFM_REG_MAX), 0, 0}, buf[:]); err != nil {
		return 0, err
	}
	// Return success
	return uint16(buf[1])<<8 | uint16(buf[2]), nil
}

func (d *device) writereg_uint16(r register, data uint16) error {
	return d.SPI.Transfer([]byte{byte((r & RFM_REG_MAX) | RFM_REG_WRITE), uint8(data >> 8), uint8(data & 0xFF)}, nil)
}
