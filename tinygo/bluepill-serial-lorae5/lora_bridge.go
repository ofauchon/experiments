// This code accepts "TEST" commands on UART1 (TX:PA9 RX:PA10 9600baud)
// When received, bluepill sends AT command on UART2 (TX: RX: 9600baud)
// Use picocom --omap crlf -b 9600 /dev/ttyUSB0 (On Mac)

package main

import (
	"fmt"
	"machine"
	"time"
)

var (
	uartConsole, uartLora             *Xuart
	uartLoraBuffer, uartConsoleBuffer []byte
)

// This contains both UART and its command buffer
type Xuart struct {
	uart   *machine.UART
	buffer []byte
}

// processSerialCmd process command bytes from uarts
func processCmd(x *Xuart) {
	s := string(x.buffer)
	x.buffer = nil // Clear buffer

	if x == uartLora {
		print(s + "\r\n")
	} else if x == uartConsole {
		//println(s +"\r\n")
		uartLora.uart.Write([]byte(s + "\r\n"))
	}
}

// handleSerial concatenates stream of incoming bytes from uarts
func handleSerial(xu *Xuart, data byte) {
	if xu == uartConsole {
		print(string(data)) // local echo
	}

	switch data {
	case '\r': // discard \r
	case '\n': // only consider \n as end of line
		processCmd(xu)
	default: // pressed any other key
		xu.buffer = append(xu.buffer, data)
	}
}

// readSerial process incoming bytes from uarts
func readSerial() {
	for {
		for _, u := range []*Xuart{uartConsole, uartLora} {
			if u.uart.Buffered() > 0 {
				data, _ := u.uart.ReadByte()
				handleSerial(u, data)
			}
		}
		time.Sleep(time.Millisecond)
	}

}

func hwInit() {
	// Console serial
	uartConsole = &Xuart{uart: machine.UART1}
	uartConsole.uart.Configure(machine.UARTConfig{TX: machine.PA9, RX: machine.PA10, BaudRate: 9600})
	uartConsole.buffer = make([]byte, 300) // serial port buffer

	// Lora serial
	uartLora = &Xuart{uart: machine.UART2}
	uartLora.uart.Configure(machine.UARTConfig{TX: machine.PA2, RX: machine.PA3, BaudRate: 9600})
	uartLora.buffer = make([]byte, 300) // lora port buffer

}

func main() {

	hwInit()

	fmt.Println("Hello\r\n")

	// Start Task for reading serial
	go readSerial()

	for {
		println("I'm alive\r\n")
		time.Sleep(time.Second * 10)
	}
}
