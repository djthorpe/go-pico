package spi

import (
	"fmt"
	"machine"

	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Bus       uint
	Frequency uint32
	Mode      Mode
}

type device struct {
	*machine.SPI
	sck, sdo, sdi, scs machine.Pin
}

type Mode uint8

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	CPOLCPHA Mode = 0b11
	CPOL     Mode = 0b10
	CPHA     Mode = 0b01
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) New() (*device, error) {
	this := new(device)

	// Pico has two SPI devices
	switch cfg.Bus {
	case 0:
		this.SPI = machine.SPI0
		this.sck = machine.SPI0_SCK_PIN
		this.sdo = machine.SPI0_SDO_PIN
		this.sdi = machine.SPI0_SDI_PIN
		this.scs = machine.Pin(17)
	case 1:
		this.SPI = machine.SPI1
		this.sck = machine.SPI1_SCK_PIN
		this.sdo = machine.SPI1_SDO_PIN
		this.sdi = machine.SPI1_SDI_PIN
		this.scs = machine.Pin(13)
	default:
		return nil, ErrBadParameter
	}

	// Initialise the device
	if err := this.SPI.Configure(machine.SPIConfig{
		Frequency: cfg.Frequency,
		Mode:      uint8(cfg.Mode),
		DataBits:  8,
		SCK:       this.sck,
		SDO:       this.sdo,
		SDI:       this.sdi,
	}); err != nil {
		return nil, err
	}

	// Configure the CS pin
	this.scs.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	this.scs.High()

	// Return success
	return this, nil
}

func (d *device) Close() error {
	// No close implementation on pico
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d *device) String() string {
	str := "<spi"
	str += fmt.Sprint(" sck=", d.sck)
	str += fmt.Sprint(" sdo=", d.sdo)
	str += fmt.Sprint(" sdi=", d.sdi)
	str += fmt.Sprint(" scs=", d.scs)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Transfer writes then reads from SPI bus
func (d *device) Transfer(w, r []byte) error {
	//fmt.Println("tx=>", hex.EncodeToString(w))
	d.scs.Low()
	err := d.SPI.Tx(w, r)
	d.scs.High()
	/*if err == nil {
		fmt.Println("  <=", hex.EncodeToString(r))
	} else {
		fmt.Println("  err: ", err)
	}*/
	return err
}
