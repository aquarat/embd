/*
 * Copyright (c) Karan Misra 2014
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

// Package mcp4725 allows interfacing with the MCP4725 DAC.
package mcp4725

import (
	"sync"

	"github.com/aquarat/embd"
	"github.com/golang/glog"
)

const (
	dacReg     = 0x40
	programReg = 0x60
	powerDown  = 0x46

	genReset = 0x06
	powerUp  = 0x09
)

// MCP4725 represents a MCP4725 DAC.
type MCP4725 struct {
	// Bus to communicate over.
	Bus embd.I2CBus
	// Addr of the sensor.
	Addr byte

	initialized bool
	mu          sync.RWMutex
}

// New creates a new MCP4725 sensor.
func New(bus embd.I2CBus, addr byte) *MCP4725 {
	return &MCP4725{
		Bus:  bus,
		Addr: addr,
	}
}

func (d *MCP4725) setup() error {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	glog.V(1).Infof("mcp4725: general call reset")

	if err := d.Bus.WriteByteToReg(d.Addr, 0x00, powerUp); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(d.Addr, 0x00, genReset); err != nil {
		return err
	}
	d.initialized = true
	return nil
}

func (d *MCP4725) setVoltage(voltage int, reg byte) error {
	if err := d.setup(); err != nil {
		return err
	}
	if voltage > 4095 {
		voltage = 4095
	}
	if voltage < 0 {
		voltage = 0
	}

	glog.V(2).Infof("mcp4725: setting voltage to %04d", voltage)

	if err := d.Bus.WriteWordToReg(d.Addr, reg, uint16(voltage<<4)); err != nil {
		return err
	}
	return nil
}

// SetVoltage sets the output voltage.
func (d *MCP4725) SetVoltage(voltage int) error {
	return d.setVoltage(voltage, dacReg)
}

// SetPersistedVoltage sets the voltage and programs the EEPROM so
// that the voltage is restored on reboot.
func (d *MCP4725) SetPersistedVoltage(voltage int) error {
	return d.setVoltage(voltage, programReg)
}

// Close puts the DAC into power down mode.
func (d *MCP4725) Close() error {
	glog.V(1).Infof("mcp4725: powering down")

	if err := d.Bus.WriteWordToReg(d.Addr, powerDown, 0); err != nil {
		return err
	}
	return nil
}
