// +build ignore

package main

import (
	"flag"
	"fmt"

	"github.com/aquarat/embd"

	_ "github.com/aquarat/embd/host/all"
)

func main() {
	flag.Parse()

	embd.InitGPIO()
	defer embd.CloseGPIO()

	val, _ := embd.AnalogRead(0)
	fmt.Printf("Reading: %v\n", val)
}
