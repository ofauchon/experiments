package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/rfm69"
)

var uart machine.UART

/*
 * Gorouting for Handling serial communication
 */
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

/*
 * Gorouting for LED blinking
 */
func blink(led machine.Pin, delay time.Duration) {
	for {
		led.Low()
		time.Sleep(delay)
		led.High()
		time.Sleep(delay)
	}
}

// start here at main function
func main() {

	// Led
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	go blink(machine.LED, 1000*time.Millisecond)

	// UART
	uart = machine.UART0
	uart.Configure(machine.UARTConfig{9600, 1, 0})
	uart.Write([]byte("Starting Golang RFM69 demo.\r\n"))
	go serial()

	// SPI
	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 250000,
		Mode:      0},
	)

	// RFM Configuration
	// (NSS => PA1, RESET => PA0, SPI => SPI0)
	machine.PA0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.PA1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d := rfm69.New(machine.SPI0, machine.PA0, machine.PA1, false)

	for {

		println("Temperature:", d.ReadTemperature(0))
		time.Sleep(10 * time.Second)

	}

}
