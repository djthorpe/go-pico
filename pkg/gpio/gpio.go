package gpio

import (
	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	In    []Pin
	Out   []Pin
	Watch []Pin
}

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
