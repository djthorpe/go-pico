package rfm69

import "strings"

type (
	Mode       uint8
	DataMode   uint8
	Modulation uint8
	IRQ1       uint8
	IRQ2       uint8
)

const (
	RFM_MODE_SLEEP Mode = 0x00
	RFM_MODE_STDBY Mode = 0x01
	RFM_MODE_FS    Mode = 0x02
	RFM_MODE_TX    Mode = 0x03
	RFM_MODE_RX    Mode = 0x04
	RFM_MODE_MAX   Mode = 0x07
)

const (
	RFM_DATAMODE_PACKET            DataMode = 0x00
	RFM_DATAMODE_CONTINUOUS_NOSYNC DataMode = 0x02
	RFM_DATAMODE_CONTINUOUS_SYNC   DataMode = 0x03
	RFM_DATAMODE_MAX               DataMode = 0x03
)

const (
	RFM_MODULATION_FSK        Modulation = 0x00 // 00000 FSK no shaping
	RFM_MODULATION_FSK_BT_1P0 Modulation = 0x01 // 01000 FSK Guassian filter, BT=1.0
	RFM_MODULATION_FSK_BT_0P5 Modulation = 0x02 // 10000 FSK Gaussian filter, BT=0.5
	RFM_MODULATION_FSK_BT_0P3 Modulation = 0x03 // 11000 FSK Gaussian filter, BT=0.3
	RFM_MODULATION_OOK        Modulation = 0x08 // 00001 OOK no shaping
	RFM_MODULATION_OOK_BR     Modulation = 0x09 // 01001 OOK Filtering with f(cutoff) = BR
	RFM_MODULATION_OOK_2BR    Modulation = 0x0A // 01010 OOK Filtering with f(cutoff) = 2BR
	RFM_MODULATION_MAX        Modulation = 0x0A
)

const (
	RFM_IRQ1_MODEREADY        IRQ1 = 0x80 // Mode has changed
	RFM_IRQ1_RXREADY          IRQ1 = 0x40
	RFM_IRQ1_TXREADY          IRQ1 = 0x20
	RFM_IRQ1_PLLLOCK          IRQ1 = 0x10
	RFM_IRQ1_RSSI             IRQ1 = 0x08
	RFM_IRQ1_TIMEOUT          IRQ1 = 0x04
	RFM_IRQ1_AUTOMODE         IRQ1 = 0x02
	RFM_IRQ1_SYNCADDRESSMATCH IRQ1 = 0x01
	RFM_IRQ1_MAX              IRQ1 = 0x80
	RFM_IRQ1_NONE             IRQ1 = 0x00
)

const (
	RFM_IRQ2_CRCOK        IRQ2 = 0x02
	RFM_IRQ2_PAYLOADREADY IRQ2 = 0x04
	RFM_IRQ2_PACKETSENT   IRQ2 = 0x08
	RFM_IRQ2_FIFOOVERRUN  IRQ2 = 0x10
	RFM_IRQ2_FIFOLEVEL    IRQ2 = 0x20
	RFM_IRQ2_FIFONOTEMPTY IRQ2 = 0x40
	RFM_IRQ2_FIFOFULL     IRQ2 = 0x80
	RFM_IRQ2_MAX          IRQ2 = 0x80
	RFM_IRQ2_NONE         IRQ2 = 0x00
)

func (m Mode) String() string {
	switch m {
	case RFM_MODE_SLEEP:
		return "RFM_MODE_SLEEP"
	case RFM_MODE_STDBY:
		return "RFM_MODE_STDBY"
	case RFM_MODE_FS:
		return "RFM_MODE_FS"
	case RFM_MODE_TX:
		return "RFM_MODE_TX"
	case RFM_MODE_RX:
		return "RFM_MODE_RX"
	default:
		return "[??]"
	}
}

func (m DataMode) String() string {
	switch m {
	case RFM_DATAMODE_PACKET:
		return "RFM_DATAMODE_PACKET"
	case RFM_DATAMODE_CONTINUOUS_NOSYNC:
		return "RFM_DATAMODE_CONTINUOUS_NOSYNC"
	case RFM_DATAMODE_CONTINUOUS_SYNC:
		return "RFM_DATAMODE_CONTINUOUS_SYNC"
	default:
		return "[??]"
	}
}

func (m Modulation) String() string {
	switch m {
	case RFM_MODULATION_FSK:
		return "RFM_MODULATION_FSK"
	case RFM_MODULATION_FSK_BT_1P0:
		return "RFM_MODULATION_FSK_BT_1P0"
	case RFM_MODULATION_FSK_BT_0P5:
		return "RFM_MODULATION_FSK_BT_0P5"
	case RFM_MODULATION_FSK_BT_0P3:
		return "RFM_MODULATION_FSK_BT_0P3"
	case RFM_MODULATION_OOK:
		return "RFM_MODULATION_OOK"
	case RFM_MODULATION_OOK_BR:
		return "RFM_MODULATION_OOK_BR"
	case RFM_MODULATION_OOK_2BR:
		return "RFM_MODULATION_OOK_2BR"
	default:
		return "[??]"
	}
}

func (f IRQ1) String() string {
	if f == RFM_IRQ1_NONE {
		return f.flagstring()
	}
	str := ""
	for v := IRQ1(1); v <= RFM_IRQ1_MAX; v <<= 1 {
		if v&f == v {
			str += "|" + v.flagstring()
		}
	}
	return strings.TrimPrefix(str, "|")
}

func (v IRQ1) flagstring() string {
	switch v {
	case RFM_IRQ1_NONE:
		return "NONE"
	case RFM_IRQ1_MODEREADY:
		return "MODEREADY"
	case RFM_IRQ1_RXREADY:
		return "RXREADY"
	case RFM_IRQ1_TXREADY:
		return "TXREADY"
	case RFM_IRQ1_PLLLOCK:
		return "PLLLOCK"
	case RFM_IRQ1_RSSI:
		return "RSSI"
	case RFM_IRQ1_TIMEOUT:
		return "TIMEOUT"
	case RFM_IRQ1_AUTOMODE:
		return "AUTOMODE"
	case RFM_IRQ1_SYNCADDRESSMATCH:
		return "SYNCADDRESSMATCH"
	default:
		return "[??]"
	}
}

func (f IRQ2) String() string {
	if f == RFM_IRQ2_NONE {
		return f.flagstring()
	}
	str := ""
	for v := IRQ2(1); v <= RFM_IRQ2_MAX; v <<= 1 {
		if v&f == v {
			str += "|" + v.flagstring()
		}
	}
	return strings.TrimPrefix(str, "|")
}

func (v IRQ2) flagstring() string {
	switch v {
	case RFM_IRQ2_NONE:
		return "NONE"
	case RFM_IRQ2_CRCOK:
		return "CRCOK"
	case RFM_IRQ2_PAYLOADREADY:
		return "PAYLOADREADY"
	case RFM_IRQ2_PACKETSENT:
		return "PACKETSENT"
	case RFM_IRQ2_FIFOOVERRUN:
		return "FIFOOVERRUN"
	case RFM_IRQ2_FIFOLEVEL:
		return "FIFOLEVEL"
	case RFM_IRQ2_FIFONOTEMPTY:
		return "FIFONOTEMPTY"
	case RFM_IRQ2_FIFOFULL:
		return "FIFOFULL"
	default:
		return "[??]"
	}
}
