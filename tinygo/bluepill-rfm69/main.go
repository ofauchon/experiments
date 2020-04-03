package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/rfm69"
)

var uart machine.UART

func serial() string {
	input := make([]byte, 64) // serial port buffer
	i := 0

	for {

		if uart.Buffered() > 0 {

			data, _ := uart.ReadByte() // read a character

			switch data {
			case 13: // pressed return key
				uart.Write([]byte("\r\n"))
				uart.Write([]byte("You typed: "))
				uart.Write(input[:i])
				uart.Write([]byte("\r\n"))
				i = 0
			default: // pressed any other key
				uart.WriteByte(data)
				input[i] = data
				i++
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

}

// blink the LED with given duration
func blink(led machine.Pin, delay time.Duration) {

	println("Hello world from Go!")

	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		led.Low()
		time.Sleep(delay)
		led.High()
		time.Sleep(delay)
	}
}

// start here at main function
func main() {

	go blink(machine.LED, 1000*time.Millisecond)

	// Init UART
	uart = machine.UART0
	uart.Configure(machine.UARTConfig{9600, 1, 0})
	uart.Write([]byte("Starting Golang RFM69 demo.\r\n"))
	go serial()

	rfm69.New(machine.SPI0, machine.PA1, machine.PA6, machine.PA7)

}
