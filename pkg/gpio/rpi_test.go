// go:build rpi

package gpio_test

import (
	"testing"

	// Namespace import
	. "github.com/djthorpe/go-pico/pkg/gpio"
)

func Test_GPIO_001(t *testing.T) {
	gpio, err := Config{}.New()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(gpio)
		if err := gpio.Close(); err != nil {
			t.Error(err)
		}
	}
}
