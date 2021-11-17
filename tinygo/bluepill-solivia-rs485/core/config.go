package core

import (
	"machine"

	"github.com/ofauchon/go-lorawan-stack"
	"tinygo.org/x/drivers/lora/sx127x"
)

var LoraConfig = sx127x.Config{
	Frequency:       868100000,
	SpreadingFactor: 12,
	Bandwidth:       125000,
	CodingRate:      6,
	TxPower:         17,
}

var (
	UartConsole *machine.UART
	UartRS485   *machine.UART
	LoraStack   lorawan.LoraWanStack
	LoraRadio   sx127x.Device
)

const (

	// Serial console
	UART1_TX_PIN = machine.PA9
	UART1_RX_PIN = machine.PA10

	// Serial to RS485
	UART2_TX_PIN = machine.PA2
	UART2_RX_PIN = machine.PA3

	// RFM95 SPI Connection to Bluepill
	SPI_SCK_PIN = machine.PA5
	SPI_SDO_PIN = machine.PA7
	SPI_SDI_PIN = machine.PA6
	SPI_CS_PIN  = machine.PB8
	SPI_RST_PIN = machine.PB9

	// DIO RFM95 Pin connection to BluePill
	DIO0_PIN        = machine.PA0
	DIO0_PIN_MODE   = machine.PinInputPulldown
	DIO0_PIN_CHANGE = machine.PinRising

	RS485_DERE_PIN = machine.PA1
)
