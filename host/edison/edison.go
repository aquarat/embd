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

// Package edison Provides Intel Edison support.
package edison

import (
	"github.com/aquarat/embd"
	"github.com/aquarat/embd/host/generic"
)

var pins = embd.PinMap{
	&embd.PinDesc{ID: "P1_12", Aliases: []string{"12", "GPIO_12"}, Caps: embd.CapDigital, DigitalLogical: 12}, // PWM?
	&embd.PinDesc{ID: "P1_13", Aliases: []string{"13", "GPIO_13"}, Caps: embd.CapDigital, DigitalLogical: 13}, // PWM?
	&embd.PinDesc{ID: "P1_14", Aliases: []string{"14", "GPIO_14"}, Caps: embd.CapDigital, DigitalLogical: 14},
	&embd.PinDesc{ID: "P1_15", Aliases: []string{"15", "GPIO_15"}, Caps: embd.CapDigital, DigitalLogical: 15},
	&embd.PinDesc{ID: "P1_44", Aliases: []string{"44", "GPIO_44"}, Caps: embd.CapDigital, DigitalLogical: 44},
	&embd.PinDesc{ID: "P1_45", Aliases: []string{"45", "GPIO_45"}, Caps: embd.CapDigital, DigitalLogical: 45},
	&embd.PinDesc{ID: "P1_46", Aliases: []string{"46", "GPIO_46"}, Caps: embd.CapDigital, DigitalLogical: 46},
	&embd.PinDesc{ID: "P1_47", Aliases: []string{"47", "GPIO_47"}, Caps: embd.CapDigital, DigitalLogical: 47},
	&embd.PinDesc{ID: "P1_48", Aliases: []string{"48", "GPIO_48"}, Caps: embd.CapDigital, DigitalLogical: 48},
	&embd.PinDesc{ID: "P1_49", Aliases: []string{"49", "GPIO_49"}, Caps: embd.CapDigital, DigitalLogical: 49},
	&embd.PinDesc{ID: "P1_128", Aliases: []string{"128", "GPIO_128"}, Caps: embd.CapDigital, DigitalLogical: 128}, //CTS?
	&embd.PinDesc{ID: "P1_129", Aliases: []string{"129", "GPIO_129"}, Caps: embd.CapDigital, DigitalLogical: 129}, //RTS?
	&embd.PinDesc{ID: "P1_130", Aliases: []string{"130", "GPIO_130", "RXD", "UART0_RXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 130},
	&embd.PinDesc{ID: "P1_131", Aliases: []string{"131", "GPIO_131", "TXD", "UART0_TXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 131},
	&embd.PinDesc{ID: "P1_182", Aliases: []string{"182", "GPIO_182"}, Caps: embd.CapDigital, DigitalLogical: 182}, // PWM?
	&embd.PinDesc{ID: "P1_183", Aliases: []string{"183", "GPIO_183"}, Caps: embd.CapDigital, DigitalLogical: 183}, // PWM?
}

var ledMap = embd.LEDMap{}

var spiDeviceMinor = 0 // Check??

func init() {
	embd.Register(embd.HostEdison, func(rev int) *embd.Descriptor {
		return &embd.Descriptor{
			GPIODriver: func() embd.GPIODriver {
				return embd.NewGPIODriver(pins, generic.NewDigitalPin, nil, nil)
			},
			I2CDriver: func() embd.I2CDriver {
				return embd.NewI2CDriver(generic.NewI2CBus)
			},
			LEDDriver: func() embd.LEDDriver {
				return embd.NewLEDDriver(ledMap, generic.NewLED)
			},
			SPIDriver: func() embd.SPIDriver {
				return embd.NewSPIDriver(spiDeviceMinor, generic.NewSPIBus, nil)
			},
		}
	})
}
