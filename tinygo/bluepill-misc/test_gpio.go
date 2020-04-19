package main

import (
	"machine"
	"time"
)

// start here at main function
func main() {

	machine.PA0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.PA4.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		machine.PA0.Low()
		machine.PA4.High()
		time.Sleep(100 * time.Millisecond)
		machine.PA0.High()
		machine.PA4.Low()
		time.Sleep(100 * time.Millisecond)

	}

}
