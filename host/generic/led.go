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

// LED control support.

package generic

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aquarat/embd"
)

type led struct {
	id string

	brightness *os.File

	initialized bool
}

func NewLED(id string) embd.LED {
	return &led{id: id}
}

func (l *led) init() error {
	if l.initialized {
		return nil
	}

	var err error
	if l.brightness, err = l.brightnessFile(); err != nil {
		return err
	}

	l.initialized = true

	return nil
}

func (l *led) brightnessFilePath() string {
	return fmt.Sprintf("/sys/class/leds/%v/brightness", l.id)
}

func (l *led) openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, os.ModeExclusive)
}

func (l *led) brightnessFile() (*os.File, error) {
	return l.openFile(l.brightnessFilePath())
}

func (l *led) On() error {
	if err := l.init(); err != nil {
		return err
	}

	_, err := l.brightness.WriteString("1")
	return err
}

func (l *led) Off() error {
	if err := l.init(); err != nil {
		return err
	}

	_, err := l.brightness.WriteString("0")
	return err
}

func (l *led) isOn() (bool, error) {
	l.brightness.Seek(0, 0)
	bytes, err := ioutil.ReadAll(l.brightness)
	if err != nil {
		return false, err
	}
	str := string(bytes)
	str = strings.TrimSpace(str)
	if str == "1" {
		return true, nil
	}
	return false, nil
}

func (l *led) Toggle() error {
	if err := l.init(); err != nil {
		return err
	}

	state, err := l.isOn()
	if err != nil {
		return err
	}

	if state {
		return l.Off()
	}
	return l.On()
}

func (l *led) Close() error {
	if !l.initialized {
		return nil
	}

	if err := l.brightness.Close(); err != nil {
		return err
	}

	l.initialized = false

	return nil
}
