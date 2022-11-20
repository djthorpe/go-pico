package bme280

import "time"

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func toOversampleNumber(value Oversample) float64 {
	switch value {
	case BME280_OVERSAMPLE_SKIP:
		return 0
	case BME280_OVERSAMPLE_1:
		return 1
	case BME280_OVERSAMPLE_2:
		return 2
	case BME280_OVERSAMPLE_4:
		return 4
	case BME280_OVERSAMPLE_8:
		return 8
	case BME280_OVERSAMPLE_16:
		return 16
	default:
		return 0
	}
}

func toMeasurementTime(osrs_t, osrs_p, osrs_h Oversample) time.Duration {
	// Measurement Time as per BME280 datasheet section 9.1
	time_ms := 1.25
	if osrs_t != BME280_OVERSAMPLE_SKIP {
		time_ms += toOversampleNumber(osrs_t) * 2.3
	}
	if osrs_p != BME280_OVERSAMPLE_SKIP {
		time_ms += toOversampleNumber(osrs_p)*2.3 + 0.575
	}
	if osrs_h != BME280_OVERSAMPLE_SKIP {
		time_ms += toOversampleNumber(osrs_h)*2.4 + 0.575
	}
	return time.Millisecond * time.Duration(time_ms)
}

func toStandbyTime(value StandbyTime) time.Duration {
	switch value {
	case BME280_STANDBY_0P5MS:
		return time.Microsecond * 500
	case BME280_STANDBY_62P5MS:
		return time.Microsecond * 62500
	case BME280_STANDBY_125MS:
		return time.Millisecond * 125
	case BME280_STANDBY_250MS:
		return time.Millisecond * 250
	case BME280_STANDBY_500MS:
		return time.Millisecond * 500
	case BME280_STANDBY_1000MS:
		return time.Millisecond * 1000
	case BME280_STANDBY_10MS:
		return time.Millisecond * 10
	case BME280_STANDBY_20MS:
		return time.Millisecond * 20
	default:
		return 0
	}
}

func toUint16(v []byte) uint16 {
	return uint16(v[1])<<8 | uint16(v[0])
}

func toInt16(v []byte) int16 {
	return int16(v[1])<<8 | int16(v[0])
}
