package main

import (
	"machine"
	"time"
)

// start here at main function
func main() {

	machine.LED_GREEN.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		machine.LED_GREEN.Low()
		time.Sleep(1 * time.Second)

		machine.LED_GREEN.High()
		time.Sleep(1 * time.Second)

	}

}
