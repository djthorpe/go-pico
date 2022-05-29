package gpio

import (
	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	In  []Pin
	Out []Pin
}
