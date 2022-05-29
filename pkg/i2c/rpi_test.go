// go:build rpi

package i2c_test

import (
	"testing"

	// Namespace import
	. "github.com/djthorpe/go-pico/pkg/i2c"
)

func Test_I2C_001(t *testing.T) {
	for bus := uint(0); bus <= 2; bus++ {
		i2c, err := Config{Bus: bus}.New()
		if err != nil {
			t.Log("Not found bus:", bus, ":", err)
		} else if err := i2c.Close(); err != nil {
			t.Error(err)
		} else {
			t.Log(i2c)
		}
	}
}
