package bme280

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Mode uint8
type Filter uint8
type StandbyTime uint8
type Oversample uint8

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// BME280 Mode
const (
	BME280_MODE_SLEEP   Mode = 0x00
	BME280_MODE_FORCED  Mode = 0x01
	BME280_MODE_FORCED2 Mode = 0x02
	BME280_MODE_NORMAL  Mode = 0x03
	BME280_MODE_MAX     Mode = 0x03
)

// BME280 Filter Co-efficient
const (
	BME280_FILTER_OFF Filter = 0x00
	BME280_FILTER_2   Filter = 0x01
	BME280_FILTER_4   Filter = 0x02
	BME280_FILTER_8   Filter = 0x03
	BME280_FILTER_16  Filter = 0x04
	BME280_FILTER_MAX Filter = 0x07
)

// BME280 Standby time
const (
	BME280_STANDBY_0P5MS  StandbyTime = 0x00
	BME280_STANDBY_62P5MS StandbyTime = 0x01
	BME280_STANDBY_125MS  StandbyTime = 0x02
	BME280_STANDBY_250MS  StandbyTime = 0x03
	BME280_STANDBY_500MS  StandbyTime = 0x04
	BME280_STANDBY_1000MS StandbyTime = 0x05
	BME280_STANDBY_10MS   StandbyTime = 0x06
	BME280_STANDBY_20MS   StandbyTime = 0x07
	BME280_STANDBY_MAX    StandbyTime = 0x07
)

// BME280 Oversampling value
const (
	BME280_OVERSAMPLE_SKIP Oversample = 0x00
	BME280_OVERSAMPLE_1    Oversample = 0x01
	BME280_OVERSAMPLE_2    Oversample = 0x02
	BME280_OVERSAMPLE_4    Oversample = 0x03
	BME280_OVERSAMPLE_8    Oversample = 0x04
	BME280_OVERSAMPLE_16   Oversample = 0x05
	BME280_OVERSAMPLE_MAX  Oversample = 0x07
)

// Sealevel pressure approximation
// http://resource.npl.co.uk/pressure/pressure.html
const (
	BME280_PRESSURE_SEALEVEL float64 = 103090 // in Pascals
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (m Mode) String() string {
	switch m {
	case BME280_MODE_SLEEP:
		return "BME280_MODE_SLEEP"
	case BME280_MODE_FORCED:
		return "BME280_MODE_FORCED"
	case BME280_MODE_FORCED2:
		return "BME280_MODE_FORCED"
	case BME280_MODE_NORMAL:
		return "BME280_MODE_NORMAL"
	default:
		return "[?? Invalid BME280Mode value]"
	}
}

func (f Filter) String() string {
	switch f {
	case BME280_FILTER_OFF:
		return "BME280_FILTER_OFF"
	case BME280_FILTER_2:
		return "BME280_FILTER_2"
	case BME280_FILTER_4:
		return "BME280_FILTER_4"
	case BME280_FILTER_8:
		return "BME280_FILTER_8"
	case BME280_FILTER_16:
		return "BME280_FILTER_16"
	default:
		return "BME280_FILTER_16"
	}
}

func (t StandbyTime) String() string {
	switch t {
	case BME280_STANDBY_0P5MS:
		return "BME280_STANDBY_0P5MS"
	case BME280_STANDBY_62P5MS:
		return "BME280_STANDBY_62P5MS"
	case BME280_STANDBY_125MS:
		return "BME280_STANDBY_125MS"
	case BME280_STANDBY_250MS:
		return "BME280_STANDBY_250MS"
	case BME280_STANDBY_500MS:
		return "BME280_STANDBY_500MS"
	case BME280_STANDBY_1000MS:
		return "BME280_STANDBY_1000MS"
	case BME280_STANDBY_10MS:
		return "BME280_STANDBY_10MS"
	case BME280_STANDBY_20MS:
		return "BME280_STANDBY_20MS"
	default:
		return "[?? Invalid BME280Standby value]"
	}
}

func (o Oversample) String() string {
	switch o {
	case BME280_OVERSAMPLE_SKIP:
		return "BME280_OVERSAMPLE_SKIP"
	case BME280_OVERSAMPLE_1:
		return "BME280_OVERSAMPLE_1"
	case BME280_OVERSAMPLE_2:
		return "BME280_OVERSAMPLE_2"
	case BME280_OVERSAMPLE_4:
		return "BME280_OVERSAMPLE_4"
	case BME280_OVERSAMPLE_8:
		return "BME280_OVERSAMPLE_8"
	case BME280_OVERSAMPLE_16:
		return "BME280_OVERSAMPLE_16"
	default:
		return "[?? Invalid BME280Oversample value]"
	}
}
