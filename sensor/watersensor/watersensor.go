/*
 * Copyright (c) Kashyap Kopparam 2013
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

// Package watersensor allows interfacing with the water sensor.
package watersensor

import (
	"sync"

	"github.com/aquarat/embd"
	"github.com/golang/glog"
)

// WaterSensor represents a water sensor.
type WaterSensor struct {
	Pin embd.DigitalPin

	initialized bool
	mu          sync.RWMutex
}

// New creates a new WaterSensor struct
func New(pin embd.DigitalPin) *WaterSensor {
	return &WaterSensor{Pin: pin}
}

func (d *WaterSensor) setup() error {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.Pin.SetDirection(embd.In); err != nil {
		return err
	}
	d.initialized = true

	return nil
}

// IsWet determines if there is water present on the sensor
func (d *WaterSensor) IsWet() (bool, error) {
	if err := d.setup(); err != nil {
		return false, err
	}

	glog.V(1).Infof("watersensor: reading")

	value, err := d.Pin.Read()
	if err != nil {
		return false, err
	}
	if value == embd.High {
		return true, nil
	} else {
		return false, nil
	}
}
