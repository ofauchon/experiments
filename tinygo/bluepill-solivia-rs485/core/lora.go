package core

import (
	"fmt"
	"machine"

	"tinygo.org/x/drivers/lora/sx127x"
)

func InitLora() {

	// Prepare gpio for RFM95 spi/dio
	SPI_CS_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	SPI_RST_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	DIO0_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Configure SPI
	machine.SPI0.Configure(machine.SPIConfig{
		SCK:       SPI_SCK_PIN,
		SDO:       SPI_SDO_PIN,
		SDI:       SPI_SDI_PIN,
		Frequency: 500000,
		Mode:      0})

	// Initialize RFM95 driver
	LoraRadio = sx127x.New(machine.SPI0, SPI_CS_PIN, SPI_RST_PIN)
	LoraRadio.Reset()

	// Setup DIO0 interrupt Handling
	err := DIO0_PIN.SetInterrupt(DIO0_PIN_CHANGE, func(machine.Pin) {
		if DIO0_PIN.Get() {
			LoraRadio.CheckIrq()
		}
	})
	if err != nil {
		println("could not configure pin interrupt:", err.Error())
	}

	// Configure Lora settings (modulation, SF... etc )
	err = LoraRadio.SetupLora(LoraConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
}
