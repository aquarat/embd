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

// Package servo allows control of servos using a PWM controller.
package servo

import (
	"github.com/aquarat/embd/util"
	"github.com/golang/glog"
)

const (
	minus = 544
	maxus = 2400
)

const (
	// DefaultFreq represents the default (preferred) freq of a PWM doing servo duties.
	DefaultFreq = 50
)

// A PWM interface implements access to a pwm controller.
type PWM interface {
	SetMicroseconds(us int) error
}

type Servo struct {
	PWM PWM

	Minus, Maxus int
}

// New creates a new Servo interface.
func New(pwm PWM) *Servo {
	return &Servo{
		PWM:   pwm,
		Minus: minus,
		Maxus: maxus,
	}
}

// SetAngle sets the servo angle.
func (s *Servo) SetAngle(angle int) error {
	us := util.Map(int64(angle), 0, 180, int64(s.Minus), int64(s.Maxus))

	glog.V(1).Infof("servo: given angle %v calculated %v us", angle, us)

	return s.PWM.SetMicroseconds(int(us))
}
