package uart

import (
	"fmt"
	"machine"

	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/errors"
)

type Config struct {
	Bus                uint
	BaudRate           uint32
	DataBits, StopBits uint8
}

type device struct {
	*machine.UART
	tx, rx machine.Pin
}

func (cfg Config) New() (*device, error) {
	this := new(device)

	// One supported UART device
	switch cfg.Bus {
	case 0:
		this.UART = machine.UART0
		this.tx = machine.UART0_TX_PIN
		this.rx = machine.UART0_RX_PIN
	default:
		return nil, ErrBadParameter
	}

	if err := this.UART.Configure(machine.UARTConfig{
		BaudRate: cfg.BaudRate,
		TX:       this.tx,
		RX:       this.rx,
	}); err != nil {
		return nil, err
	}
	if cfg.DataBits != 0 {
		if err := this.UART.SetFormat(cfg.DataBits, cfg.StopBits, machine.ParityNone); err != nil {
			return nil, err
		}
	}

	// Return success
	return this, nil
}

// Printf writes formatted string to UART
func (d *device) Printf(v string, args ...interface{}) {
	d.UART.Write([]byte(fmt.Sprintf(v, args...)))
}

// Print writes string to UART
func (d *device) Print(args ...interface{}) {
	d.UART.Write([]byte(fmt.Sprint(args...)))
}

// Println writes string to UART
func (d *device) Println(args ...interface{}) {
	d.UART.Write([]byte(fmt.Sprintln(args...)))
}
