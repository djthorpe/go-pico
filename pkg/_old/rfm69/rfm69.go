package rfm69

import (
	"fmt"

	// Package imports
	spi "github.com/djthorpe/go-pico/pkg/spi"
	multierror "github.com/hashicorp/go-multierror"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
	. "github.com/djthorpe/go-pico/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Bus   uint   // SPI Bus (0 or 1)
	Slave uint   // SPI Slave (0 or 1) not used on the Pico
	Speed uint32 // SPI Communication Speed in Hz, optional
}

type device struct {
	SPI
	version       uint8
	mode          Mode
	sequencer_off bool
	listen_on     bool
	data_mode     DataMode
	modulation    Modulation
	bitrate       uint16
	/*
		frf                   uint32
		fdev                  uint16
		aes_key               []byte
		aes_on                bool
		sync_word             []byte
		sync_on               bool
		sync_size             uint8
		sync_tol              uint8
		rx_inter_packet_delay uint8
		rx_auto_restart       bool
		tx_start              TXStart
		fifo_threshold        uint8
		fifo_fill_condition   bool
		node_address          uint8
		broadcast_address     uint8
		preamble_size         uint16
		payload_size          uint8
		packet_format         PacketFormat
		packet_coding         PacketCoding
		packet_filter         PacketFilter
		crc_enabled           bool
		crc_auto_clear_off    bool
		afc                   int16
		afc_mode              AFCMode
		afc_routine           AFCRoutine
		lna_impedance         LNAImpedance
		lna_gain              LNAGain
		rxbw_frequency        RXBWFrequency
		rxbw_cutoff           RXBWCutoff
	*/
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SPI_MODE        = 0
	SPI_SPEEDHZ     = 115200     // Hz
	RFM_FXOSC_MHZ   = 32         // Crystal oscillator frequency MHz
	RFM_FSTEP_HZ    = 61         // Frequency synthesizer step
	RFM_BITRATE_MIN = 500        // bits per second (Hz)
	RFM_BITRATE_MAX = 300 * 1024 // bits per second (Hz)
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) New() (*device, error) {
	this := new(device)

	// Create SPI device
	spiconfig := spi.Config{
		Bus:   cfg.Bus,
		Slave: cfg.Slave,
		Speed: cfg.Speed | SPI_SPEEDHZ,
		Mode:  SPI_MODE,
	}
	if device, err := spiconfig.New(); err != nil {
		return nil, err
	} else {
		this.SPI = device
	}

	// Syncronize registers
	if err := this.sync(); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

func (d *device) Close() error {
	var result error

	if d.SPI != nil {
		if err := d.SPI.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	d.SPI = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d *device) String() string {
	str := "<rfm69"
	str += fmt.Sprintf(" version=0x%02X", d.version)
	str += fmt.Sprint(" mode=", d.mode)
	str += fmt.Sprint(" data_mode=", d.data_mode)
	str += fmt.Sprint(" modulation=", d.modulation)
	str += fmt.Sprintf(" bitrate=0x%04X bitrate_hz=%v", d.bitrate, to_bitratehz(d.bitrate))

	if d.listen_on {
		str += " listen_on"
	}
	if d.sequencer_off {
		str += " sequencer_off="
	}
	if d.SPI != nil {
		str += fmt.Sprint(" spi=", d.SPI)
	}

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

	// Get operational mode
	if mode, listen_on, sequencer_off, err := d.getOpMode(); err != nil {
		return err
	} else {
		d.mode = mode
		d.listen_on = listen_on
		d.sequencer_off = sequencer_off
	}

	// Get data mode and modulation
	if data_mode, modulation, err := d.getDataModul(); err != nil {
		return err
	} else {
		d.data_mode = data_mode
		d.modulation = modulation
	}

	// Get bitrate
	if bitrate, err := d.getBitrate(); err != nil {
		return err
	} else {
		d.bitrate = bitrate
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MODE, DATA MODE AND MODULATION

// Return device mode
func (d *device) Mode() Mode {
	return d.mode
}

// Return data mode
func (d *device) DataMode() DataMode {
	return d.data_mode
}

// Return modulation
func (d *device) Modulation() Modulation {
	return d.modulation
}

// Set device mode
func (d *device) SetMode(v Mode) error {
	// Unset listen_on if not going to sleep
	if v != RFM_MODE_SLEEP && d.listen_on {
		if err := d.setOpMode(v, false, true, d.sequencer_off); err != nil {
			return err
		} else {
			d.listen_on = false
		}
	}

	// Unset sequencer if going into standby
	if v == RFM_MODE_STDBY && d.sequencer_off == false {
		if err := d.setOpMode(d.mode, d.listen_on, false, true); err != nil {
			return err
		} else {
			d.sequencer_off = true
		}
	}

	// Write mode and read back again
	if err := d.setOpMode(v, false, false, d.sequencer_off); err != nil {
		return err
	}

	// Wait for device ready bit
	if err := wait_for(func() (bool, error) {
		value, err := d.getIRQ1(RFM_IRQ1_MODEREADY)
		return to_uint8_bool(uint8(value)), err
	}); err != nil {
		return err
	}

	// Wait for RX or TX Ready
	if v == RFM_MODE_RX {
		if err := wait_for(func() (bool, error) {
			value, err := d.getIRQ1(RFM_IRQ1_RXREADY)
			return to_uint8_bool(uint8(value)), err
		}); err != nil {
			return err
		}
	} else if v == RFM_MODE_TX {
		if err := wait_for(func() (bool, error) {
			value, err := d.getIRQ1(RFM_IRQ1_TXREADY)
			return to_uint8_bool(uint8(value)), err
		}); err != nil {
			return err
		}
	}

	// Read back register
	if mode, listen_on, sequencer_off, err := d.getOpMode(); err != nil {
		return err
	} else if mode != v {
		return ErrUnexpectedValue.With("SetMode")
	} else if listen_on != d.listen_on {
		return ErrUnexpectedValue.With("SetMode")
	} else {
		d.mode = mode
		d.listen_on = listen_on
		d.sequencer_off = sequencer_off
	}

	// If RX mode then read AFC value
	// TODO
	/*if this.mode == sensors.RFM_MODE_RX {
		if afc, err := this.getAFC(); err != nil {
			return err
		} else {
			this.afc = afc
		}
	}*/

	// Return success
	return nil
}

// Set data mode
func (d *device) SetDataMode(v DataMode) error {
	// Write
	if err := d.setDataModul(v, d.modulation); err != nil {
		return err
	}

	// Read
	if data_mode, modulation, err := d.getDataModul(); err != nil {
		return err
	} else if data_mode != v {
		return ErrUnexpectedValue.With("SetDataMode")
	} else if modulation != d.modulation {
		return ErrUnexpectedValue.With("SetDataMode")
	} else {
		d.data_mode = data_mode
	}

	// Return success
	return nil
}

// Set modulation
func (d *device) SetModulation(v Modulation) error {
	// Write
	if err := d.setDataModul(d.data_mode, v); err != nil {
		return err
	}

	// Read
	if data_mode, modulation, err := d.getDataModul(); err != nil {
		return err
	} else if modulation != v {
		return ErrUnexpectedValue.With("SetModulation")
	} else if data_mode != d.data_mode {
		return ErrUnexpectedValue.With("SetModulation")
	} else {
		d.modulation = modulation
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BITRATE

// Return bitrate in bits per second
func (d *device) Bitrate() uint16 {
	return d.bitrate
}
