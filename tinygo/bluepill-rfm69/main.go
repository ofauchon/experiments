package main

import (
	"device/stm32"
	"errors"
	"machine"
	"runtime/interrupt"
	"strconv"
	"strings"
	"time"

	"tinygo.org/x/drivers/rfm69"
)

var uart machine.UART
var d *rfm69.Device

var send_data = string("")
var send_delay = int(0)

// processCmd parses commands and execute actions
func processCmd(cmd string) error {
	ss := strings.Split(cmd, " ")
	switch ss[0] {
	case "help":
		println("reset: reset rfm69 device")
		println("send xxxxxxx: send string over the air ")
		println("get: temp|mode|freq|regs")
		println("set: freq <433900000> set transceiver frequency (in Hz)")
		println("mode: <rx,tx,standby,sleep>")

	case "reset":
		d.Reset()
		println("Reset done !")

	case "send":
		if len(ss) == 2 {
			println("Scheduled data to send :", ss[1])
			send_data = ss[1]
			err := d.Send([]byte(send_data))
			if err != nil {
				println("Send error", err)
			}
		}
	case "get":
		switch ss[1] {
		case "freq":
			println("Freq:", d.GetFrequency())
		case "temp":
			temp, _ := d.ReadTemperature(0)
			println("Temperature:", temp)
		case "mode":
			mode := d.GetMode()
			println(" Mode:", mode)
		case "regs":
			for i := uint8(0); i < 0x60; i++ {
				val, _ := d.ReadReg(i)
				println(" Reg: ", strconv.FormatInt(int64(i), 16), " -> ", strconv.FormatInt(int64(val), 16))
			}
		default:
			return errors.New("Unknown command get")
		}

	case "set":
		switch ss[1] {
		case "freq":
			val, _ := strconv.ParseUint(ss[2], 10, 32)
			d.SetFrequency(uint32(val))
			println("Freq set to ", val)
		case "power":
			val, _ := strconv.ParseUint(ss[2], 10, 32)
			d.SetTxPower(uint8(val))
			println("TxPower set to ", val)
		}

	case "mode":
		switch ss[1] {
		case "standby":
			d.SetMode(rfm69.RFM69_MODE_STANDBY)
			d.WaitForMode()
			println("Mode changed !")
		case "sleep":
			d.SetMode(rfm69.RFM69_MODE_SLEEP)
			d.WaitForMode()
			println("Mode changed !")
		case "tx":
			d.SetMode(rfm69.RFM69_MODE_TX)
			d.WaitForMode()
			println("Mode changed !")
		case "rx":
			d.SetMode(rfm69.RFM69_MODE_RX)
			d.WaitForMode()
			println("Mode changed !")
		default:
			return errors.New("Unknown command mode")
		}
	default:
		return errors.New("Unknown command")
	}
	return nil
}

