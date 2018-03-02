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

// Package pcal9535a adds support for the low volage GPIO expander as found in the Raspberry
// Pi Relay board by Seeed Studio.
package pcal9535a

import (
	"github.com/aquarat/embd"
)

// More details at - http://wiki.seeedstudio.com/wiki/Raspberry_Pi_Relay_Board_v1.0
const (
	REG_MODE = 0x06
)

type PCAL9535A struct {
	Bus  embd.I2CBus
	Addr byte
	D    byte
}

// New creates and connects to a PCAL9535A GPIO expander.
func New(bus embd.I2CBus, addr byte) (*PCAL9535A, error) {
	return &PCAL9535A{
		Bus:  bus,
		Addr: addr,
		D:    0xff,
	}, bus.WriteByteToReg(addr, REG_MODE, 0xff)
}

// Sets the nominated GPIO pin to either high (on = true) or low (on = false)
func (c *PCAL9535A) SetPin(pin uint, on bool) error {
	if on {
		c.D &= ^(byte(0x1) << pin)
	} else {
		c.D |= (byte(0x1) << pin)
	}

	return c.Bus.WriteByteToReg(c.Addr, REG_MODE, c.D)
}

// Gets the state of supplied pin true = high or on, while false = low or off.
func (c *PCAL9535A) GetPin(pin uint) bool {
	return (((c.D >> pin) & 1) == 0)
}
