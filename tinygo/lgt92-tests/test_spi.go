package main

import (
	"machine"
	"time"
)

// start here at main function
func main() {

	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 250000,
		Mode:      0},
	)

	i := byte(0)
	for {
		time.Sleep(20 * time.Millisecond)
		i = i + 1
		machine.SPI0.Transfer(i)
	}

}