// serial() function is a gorouting for handling USART rx data
func serial() string {
	input := make([]byte, 64) // serial port buffer
	i := 0

	for {

		if uart.Buffered() > 0 {

			data, _ := uart.ReadByte() // read a character

			switch data {
			case 13: // pressed return key
				uart.Write([]byte("\r\n"))
				cmd := string(input[:i])
				err := processCmd(cmd)
				if err != nil {
					println(err)
				}
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

// blink is endless loop to make status led blink
func blink(led machine.Pin, delay time.Duration) {
	for {
		led.Low()
		time.Sleep(delay)
		led.High()
		time.Sleep(delay)
	}
}

// PB1_Int_Handler is interrupt handler from DIO1 signals
func PB1_Int_Handler(intr interrupt.Interrupt) {
	println("INTB1: ", machine.PB1.Get())
	stm32.EXTI.PR.SetBits(stm32.EXTI_PR_PR1_Msk) // Clear interrupt
	d.Receive()
}

func init_radio() {

	// SPI
	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 250000,
		Mode:      0},
	)

	// GPIO PB0 Interrupt DIO0 (RFM69 Data in Continuous mode)
	// Todo : Enable interrupt and timer to decode RF Pulses
	machine.PB0.Configure(machine.PinConfig{Mode: machine.PinInputModeFloating})

	// GPIO PB1 Interrupt DIO0 (RFM69 Interrupt on RX Packet)
	machine.PB1.Configure(machine.PinConfig{Mode: machine.PinInputModeFloating})
	stm32.RCC.APB2ENR.SetBits(stm32.RCC_APB2ENR_AFIOEN) // Enable AFIO
	stm32.AFIO.EXTICR1.Set(0x0001 << 4)                 // EXTICR1 configuration to enable PORTB1 trigger on EXTI1 line
	stm32.EXTI.RTSR.SetBits(stm32.EXTI_RTSR_TR1)        // Detect Rising Edge on EXTI1 Line
	//stm32.EXTI.FTSR.SetBits(stm32.EXTI_FTSR_TR1)      // Detect Falling Edge on EXTI1 Line
	stm32.EXTI.IMR.SetBits(stm32.EXTI_IMR_MR1) // Enable EXTI2 line
	intr := interrupt.New(stm32.IRQ_EXTI1, PB1_Int_Handler)
	intr.SetPriority(0xc0)
	intr.Enable()

	// RFM Configuration
	// (NSS => PA1, RESET => PA0, SPI => SPI0)
	machine.PA0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.PA1.Configure(machine.PinConfig{Mode: machine.PinOutput})

	cfg := [][]byte{
		/* 0x01 */ {rfm69.REG_OPMODE, rfm69.RF_OPMODE_SEQUENCER_ON | rfm69.RF_OPMODE_LISTEN_OFF | rfm69.RF_OPMODE_STANDBY},
		/* 0x02 */ {rfm69.REG_DATAMODUL, rfm69.RF_DATAMODUL_DATAMODE_PACKET | rfm69.RF_DATAMODUL_MODULATIONTYPE_OOK | rfm69.RF_DATAMODUL_MODULATIONSHAPING_00}, // no shaping
		/* 0x03 */ {rfm69.REG_BITRATEMSB, rfm69.RF_BITRATEMSB_4800}, // default: 4.8 KBPS
		/* 0x04 */ {rfm69.REG_BITRATELSB, rfm69.RF_BITRATELSB_4800},

		/* 0x05 */ {rfm69.REG_FDEVMSB, rfm69.RF_FDEVMSB_50000}, // 5khz
		/* 0x06 */ {rfm69.REG_FDEVLSB, rfm69.RF_FDEVLSB_50000},
		/* 0x07 */ {rfm69.REG_FRFMSB, rfm69.RF_FRFMSB_433}, // 433 Mhz
		/* 0x08 */ {rfm69.REG_FRFMID, rfm69.RF_FRFMID_433},
		/* 0x09 */ {rfm69.REG_FRFLSB, rfm69.RF_FRFLSB_433},

		/* 0x19 */ {rfm69.REG_RXBW, rfm69.RF_RXBW_DCCFREQ_001 | rfm69.RF_RXBW_MANT_16 | rfm69.RF_RXBW_EXP_0}, // (BitRate < 2 * RxBw)
		/* 0x25 */ {rfm69.REG_DIOMAPPING1, rfm69.RF_DIOMAPPING1_DIO0_01}, // DIO0 is the only IRQ we're using
		/* 0x26 */ {rfm69.REG_DIOMAPPING2, rfm69.RF_DIOMAPPING2_CLKOUT_OFF}, // DIO5 ClkOut disable for power saving
		/* 0x28 */ {rfm69.REG_IRQFLAGS2, rfm69.RF_IRQFLAGS2_FIFOOVERRUN}, // writing to this bit ensures that the FIFO & status flags are reset
		/* 0x29 */ {rfm69.REG_RSSITHRESH, 240}, //must be set to dBm = (-Sensitivity / 2) - default is 0xE4=228 so -114dBm
		/* 0x2E */ {rfm69.REG_SYNCCONFIG, rfm69.RF_SYNC_ON | rfm69.RF_SYNC_FIFOFILL_AUTO | rfm69.RF_SYNC_SIZE_2 | rfm69.RF_SYNC_TOL_0},
		/* 0x2d */ {rfm69.REG_PREAMBLELSB, 2}, // RF_SYNC_SIZE_2 => RF_SYNC_SIZE_1
		/* 0x2F */ {rfm69.REG_SYNCVALUE1, 0x2D}, // attempt to make this compatible with sync1 byte of RFM12B lib
		/* 0x2e */ {rfm69.REG_SYNCCONFIG, rfm69.RF_SYNC_ON | rfm69.RF_SYNC_FIFOFILL_AUTO | rfm69.RF_SYNC_SIZE_1 | rfm69.RF_SYNC_TOL_0},
		/* 0x2f */ {rfm69.REG_SYNCVALUE1, 0x9C}, //attempt to make this compatible with sync1 byte of RFM12B lib
		/* 0x37 */ {rfm69.REG_PACKETCONFIG1, rfm69.RF_PACKET1_FORMAT_VARIABLE | rfm69.RF_PACKET1_DCFREE_MANCHESTER | rfm69.RF_PACKET1_CRC_ON | rfm69.RF_PACKET1_CRCAUTOCLEAR_ON | rfm69.RF_PACKET1_ADRSFILTERING_OFF},
		/* 0x38 */ {rfm69.REG_PAYLOADLENGTH, 0xFE}, // in variable length mode: the max frame size, not used in TX
		/* 0x3C */ {rfm69.REG_FIFOTHRESH, rfm69.RF_FIFOTHRESH_TXSTART_FIFONOTEMPTY | rfm69.RF_FIFOTHRESH_VALUE}, // TX on FIFO not empty
		/* 0x3D */ {rfm69.REG_PACKETCONFIG2, rfm69.RF_PACKET2_RXRESTARTDELAY_NONE | rfm69.RF_PACKET2_AUTORXRESTART_ON | rfm69.RF_PACKET2_AES_OFF}, // RXRESTARTDELAY must match transmitter PA ramp-down time (bitrate dependent)
		/* 0x6F */ {rfm69.REG_TESTDAGC, rfm69.RF_DAGC_IMPROVED_LOWBETA0}, // run DAGC continuously in RX mode for Fading Margin Improvement, recommended default for AfcLowBetaOn=0
		{rfm69.REG_TESTAFC, 0},
	}

	/*
		cfg2 := [][]byte{
			{0x01, 0x04}, // RegOpMode: Standby Mode
			{0x02, 0x00}, // RegDataModul: Packet mode, FSK, no shaping
			{0x03, 0x0C}, // RegBitrateMsb: 10 kbps
			{0x04, 0x80}, // RegBitrateLsb
			{0x05, 0x01}, // RegFdevMsb: 20 kHz
			{0x06, 0x48}, // RegFdevLsb
			{0x07, 0xD9}, // RegFrfMsb: 868,15 MHz
			{0x08, 0x09}, // RegFrfMid
			{0x09, 0x9A}, // RegFrfLsb
			{0x18, 0x88}, // RegLNA: 200 Ohm impedance, gain set by AGC loop
			{0x19, 0x4C}, // RegRxBw: 25 kHz
			{0x2C, 0x00}, // RegPreambleMsb: 3 bytes preamble
			{0x2D, 0x03}, // RegPreambleLsb
			{0x2E, 0x88}, // RegSyncConfig: Enable sync word, 2 bytes sync word
			{0x2F, 0x41}, // RegSyncValue1: 0x4148
			{0x30, 0x48}, // RegSyncValue2
			{0x37, 0xD0}, // RegPacketConfig1: Variable length, CRC on, whitening
			{0x38, 0x40}, // RegPayloadLength: 64 bytes max payload
			{0x3C, 0x8F}, // RegFifoThresh: TxStart on FifoNotEmpty, 15 bytes FifoLevel
			{0x58, 0x1B}, // RegTestLna: Normal sensitivity mode
			{0x6F, 0x30}, // RegTestDagc: Improved margin, use if AfcLowBetaOn=0 (default)
		}
	*/
	d = rfm69.New(machine.SPI0, machine.PA0, machine.PA1, true)

	// Reset RFM
	d.Reset()
	d.Configure(cfg) // TODO FIXME
	d.SetTxPower(30)
	d.SetHighPower(true)

}

// main is .... main
func main() {

	// Led
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	go blink(machine.LED, 1000*time.Millisecond)

	// UART
	uart = machine.UART0
	uart.Configure(machine.UARTConfig{9600, 1, 0})
	go serial()
	println("GoTiny RFM69 Demo")

	init_radio()

	var cycle uint32

	for {
		/*
			println("Loop")
			d.SetMode(rfm69.RFM69_MODE_STANDBY)
			temp, err := d.ReadTemperature(0)
			if err != nil {
				println(err)
			}
			println("PB0 temp: ", temp, "  state:", machine.PB0.Get())
			d.SetFrequency(433900)

			if len(send_data) > 0 {
				println("Send ", len(send_data), "bytes")
				err := d.Send([]byte(send_data))
				if err != nil {
					println(err)
				}
			}
			//d.SetMode(rfm69.RFM69_MODE_RX)
		*/
		time.Sleep(5 * time.Second)
		cycle++
	}

}
