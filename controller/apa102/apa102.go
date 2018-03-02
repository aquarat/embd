/*
 * Copyright (c) Clinton Freeman 2017
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

package apa102

import (
	"github.com/aquarat/embd"
)

const (
	MAX_BRIGHTNESS = 31
)

type APA102 struct {
	Bus  embd.SPIBus
	LED  []uint8
	nLED int
}

func New(bus embd.SPIBus, numLEDs int) *APA102 {
	res := &APA102{
		Bus:  bus,
		LED:  make([]uint8, numLEDs*4),
		nLED: numLEDs,
	}

	// Init the intensity field for each LED which is
	// 0b111 + 5 intensity bits.
	for i := 0; i < (numLEDs * 4); i = i + 4 {
		res.LED[i] = 224
	}

	return res
}

func (a *APA102) SetLED(index int, r uint8, g uint8, b uint8, i uint8) error {
	intensity := i
	if i > 31 {
		intensity = 31
	}

	ind := index * 4

	a.LED[ind+0] = 224 + intensity // Brightness is 0b111 + 5 intensity bits
	a.LED[ind+1] = b
	a.LED[ind+2] = g
	a.LED[ind+3] = r

	return a.Show()
}

func (a *APA102) Show() error {
	// Start frame is 32 zero bits.
	err := a.Bus.TransferAndReceiveData([]uint8{0, 0, 0, 0})

	// LED data.
	// 111+5bits(intensity) + 1byte(Red) + 1byte(Green) + 1byte(Blue)
	err = a.Bus.TransferAndReceiveData(a.LED)
	if err != nil {
		return err
	}

	err = a.Bus.TransferAndReceiveData([]uint8{1, 1, 1, 1})
	if err != nil {
		return err
	}

	return nil
}
