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

// Analog I/O support on the BBB.

package bbb

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/aquarat/embd"
)

type analogPin struct {
	id string
	n  int

	drv embd.GPIODriver

	val *os.File

	initialized bool
}

func newAnalogPin(pd *embd.PinDesc, drv embd.GPIODriver) embd.AnalogPin {
	return &analogPin{id: pd.ID, n: pd.AnalogLogical, drv: drv}
}

func (p *analogPin) N() int {
	return p.n
}

func (p *analogPin) init() error {
	if p.initialized {
		return nil
	}

	var err error
	if err = p.ensureEnabled(); err != nil {
		return err
	}
	if p.val, err = p.valueFile(); err != nil {
		return err
	}

	p.initialized = true

	return nil
}

func (p *analogPin) ensureEnabled() error {
	return ensureFeatureEnabled("cape-bone-iio")
}

func (p *analogPin) valueFilePath() (string, error) {
	pattern := fmt.Sprintf("/sys/devices/ocp.*/helper.*/AIN%v", p.n)
	return embd.FindFirstMatchingFile(pattern)
}

func (p *analogPin) openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY, os.ModeExclusive)
}

func (p *analogPin) valueFile() (*os.File, error) {
	path, err := p.valueFilePath()
	if err != nil {
		return nil, err
	}
	return p.openFile(path)
}

func (p *analogPin) Read() (int, error) {
	if err := p.init(); err != nil {
		return 0, err
	}

	p.val.Seek(0, 0)
	bytes, err := ioutil.ReadAll(p.val)
	if err != nil {
		return 0, err
	}
	str := string(bytes)
	str = strings.TrimSpace(str)
	return strconv.Atoi(str)
}

func (p *analogPin) Close() error {
	if err := p.drv.Unregister(p.id); err != nil {
		return err
	}

	if !p.initialized {
		return nil
	}

	if err := p.val.Close(); err != nil {
		return err
	}

	p.initialized = false

	return nil
}
