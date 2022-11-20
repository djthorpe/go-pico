package bme280

import (
	"fmt"

	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/errors"

	// Package imports
	//math32 "github.com/chewxy/math32"
	math32 "github.com/djthorpe/go-pico/pkg/math32"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Calibration Coefficients
type cal struct {
	t1 uint16
	t2 int16
	t3 int16
	p1 uint16
	p2 int16
	p3 int16
	p4 int16
	p5 int16
	p6 int16
	p7 int16
	p8 int16
	p9 int16
	h1 uint8
	h2 int16
	h3 uint8
	h4 int16
	h5 int16
	h6 int8
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c cal) String() string {
	str := "<cal"
	str += fmt.Sprint(" t1=", c.t1)
	str += fmt.Sprint(" t2=", c.t2)
	str += fmt.Sprint(" t3=", c.t3)
	str += fmt.Sprint(" p1=", c.p1)
	str += fmt.Sprint(" p2=", c.p2)
	str += fmt.Sprint(" p3=", c.p3)
	str += fmt.Sprint(" p4=", c.p4)
	str += fmt.Sprint(" p5=", c.p5)
	str += fmt.Sprint(" p6=", c.p6)
	str += fmt.Sprint(" p7=", c.p7)
	str += fmt.Sprint(" p8=", c.p8)
	str += fmt.Sprint(" p9=", c.p9)
	str += fmt.Sprint(" h1=", c.h1)
	str += fmt.Sprint(" h2=", c.h2)
	str += fmt.Sprint(" h3=", c.h3)
	str += fmt.Sprint(" h4=", c.h4)
	str += fmt.Sprint(" h5=", c.h5)
	str += fmt.Sprint(" h6=", c.h6)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Read and return calibration coefficients
func (d *device) calibrate() (cal, error) {
	var data [24]byte
	var v cal

	// Bulk read T1-T3, P1-P9
	if err := d.readRegister(BME280_REG_DIG_T1, data[:]); err != nil {
		return v, err
	} else {
		v.t1 = toUint16(data[0:])
		v.t2 = toInt16(data[2:])
		v.t3 = toInt16(data[4:])
		v.p1 = toUint16(data[6:])
		v.p2 = toInt16(data[8:])
		v.p3 = toInt16(data[10:])
		v.p4 = toInt16(data[12:])
		v.p5 = toInt16(data[14:])
		v.p6 = toInt16(data[16:])
		v.p7 = toInt16(data[18:])
		v.p8 = toInt16(data[20:])
		v.p9 = toInt16(data[22:])
	}

	// Read H1
	if err := d.readRegister(BME280_REG_DIG_H1, data[0:1]); err != nil {
		return v, err
	} else {
		v.h1 = data[0]
	}

	// Read H2-H6
	if err := d.readRegister(BME280_REG_DIG_H2, data[0:9]); err != nil {
		return v, err
	} else {
		v.h2 = toInt16(data[0:])
		v.h3 = data[2]
		v.h4 = int16(data[3])<<4 | int16(data[4])&0x0F
		v.h5 = int16(data[5])<<4 | int16(data[4])>>4
		v.h6 = int8(data[6])
	}

	// return success
	return v, nil
}

// Test data from the Bosch datasheet. Used to check calculations.
// T1=27504 : T2=26435 : T3=-1000
// P1=36477 : P2=-10685 : P3=3024
// P4=2855 : P5=140 : P6=-7
// P7=15500 : P8=-14600 : P9=6000

// Test data from another sensor
// calibration=<calibration{
// T1=28244 T2=26571 T3=50
// P1=37759 P2=-10679 P3=3024 P4=8281 P5=-140 P6=-7 P7=9900 P8=-10230 P9=4285
// H1=75 H2=353 H3=0 H4=340 H5=0 H6=30

// toTemperature returns tvalue (milli-degrees) and tfine. ErrSampleSkipped is returned
// if temperature sample was not available
func toTemperature(data [8]byte, coefficients cal) (int32, int32, error) {
	raw := ((int32(data[3]) << 16) | (int32(data[4]) << 8) | int32(data[5])) >> 4
	if raw == BME280_SKIPTEMP_VALUE {
		return 0, 0, ErrSampleSkipped
	}
	var1 := (((raw >> 3) - (int32(coefficients.t1) << 1)) * int32(coefficients.t2)) >> 11
	var2 := (((((raw >> 4) - int32(coefficients.t1)) * ((raw >> 4) - int32(coefficients.t1))) >> 12) * int32(coefficients.t3)) >> 14
	tFine := var1 + var2
	T := (tFine*5 + 128) >> 8
	return (10 * T), tFine, nil
}

// toPressure returns pvalue (milli-pascals). ErrSampleSkipped is returned if
// pressure sample was not available
func toPressure(data [8]byte, tfine int32, coefficients cal) (int32, error) {
	raw := ((int32(data[0]) << 16) | (int32(data[1]) << 8) | int32(data[2])) >> 4
	if raw == BME280_SKIPPRESSURE_VALUE {
		return 0, ErrSampleSkipped
	}
	var1 := int64(tfine) - 128000
	var2 := var1 * var1 * int64(coefficients.p6)
	var2 = var2 + ((var1 * int64(coefficients.p5)) << 17)
	var2 = var2 + (int64(coefficients.p4) << 35)
	var1 = ((var1 * var1 * int64(coefficients.p3)) >> 8) + ((var1 * int64(coefficients.p2)) << 12)
	var1 = ((int64(1) << 47) + var1) * int64(coefficients.p1) >> 33

	if var1 == 0 {
		// avoid exception caused by division by zero
		return 0, ErrSampleSkipped
	}
	p := int64(1048576 - raw)
	p = (((p << 31) - var2) * 3125) / var1
	var1 = (int64(coefficients.p9) * (p >> 13) * (p >> 13)) >> 25
	var2 = (int64(coefficients.p8) * p) >> 19

	p = ((p + var1 + var2) >> 8) + (int64(coefficients.p7) << 4)
	p = (p / 256)
	return int32(1000 * p), nil
}

// toHumidity returns hvalue which is relative humidity in hundredths of a percent.
// ErrSampleSkipped is returned if humidity is not available
func toHumidity(data [8]byte, tfine int32, coefficients cal) (int32, error) {
	raw := int32(uint32(data[6])<<8 | uint32(data[7]))
	if raw == BME280_SKIPHUMID_VALUE {
		return 0, ErrSampleSkipped
	}

	// Offset tfine
	h := float32(tfine) - 76800.0
	if h == 0 {
		return 0, ErrSampleSkipped
	}

	// Calibrate
	var1 := float32(raw) - (float32(coefficients.h4)*64.0 +
		(float32(coefficients.h5) / 16384.0 * h))

	var2 := float32(coefficients.h2) / 65536.0 *
		(1.0 + float32(coefficients.h6)/67108864.0*h*
			(1.0+float32(coefficients.h3)/67108864.0*h))

	h = var1 * var2
	h = h * (1 - float32(coefficients.h1)*h/524288)
	return int32(100 * h), nil
}

const (
	pow  = 0.19022256      // 1.0 / 5.257
	koff = 273.15 * 1000.0 // Kelvin offset
	pdiv = 0.0065
)

// toAltitude returns height in centimetres based on temperature, current pressure and sealevel pressure
func toAltitude(tvalue, pvalue, p0 int32) (int32, error) {
	if pvalue == 0 {
		return 0, ErrSampleSkipped
	}
	// Hypsometric formula
	// https://keisan.casio.com/exec/system/1224585971
	r := float32(p0) / float32(pvalue)
	k := float32(tvalue) + koff
	alt := ((math32.Pow(r, pow) - 1.0) * k) / (pdiv * 10.0)
	return int32(alt), nil
}
