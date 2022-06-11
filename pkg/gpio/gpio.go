package gpio

import (
	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	In    []Pin // Pins in input mode
	Out   []Pin // Pins in output mode
	PWM   []Pin // Output pins with PWM
	Watch []Pin // Input pins to watch
}

const (
	DEFAULT_PWM = 1000 // Default PWM frequency
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create new GPIO, with channel. Panics on error
func New(cfg Config, ch chan<- Event) *device {
	if d, err := cfg.New(ch); err != nil {
		panic(err)
	} else {
		return d
	}
}
