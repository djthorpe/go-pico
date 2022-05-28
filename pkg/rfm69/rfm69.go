package rfm69

import (

	// Package imports
	"fmt"

	pico "github.com/djthorpe/go-pico"
	spi "github.com/djthorpe/go-pico/pkg/spi"
	// Namespace imports
	//. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Bus   uint   // SPI Bus (0 or 1)
	Speed uint32 // SPI Communication Speed in Hz, optional
}

type device struct {
	spi     pico.SPI
	version uint8
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SPI_MODE    = 0
	SPI_SPEEDHZ = 115200 // Hz
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) New() (*device, error) {
	this := new(device)

	// Create SPI device
	device, err := spi.Config{
		Bus:       cfg.Bus,
		Frequency: SPI_SPEEDHZ,
		Mode:      SPI_MODE,
	}.New()
	if err != nil {
		return nil, err
	} else {
		this.spi = device
	}

	// Syncronize registers
	if err := this.sync(); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d *device) String() string {
	str := "<rfm69"
	str += fmt.Sprintf(" version=%02X", d.version)
	str += fmt.Sprint(" spi=", d.spi)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// SYNC

func (d *device) sync() error {
	if version, err := d.Version(); err != nil {
		return err
	} else {
		d.version = version
	}

	// Return success
	return nil
}
