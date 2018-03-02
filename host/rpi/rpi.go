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

/*
	Package rpi provides Raspberry Pi (including A+/B+) support.
	The following features are supported on Linux kernel 3.8+

	GPIO (digital (rw))
	I²C
	LED
*/
package rpi

import (
	"github.com/cfreeman/embd"
	"github.com/cfreeman/embd/host/generic"
	"log"
)

var spiDeviceMinor = 0

var rev1Pins = embd.PinMap{
	&embd.PinDesc{ID: "P1_3", Aliases: []string{"0", "GPIO_0", "SDA", "I2C0_SDA"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 0},
	&embd.PinDesc{ID: "P1_5", Aliases: []string{"1", "GPIO_1", "SCL", "I2C0_SCL"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 1},
	&embd.PinDesc{ID: "P1_7", Aliases: []string{"4", "GPIO_4", "GPCLK0"}, Caps: embd.CapDigital, DigitalLogical: 4},
	&embd.PinDesc{ID: "P1_8", Aliases: []string{"14", "GPIO_14", "TXD", "UART0_TXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 14},
	&embd.PinDesc{ID: "P1_10", Aliases: []string{"15", "GPIO_15", "RXD", "UART0_RXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 15},
	&embd.PinDesc{ID: "P1_11", Aliases: []string{"17", "GPIO_17"}, Caps: embd.CapDigital, DigitalLogical: 17},
	&embd.PinDesc{ID: "P1_12", Aliases: []string{"18", "GPIO_18", "PCM_CLK"}, Caps: embd.CapDigital, DigitalLogical: 18},
	&embd.PinDesc{ID: "P1_13", Aliases: []string{"21", "GPIO_21"}, Caps: embd.CapDigital, DigitalLogical: 21},
	&embd.PinDesc{ID: "P1_15", Aliases: []string{"22", "GPIO_22"}, Caps: embd.CapDigital, DigitalLogical: 22},
	&embd.PinDesc{ID: "P1_16", Aliases: []string{"23", "GPIO_23"}, Caps: embd.CapDigital, DigitalLogical: 23},
	&embd.PinDesc{ID: "P1_18", Aliases: []string{"24", "GPIO_24"}, Caps: embd.CapDigital, DigitalLogical: 24},
	&embd.PinDesc{ID: "P1_19", Aliases: []string{"10", "GPIO_10", "MOSI", "SPI0_MOSI"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 10},
	&embd.PinDesc{ID: "P1_21", Aliases: []string{"9", "GPIO_9", "MISO", "SPI0_MISO"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 9},
	&embd.PinDesc{ID: "P1_22", Aliases: []string{"25", "GPIO_25"}, Caps: embd.CapDigital, DigitalLogical: 25},
	&embd.PinDesc{ID: "P1_23", Aliases: []string{"11", "GPIO_11", "SCLK", "SPI0_SCLK"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 11},
	&embd.PinDesc{ID: "P1_24", Aliases: []string{"8", "GPIO_8", "CE0", "SPI0_CE0_N"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 8},
	&embd.PinDesc{ID: "P1_26", Aliases: []string{"7", "GPIO_7", "CE1", "SPI0_CE1_N"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 7},
}

var rev2Pins = embd.PinMap{
	&embd.PinDesc{ID: "P1_3", Aliases: []string{"2", "GPIO_2", "SDA", "I2C1_SDA"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 2},
	&embd.PinDesc{ID: "P1_5", Aliases: []string{"3", "GPIO_3", "SCL", "I2C1_SCL"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 3},
	&embd.PinDesc{ID: "P1_7", Aliases: []string{"4", "GPIO_4", "GPCLK0"}, Caps: embd.CapDigital, DigitalLogical: 4},
	&embd.PinDesc{ID: "P1_8", Aliases: []string{"14", "GPIO_14", "TXD", "UART0_TXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 14},
	&embd.PinDesc{ID: "P1_10", Aliases: []string{"15", "GPIO_15", "RXD", "UART0_RXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 15},
	&embd.PinDesc{ID: "P1_11", Aliases: []string{"17", "GPIO_17"}, Caps: embd.CapDigital, DigitalLogical: 17},
	&embd.PinDesc{ID: "P1_12", Aliases: []string{"18", "GPIO_18", "PCM_CLK"}, Caps: embd.CapDigital, DigitalLogical: 18},
	&embd.PinDesc{ID: "P1_13", Aliases: []string{"27", "GPIO_27"}, Caps: embd.CapDigital, DigitalLogical: 27},
	&embd.PinDesc{ID: "P1_15", Aliases: []string{"22", "GPIO_22"}, Caps: embd.CapDigital, DigitalLogical: 22},
	&embd.PinDesc{ID: "P1_16", Aliases: []string{"23", "GPIO_23"}, Caps: embd.CapDigital, DigitalLogical: 23},
	&embd.PinDesc{ID: "P1_18", Aliases: []string{"24", "GPIO_24"}, Caps: embd.CapDigital, DigitalLogical: 24},
	&embd.PinDesc{ID: "P1_19", Aliases: []string{"10", "GPIO_10", "MOSI", "SPI0_MOSI"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 10},
	&embd.PinDesc{ID: "P1_21", Aliases: []string{"9", "GPIO_9", "MISO", "SPI0_MISO"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 9},
	&embd.PinDesc{ID: "P1_22", Aliases: []string{"25", "GPIO_25"}, Caps: embd.CapDigital, DigitalLogical: 25},
	&embd.PinDesc{ID: "P1_23", Aliases: []string{"11", "GPIO_11", "SCLK", "SPI0_SCLK"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 11},
	&embd.PinDesc{ID: "P1_24", Aliases: []string{"8", "GPIO_8", "CE0", "SPI0_CE0_N"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 8},
	&embd.PinDesc{ID: "P1_26", Aliases: []string{"7", "GPIO_7", "CE1", "SPI0_CE1_N"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 7},
}

// This is the same as the Rev 2 for the first 26 pins.
var rev3Pins = append(append(embd.PinMap(nil), rev2Pins...), embd.PinMap{
	&embd.PinDesc{ID: "P1_29", Aliases: []string{"5", "GPIO_5"}, Caps: embd.CapDigital, DigitalLogical: 5},
	&embd.PinDesc{ID: "P1_31", Aliases: []string{"6", "GPIO_6"}, Caps: embd.CapDigital, DigitalLogical: 6},
	&embd.PinDesc{ID: "P1_32", Aliases: []string{"12", "GPIO_12"}, Caps: embd.CapDigital, DigitalLogical: 12},
	&embd.PinDesc{ID: "P1_33", Aliases: []string{"13", "GPIO_13"}, Caps: embd.CapDigital, DigitalLogical: 13},
	&embd.PinDesc{ID: "P1_35", Aliases: []string{"19", "GPIO_19"}, Caps: embd.CapDigital, DigitalLogical: 19},
	&embd.PinDesc{ID: "P1_36", Aliases: []string{"16", "GPIO_16"}, Caps: embd.CapDigital, DigitalLogical: 16},
	&embd.PinDesc{ID: "P1_37", Aliases: []string{"26", "GPIO_26"}, Caps: embd.CapDigital, DigitalLogical: 26},
	&embd.PinDesc{ID: "P1_38", Aliases: []string{"20", "GPIO_20"}, Caps: embd.CapDigital, DigitalLogical: 20},
	&embd.PinDesc{ID: "P1_40", Aliases: []string{"21", "GPIO_21"}, Caps: embd.CapDigital, DigitalLogical: 21},
}...)

var ledMap = embd.LEDMap{
	"led0": []string{"0", "led0", "LED0"},
}

func init() {
	embd.Register(embd.HostRPi, func(rev int) *embd.Descriptor {
		// Refer to http://elinux.org/RPi_HardwareHistory#Board_Revision_History
		// for details.

		pins := rev3Pins
		switch rev {
		case 0x2 : pins = rev1Pins
		case 0x3 : pins = rev1Pins
		case 0x4 : pins = rev2Pins
		case 0x5 : pins = rev2Pins
		case 0x6 : pins = rev2Pins
		case 0x7 : pins = rev2Pins
		case 0x8 : pins = rev2Pins
		case 0x9 : pins = rev2Pins
		case 0x000d : pins = rev2Pins
		case 0x000e : pins = rev2Pins
		case 0x000f : pins = rev2Pins
		case 0x10 : pins = rev1Pins
		case 0x11 : pins = rev1Pins
		case 0x12 : pins = rev1Pins
		case 0x13 : pins = rev1Pins
		case 0x14 : pins = rev1Pins
		case 0x15 : pins = rev1Pins
		case 0xa01040 : pins = rev1Pins
		case 0xa01041 : pins = rev1Pins
		case 0xa21041 : pins = rev1Pins
		case 0xa22042 : pins = rev1Pins
		case 0x900021 : pins = rev1Pins
		case 0x900032 : pins = rev1Pins
		case 0x900092 : pins = rev1Pins
		case 0x900093 : pins = rev1Pins
		case 0x920093 : pins = rev1Pins
		case 0x9000c1 : pins = rev1Pins
		case 0xa02082 : pins = rev3Pins
		case 0xa020a0 : pins = rev3Pins
		case 0xa22082 : pins = rev3Pins
		case 0xa32082 : pins = rev3Pins
		default : log.Println("Can't determine the revision of your board.")
		}

		switch {
		case pins[0].Aliases[0] == "0": log.Println("Selected revision 1 pin mapping.")
		case pins[0].Aliases[0] == "2": log.Println("Selected revision 2 pin mapping.")
		case pins[0].Aliases[0] == "5": log.Println("Selected revision 3 pin mapping.")
		}

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
