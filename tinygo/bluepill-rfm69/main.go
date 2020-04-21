package main

import (
	"errors"
	"machine"
	"strconv"
	"strings"
	"time"

	"tinygo.org/x/drivers/rfm69"
)

var uart machine.UART
var d *rfm69.Device

var send_data = string("")
var send_delay = int(0)

func dumpIrqFlags(irqflags1, irqflags2 uint8)
{

	

}

func processCmd(cmd string) error {
	ss := strings.Split(cmd, " ")
	switch ss[0] {
	case "send":
		if len(ss) == 2 {
			send_data = ss[1]
		}
	case "get":
		switch ss[1] {
		case "temp":
			temp, _ := d.ReadTemperature(0)
			println("Temperature:", temp)
		case "mode":
			mode, _ := d.GetMode()
			println(" Mode:", mode)
		case "registers":
			for i := uint8(0); i < 0x60; i++ {
				val, _ := d.ReadReg(i)
				println(" Reg: ", strconv.FormatInt(int64(i), 16), " -> ", strconv.FormatInt(int64(val), 16))
				//				fmt.Printf(" Reg: %02X = %02X", i, val)
			}
		default:
			return errors.New("Unknown command")
		}

	case "rfm":
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
		case "setfreq":
			val, _ := strconv.ParseUint(ss[2], 10, 32)
			//i32 := uint32(val)
			d.SetFrequency(uint32(val))
		case "getfreq":
			println("Freq:", d.GetFrequency())
		default:
			return errors.New("Unknown parameter")
		}
	default:
		return errors.New("Unknown command")
	}
	return nil
}

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
	go serial()
	println("GoTiny Demo")

	// SPI
	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 250000,
		Mode:      0},
	)

	// RFM Configuration
	// (NSS => PA1, RESET => PA0, SPI => SPI0)
	machine.PA0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.PA1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d = rfm69.New(machine.SPI0, machine.PA0, machine.PA1, false)

	// Initialiqe RFM
	d.Reset()
	cfg := [][]byte{
		/* 0x01 */ {rfm69.REG_OPMODE, rfm69.RF_OPMODE_SEQUENCER_ON | rfm69.RF_OPMODE_LISTEN_OFF | rfm69.RF_OPMODE_STANDBY},
		/* 0x02 */ {rfm69.REG_DATAMODUL, rfm69.RF_DATAMODUL_DATAMODE_PACKET | rfm69.RF_DATAMODUL_MODULATIONTYPE_OOK | rfm69.RF_DATAMODUL_MODULATIONSHAPING_00}, // no shaping
		/* 0x03 */ {rfm69.REG_BITRATEMSB, rfm69.RF_BITRATEMSB_4800}, // default: 4.8 KBPS
		/* 0x04 */ {rfm69.REG_BITRATELSB, rfm69.RF_BITRATELSB_4800},
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
	d.Configure(cfg) // TODO FIXME
	d.SetFrequency(433900)
	d.SetTxPower(16)

	var cycle uint32

	for {

		if len(send_data) > 0 {
			println("Send ", len(send_data), "bytes")
			err := d.Send([]byte(send_data))
			if err != nil {
				println(err)
			}
		}

		time.Sleep(5 * time.Second)
		cycle++
	}

}
