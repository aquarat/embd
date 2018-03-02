/*
 * Copyright (c) Kunal Powar 2014
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

// Package mcp3008 allows interfacing with the mcp3008 8-channel, 10-bit ADC through SPI protocol.
package mcp3008

import (
	"github.com/aquarat/embd"
	"github.com/golang/glog"
)

// MCP3008 represents a mcp3008 8bit DAC.
type MCP3008 struct {
	Mode byte

	Bus embd.SPIBus
}

const (
	// SingleMode represents the single-ended mode for the mcp3008.
	SingleMode = 1

	// DifferenceMode represents the diffenrential mode for the mcp3008.
	DifferenceMode = 0
)

// New creates a representation of the mcp3008 convertor
func New(mode byte, bus embd.SPIBus) *MCP3008 {
	return &MCP3008{mode, bus}
}

const (
	startBit = 1
)

// AnalogValueAt returns the analog value at the given channel of the convertor.
func (m *MCP3008) AnalogValueAt(chanNum int) (int, error) {
	var data [3]uint8
	data[0] = startBit
	data[1] = uint8(m.Mode)<<7 | uint8(chanNum)<<4
	data[2] = 0

	glog.V(2).Infof("mcp3008: sendingdata buffer %v", data)
	if err := m.Bus.TransferAndReceiveData(data[:]); err != nil {
		return 0, err
	}

	return int(uint16(data[1]&0x03)<<8 | uint16(data[2])), nil
}
