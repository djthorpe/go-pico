package rfm69

type (
	Mode       uint8
	DataMode   uint8
	Modulation uint8
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
