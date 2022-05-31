//go:build !production

package rfm69

import "fmt"

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - UINT8

func (d *device) readreg_uint8(r register) (uint8, error) {
	var buf [2]uint8
	fmt.Print("readreg_uint8(", r, ")")
	if err := d.SPI.Transfer([]byte{uint8(r & RFM_REG_MAX), 0}, buf[:]); err != nil {
		fmt.Println("  err:", err)
		return 0, err
	}
	fmt.Printf(" => 0x%02X\n", buf[1])
	// Return success
	return buf[1], nil
}

func (d *device) writereg_uint8(r register, data uint8) error {
	fmt.Printf("writereg_uint8(%v, 0x%02X)", r, data)
	if err := d.SPI.Transfer([]byte{byte((r & RFM_REG_MAX) | RFM_REG_WRITE), data}, nil); err != nil {
		fmt.Println("  err:", err)
		return err
	} else {
		fmt.Println(" OK")
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - UINT16

func (d *device) readreg_uint16(r register) (uint16, error) {
	var buf [3]uint8
	fmt.Print("readreg_uint16(", r, ")")
	if err := d.SPI.Transfer([]byte{uint8(r & RFM_REG_MAX), 0, 0}, buf[:]); err != nil {
		fmt.Println("  err:", err)
		return 0, err
	}
	v := uint16(buf[1])<<8 | uint16(buf[2])
	fmt.Printf(" => 0x%04X\n", v)
	// Return success
	return v, nil
}

func (d *device) writereg_uint16(r register, data uint16) error {
	fmt.Printf("writereg_uint16(%v, 0x%04X)", r, data)
	if err := d.SPI.Transfer([]byte{byte((r & RFM_REG_MAX) | RFM_REG_WRITE), uint8(data >> 8), uint8(data & 0xFF)}, nil); err != nil {
		fmt.Println("  err:", err)
		return err
	} else {
		fmt.Println(" OK")
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - DEBUGGING

func println(v ...interface{}) {
	fmt.Println(v...)
}
