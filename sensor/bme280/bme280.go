/*
 * Copyright (c) Clinton Freeman 2016
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
 * associated documentation files (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge, publish, distribute,
 * sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or
 * substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
 * NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

// Package bme280 allows interfacing with Bosch BME280 digital humidity, pressure
// and temperature sensor.
package bme280

import (
	"github.com/aquarat/embd"
)

const (
	CAL_T1_LSB_REG = 0x88
	CAL_T1_MSB_REG = 0x89
	CAL_T2_LSB_REG = 0x8A
	CAL_T2_MSB_REG = 0x8B
	CAL_T3_LSB_REG = 0x8C
	CAL_T3_MSB_REG = 0x8D

	CAL_P1_LSB_REG = 0x8E
	CAL_P1_MSB_REG = 0x8F
	CAL_P2_LSB_REG = 0x90
	CAL_P2_MSB_REG = 0x91
	CAL_P3_LSB_REG = 0x92
	CAL_P3_MSB_REG = 0x93
	CAL_P4_LSB_REG = 0x94
	CAL_P4_MSB_REG = 0x95
	CAL_P5_LSB_REG = 0x96
	CAL_P5_MSB_REG = 0x97
	CAL_P6_LSB_REG = 0x98
	CAL_P6_MSB_REG = 0x99
	CAL_P7_LSB_REG = 0x9A
	CAL_P7_MSB_REG = 0x9B
	CAL_P8_LSB_REG = 0x9C
	CAL_P8_MSB_REG = 0x9D
	CAL_P9_LSB_REG = 0x9E
	CAL_P9_MSB_REG = 0x9F

	CAL_H1_REG     = 0xA1
	CAL_H2_LSB_REG = 0xE1
	CAL_H2_MSB_REG = 0xE2
	CAL_H3_REG     = 0xE3
	CAL_H4_MSB_REG = 0xE4
	CAL_H4_LSB_REG = 0xE5
	CAL_H5_MSB_REG = 0xE6
	CAL_H6_REG     = 0xE7

	MEASUREMENT_START_REG = 0xF7
	PRESSURE_MSB_IDX      = 0
	PRESSURE_LSB_IDX      = 1
	PRESSURE_XLSB_IDX     = 2

	TMP_MSB_IDX  = 3
	TMP_LSB_IDX  = 4
	TMP_XLSB_IDX = 5

	HUMIDITY_MSB_IDX = 6
	HUMIDITY_LSB_IDX = 7

	CTRL_MEAS_REG     = 0xF4
	CONFIG_REG        = 0xF5
	CTRL_HUMIDITY_REG = 0xF2
	RESET_REG         = 0xE0

	//RunMode can be:
	//  0, Sleep mode
	//  1 or 2, Forced mode
	//  3, Normal mode
	RunMode = uint8(3)

	//Standby can be:
	//  0, 0.5ms
	//  1, 62.5ms
	//  2, 125ms
	//  3, 250ms
	//  4, 500ms
	//  5, 1000ms
	//  6, 10ms
	//  7, 20ms
	Standby = uint8(0)

	//Filter can be off or number of FIR coefficients to use:
	//  0, filter off
	//  1, coefficients = 2
	//  2, coefficients = 4
	//  3, coefficients = 8
	//  4, coefficients = 16
	Filter = uint8(0)

	//TempOverSample can be:
	//  0, skipped
	//  1 through 5, oversampling *1, *2, *4, *8, *16 respectivel
	TempOverSample = uint8(1)

	//PressOverSample can be:
	//  0, skipped
	//  1 through 5, oversampling *1, *2, *4, *8, *16 respectively
	PressOverSample = uint8(1)

	//HumidOverSample can be:
	//  0, skipped
	//  1 through 5, oversampling *1, *2, *4, *8, *16 respectively
	HumidOverSample = uint8(1)
)

type Calibration struct {
	T1     uint16
	T2, T3 int16

	P1, P2, P3, P4, P5, P6, P7, P8, P9 int64
	H1, H2, H3, H4, H5, H6             float64
}

type BME280 struct {
	Bus  embd.I2CBus
	Addr byte

	Cal Calibration
}

func readUInt16(lsb byte, msb byte, bus embd.I2CBus, addr byte) (uint16, error) {
	lsbv, err := bus.ReadByteFromReg(addr, lsb)
	if err != nil {
		return 0, err
	}
	msbv, err := bus.ReadByteFromReg(addr, msb)
	if err != nil {
		return 0, err
	}

	return (uint16(msbv) << 8) | uint16(lsbv), nil
}

func readInt16(lsb byte, msb byte, bus embd.I2CBus, addr byte) (int16, error) {
	lsbv, err := bus.ReadByteFromReg(addr, lsb)
	if err != nil {
		return 0, err
	}
	msbv, err := bus.ReadByteFromReg(addr, msb)
	if err != nil {
		return 0, err
	}

	return (int16(msbv) << 8) | int16(lsbv), nil
}

// New creates and calibrates a connection to a BME280 sensor on the supplied i2c bus
// at the nominated i2c address.
func New(bus embd.I2CBus, addr byte) (*BME280, error) {
	s := &BME280{
		Bus:  bus,
		Addr: addr,
	}

	var err error
	var msb, lsb byte

	// Get calibrate information.
	s.Cal.T1, err = readUInt16(CAL_T1_LSB_REG, CAL_T1_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}

	s.Cal.T2, err = readInt16(CAL_T2_LSB_REG, CAL_T2_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}

	s.Cal.T3, err = readInt16(CAL_T3_LSB_REG, CAL_T3_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}

	pu, err := readUInt16(CAL_P1_LSB_REG, CAL_P1_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}
	s.Cal.P1 = int64(pu)

	p, err := readInt16(CAL_P2_LSB_REG, CAL_P2_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}
	s.Cal.P2 = int64(p)

	p, err = readInt16(CAL_P3_LSB_REG, CAL_P3_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}
	s.Cal.P3 = int64(p)

	p, err = readInt16(CAL_P4_LSB_REG, CAL_P4_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}
	s.Cal.P4 = int64(p)

	p, err = readInt16(CAL_P5_LSB_REG, CAL_P5_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}
	s.Cal.P5 = int64(p)

	p, err = readInt16(CAL_P6_LSB_REG, CAL_P6_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}
	s.Cal.P6 = int64(p)

	p, err = readInt16(CAL_P7_LSB_REG, CAL_P7_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}
	s.Cal.P7 = int64(p)

	p, err = readInt16(CAL_P8_LSB_REG, CAL_P8_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}
	s.Cal.P8 = int64(p)

	p, err = readInt16(CAL_P9_LSB_REG, CAL_P9_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}
	s.Cal.P9 = int64(p)

	msb, err = bus.ReadByteFromReg(addr, CAL_H1_REG)
	if err != nil {
		return s, err
	}
	s.Cal.H1 = float64(msb)

	h2, err := readInt16(CAL_H2_LSB_REG, CAL_H2_MSB_REG, bus, addr)
	if err != nil {
		return s, err
	}
	s.Cal.H2 = float64(h2)

	msb, err = bus.ReadByteFromReg(addr, CAL_H3_REG)
	if err != nil {
		return s, err
	}
	s.Cal.H3 = float64(msb)

	// H4 and H5 share three registers.
	msb, err = bus.ReadByteFromReg(addr, CAL_H4_MSB_REG)
	if err != nil {
		return s, err
	}
	lsb, err = bus.ReadByteFromReg(addr, CAL_H4_LSB_REG)
	if err != nil {
		return s, err
	}
	s.Cal.H4 = float64((int16(msb) << 4) | int16(lsb)&0x0F)

	msb, err = bus.ReadByteFromReg(addr, CAL_H5_MSB_REG)
	if err != nil {
		return s, err
	}
	lsb, err = bus.ReadByteFromReg(addr, CAL_H4_LSB_REG)
	if err != nil {
		return s, err
	}
	s.Cal.H5 = float64((int16(msb) << 4) | (int16(lsb)>>4)&0x0F)

	msb, err = bus.ReadByteFromReg(addr, CAL_H6_REG)
	if err != nil {
		return s, err
	}
	s.Cal.H6 = float64(msb)

	// Put the sensor in sleep mode and configure.
	err = bus.WriteByteToReg(addr, CTRL_MEAS_REG, 0x00)
	if err != nil {
		return s, err
	}

	// Set the config word.
	dataToWrite := (Standby << 0x5) & 0xE0
	dataToWrite = dataToWrite | ((Filter << 0x02) & 0x1C)
	err = bus.WriteByteToReg(addr, CONFIG_REG, dataToWrite)
	if err != nil {
		return s, err
	}

	dataToWrite = HumidOverSample & 0x07
	err = bus.WriteByteToReg(addr, CTRL_HUMIDITY_REG, dataToWrite)
	if err != nil {
		return s, err
	}

	dataToWrite = (TempOverSample << 0x5) & 0xE0
	dataToWrite = dataToWrite | ((PressOverSample << 0x02) & 0x1C)
	dataToWrite = dataToWrite | (RunMode & 0x03)
	err = bus.WriteByteToReg(addr, CTRL_MEAS_REG, dataToWrite)
	if err != nil {
		return s, err
	}

	_, err = bus.ReadByteFromReg(s.Addr, 0xD0)
	return s, err
}

func (s *BME280) fineT(d []byte) int32 {
	adcT := int32((uint32(d[TMP_MSB_IDX])<<12)|(uint32(d[TMP_LSB_IDX])<<4)) | ((int32(d[TMP_XLSB_IDX]) >> 4) & 0x0F)

	var1 := (((adcT >> 3) - (int32(s.Cal.T1) << 1)) * (int32(s.Cal.T2))) >> 11
	var2 := (((((adcT >> 4) - (int32(s.Cal.T1))) * ((adcT >> 4) - (int32(s.Cal.T1)))) >> 12) * (int32(s.Cal.T3))) >> 14

	return (var1 + var2)
}

// Reads all measurements from the sensor. Call Humidity, Pressure or Temperature to retrieve calibrated readings.
func (s *BME280) Measurements() ([]byte, error) {
	buf := make([]byte, 8)
	err := s.Bus.ReadFromReg(s.Addr, MEASUREMENT_START_REG, buf)

	return buf, err
}

// Humdity returns the relative humidity from the supplied measurements. Output value of "46.332" represents 46.332 %rH.
func (s *BME280) Humidity(d []byte) float64 {
	fineT := s.fineT(d)
	adcH := (uint16(d[HUMIDITY_MSB_IDX]) << 8) | uint16(d[HUMIDITY_LSB_IDX])

	varH := float64(fineT) - 76800.0
	varH = (float64(adcH) - (s.Cal.H4*64.0 + s.Cal.H5/16384.0*varH)) *
		(s.Cal.H2 / 65536.0 * (1.0 + s.Cal.H6/67108864.0*varH*(1.0+s.Cal.H3/67108864*varH)))
	varH = varH * (1.0 - s.Cal.H1*varH/524288.0)
	if varH > 100.0 {
		varH = 100.0
	} else if varH < 0.0 {
		varH = 0.0
	}

	return varH
}

// Returns the pressure in Pascals from the supplied measurements. A value of "96386.2" equals 963.862 hPa.
func (s *BME280) Pressure(d []byte) float64 {
	fineT := s.fineT(d)

	adcP := int32((uint32(d[PRESSURE_MSB_IDX])<<12)|
		(uint32(d[PRESSURE_LSB_IDX])<<4)) |
		((int32(d[PRESSURE_XLSB_IDX]) >> 4) & 0x0F)

	var1 := int64(fineT) - 128000
	var2 := var1 * var1 * s.Cal.P6
	var2 = var2 + (var1 * s.Cal.P5 << 17)
	var2 = var2 + (s.Cal.P4 << 35)
	var1 = (var1 * var1 * s.Cal.P3 >> 8) + (var1 * s.Cal.P2 << 12)
	var1 = ((int64(1) << 47) + var1) * s.Cal.P1 >> 33
	if var1 == 0 {
		return 0 // avoid exception caused by division by zero
	}
	p_acc := 1048576 - int64(adcP)
	p_acc = (((p_acc << 31) - var2) * 3125) / var1
	var1 = (s.Cal.P9 * (p_acc >> 13) * (p_acc >> 13)) >> 25
	var2 = (s.Cal.P8 * p_acc) >> 19
	p_acc = ((p_acc + var1 + var2) >> 8) + (s.Cal.P7 << 4)

	p_acc = p_acc >> 8 // /256
	return float64(p_acc)
}

// Temperature returns the temperature in Degrees Celcius from the supplied measurements. Output value of "30.33" equals 30.33°C.
func (s *BME280) Temperature(d []byte) float64 {
	fineT := s.fineT(d)
	return (float64(fineT) / 5120.0)
}
