package rfm69

import (
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/errors"
)

const (
	wait_step    = time.Millisecond * 50
	wait_retries = 10
)

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

// wait_for for a callback function to return true
// timeout after minimum of 500ms
func wait_for(callback func() (bool, error)) error {
	for i := 0; i < wait_retries; i++ {
		if r, err := callback(); err != nil {
			return err
		} else if r {
			return nil
		}
		time.Sleep(wait_step)
	}
	return ErrTimeout
}
