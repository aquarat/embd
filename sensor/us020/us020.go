/*
 * Copyright (c) Nikesh Vora 2013
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

// Package us020 allows interfacing with the US020 ultrasonic range finder.
package us020

import (
	"sync"
	"time"

	"github.com/aquarat/embd"
	"github.com/golang/glog"
)

const (
	pulseDelay  = 30000 * time.Nanosecond
	defaultTemp = 25
)

type Thermometer interface {
	Temperature() (float64, error)
}

type nullThermometer struct {
}

func (*nullThermometer) Temperature() (float64, error) {
	return defaultTemp, nil
}

var NullThermometer = &nullThermometer{}

// US020 represents a US020 ultrasonic range finder.
type US020 struct {
	EchoPin, TriggerPin embd.DigitalPin

	Thermometer Thermometer

	speedSound float64

	initialized bool
	mu          sync.RWMutex
}

// New creates a new US020 interface. The bus variable controls
// the I2C bus used to communicate with the device.
func New(echoPin, triggerPin embd.DigitalPin, thermometer Thermometer) *US020 {
	return &US020{EchoPin: echoPin, TriggerPin: triggerPin, Thermometer: thermometer}
}

func (d *US020) setup() error {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.TriggerPin.SetDirection(embd.Out)
	d.EchoPin.SetDirection(embd.In)

	if d.Thermometer == nil {
		d.Thermometer = NullThermometer
	}

	if temp, err := d.Thermometer.Temperature(); err == nil {
		d.speedSound = 331.3 + 0.606*temp

		glog.V(1).Infof("us020: read a temperature of %v, so speed of sound = %v", temp, d.speedSound)
	} else {
		d.speedSound = 340
	}

	d.initialized = true

	return nil
}

// Distance computes the distance of the bot from the closest obstruction.
func (d *US020) Distance() (float64, error) {
	if err := d.setup(); err != nil {
		return 0, err
	}

	glog.V(2).Infof("us020: trigerring pulse")

	// Generate a TRIGGER pulse
	d.TriggerPin.Write(embd.High)
	time.Sleep(pulseDelay)
	d.TriggerPin.Write(embd.Low)

	glog.V(2).Infof("us020: waiting for echo to go high")

	duration, err := d.EchoPin.TimePulse(embd.High)
	if err != nil {
		return 0, err
	}

	// Calculate the distance based on the time computed
	distance := float64(duration.Nanoseconds()) / 10000000 * (d.speedSound / 2)

	return distance, nil
}

// Close.
func (d *US020) Close() error {
	return d.EchoPin.SetDirection(embd.Out)
}
