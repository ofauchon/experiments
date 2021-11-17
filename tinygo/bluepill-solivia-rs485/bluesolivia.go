// Receives data with LoRa.
package main

import (
	"machine"
	"time"

	"./core"
)

func main() {

	// Some Init
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED.Set(true)

	// Enable console
	core.ConsoleInit(machine.UART1)
	core.ConsoleStartTask()

	println("Delta Solivia Lora gateway")

	// Init RS485
	core.RS485Init(machine.UART2)

	// Initialize Lora
	core.InitLora()

	// Loop forever
	for {
		time.Sleep(time.Second)
		machine.LED.Set(!machine.LED.Get())
	}
}
