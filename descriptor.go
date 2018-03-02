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

// Host descriptor data structures.

package embd

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"log"
)

// Descriptor represents a host descriptor.
type Descriptor struct {
	GPIODriver func() GPIODriver
	I2CDriver  func() I2CDriver
	LEDDriver  func() LEDDriver
	SPIDriver  func() SPIDriver
}

// The Describer type is a Descriptor provider.
type Describer func(rev int) *Descriptor

// Describers is a global list of registered host Describers.
var describers = make(map[Host]Describer)

// Register makes a host describer available by the provided host key.
// If Register is called twice with the same host or if describer is nil,
// it panics.
func Register(host Host, describer Describer) {
	if describer == nil {
		panic("embd: describer is nil")
	}
	if _, dup := describers[host]; dup {
		return
		log.Println("embd: describer already registered")
	}
	describers[host] = describer

	glog.V(1).Infof("embd: host %v is registered", host)
}

var hostOverride Host
var hostRevOverride int
var hostOverriden bool

// SetHost overrides the host and revision no.
func SetHost(host Host, rev int) {
	hostOverride = host
	hostRevOverride = rev

	hostOverriden = true
}

// DescribeHost returns the detected host descriptor.
// Can be overriden by calling SetHost though.
func DescribeHost() (*Descriptor, error) {
	var host Host
	var rev int

	if hostOverriden {
		host, rev = hostOverride, hostRevOverride
	} else {
		var err error
		host, rev, err = DetectHost()
		if err != nil {
			return nil, err
		}
	}

	describer, ok := describers[host]
	if !ok {
		return nil, fmt.Errorf("host: invalid host %q", host)
	}

	return describer(rev), nil
}

// ErrFeatureNotSupported is returned when the host does not support a
// particular feature.
var ErrFeatureNotSupported = errors.New("embd: requested feature is not supported")

// ErrFeatureNotImplemented is returned when a particular feature is supported
// by the host but not implemented yet.
var ErrFeatureNotImplemented = errors.New("embd: requested feature is not implemented")
