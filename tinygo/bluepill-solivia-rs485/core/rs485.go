package core

import (
	"fmt"
	"machine"
	"time"
)

// RS485Init() Configures uart and DE/RE gpio for RS485-TTL Adapter use
func RS485Init(port *machine.UART) {
	UartRS485 = port
	UartRS485.Configure(machine.UARTConfig{TX: UART2_TX_PIN, RX: UART2_RX_PIN, BaudRate: 19200})
	RS485_DERE_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

// RS485Send() sends byte to the RS485 UART
func RS485Send(request []uint8) {
	RS485_DERE_PIN.Set(true)
	print("RS485: Send ", len(request), " bytes :")
	for _, v := range request {
		fmt.Printf("%02X ", v)
	}
	println("")
	UartRS485.Write(request)
	RS485_DERE_PIN.Set(false)
}

// RS485Read reads bytes from serial port and until timeout
func RS485Read(timeoutSec int) []uint8 {

	input := make([]byte, 0) // serial port buffer
	cnt := 0

	for cnt < (timeoutSec * 100) {

		if UartRS485.Buffered() > 0 {
			stop := bool(false)
			for !stop {
				data, err := UartRS485.ReadByte()
				if err == nil {
					input = append(input, data)
					fmt.Printf("dbg/RS485: RX %02X\r\n ", data)
				} else {
					stop = true
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
		cnt++
	}

	return input
}
