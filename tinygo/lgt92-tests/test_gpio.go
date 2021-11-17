package main

import (
	"machine"
	"time"
)

// start here at main function
func main() {

	pinLed1 := machine.PA12
	pinLed2 := machine.PA8
	pinLed3 := machine.PA11

	pinLed1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pinLed2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pinLed3.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		pinLed1.High()
		pinLed2.High()
		pinLed3.High()
		time.Sleep(1000 * time.Millisecond)
		pinLed1.Low()
		pinLed2.Low()
		pinLed3.Low()
		time.Sleep(1000 * time.Millisecond)

	}

}
