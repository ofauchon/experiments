package main

import (
	"errors"
	"fmt"
	"machine"
	"strconv"
	"strings"
	"time"

	nmea "./tools"

	"tinygo.org/x/drivers/lora/sx127x"
)

var loraConfig = sx127x.Config{
	Frequency:       868000000,
	SpreadingFactor: 7,
	Bandwidth:       125000,
	CodingRate:      6,
	TxPower:         17,
}

var uartConsole machine.UART
var uartGps machine.UART

var loraRadio sx127x.Device

var send_data = string("")
var send_delay = int(0)

// processCmd parses commands and execute actions
func processCmd(cmd string) error {
	ss := strings.Split(cmd, " ")
	switch ss[0] {
	case "help":
		println("reset: reset sx127x device")
		println("send xxxxxxx: send string over the air ")
		println("get: mode|freq|regs")
		println("set: freq <433900000> set transceiver frequency (in Hz)")
		println("mode: <rx,tx,standby,sleep>")

	case "reset":
		loraRadio.Reset()
		println("Reset done !")

	case "send":
		if len(ss) == 2 {
			println("Scheduled data to send :", ss[1])
			send_data = ss[1]
			loraRadio.SendPacket([]byte(send_data))
		}
	case "get":
		if len(ss) == 2 {
			switch ss[1] {
			case "freq":
				println("Freq:", loraRadio.GetFrequency())
				/*
					case "mode":
						mode := loraRadio.GetMode()
						println(" Mode:", mode)
				*/
				/*
					case "regs":
						for i := uint8(0); i < 0x60; i++ {
							val := loraRadio.ReadRegister(i)
							println(" Reg: ", i, " => ", val)
						}
				*/

			default:
				return errors.New("Unknown command get")
			}
		}

	case "set":
		if len(ss) == 3 {
			switch ss[1] {
			case "freq":
				val, _ := strconv.ParseUint(ss[2], 10, 32)
				loraRadio.SetFrequency(uint32(val))
				println("Freq set to ", val)
			case "power":
				val, _ := strconv.ParseUint(ss[2], 10, 32)
				loraRadio.SetTxPower(int8(val))
				println("TxPower set to ", val)
			}
		} else {
			println("invalid use of set command")
		}
	default:
		return errors.New("Unknown command")
	}
	return nil
}

// serial() function is a gorouting for handling USART rx data
func serial() string {
	inputConsole := make([]byte, 128) // serial port buffer
	inputGps := make([]byte, 128)     // serial port buffer

	for {

		// Process console messages
		for uartConsole.Buffered() > 0 {

			data, _ := uartConsole.ReadByte() // read a character

			switch data {
			case 13: // pressed return key
				uartConsole.Write([]byte("\r\n"))
				cmd := string(inputConsole)
				err := processCmd(cmd)
				if err != nil {
					println(err)
				}
				inputConsole = nil
			default: // pressed any other key
				uartConsole.WriteByte(data)
				inputConsole = append(inputConsole, data)
			}
		}

		// Process GPS messages
		for uartGps.Buffered() > 0 {

			data, _ := uartGps.ReadByte() // read a character
			if data == 10 {

				//fmt.Printf("\r\n***[GPS] %s", inputGps, " *** \r\n")

				//				s, _ := nmea.Parse(string(inputGps[:iGps]))
				//				if s.DataType() == nmea.TypeRMC {
				/*
					m := s.(nmea.RMC)
					fm	t.Printf("Latitude GPS: %s\n", nmea.FormatGPS(m.Latitude))
					fmt.Printf("Latitude DMS: %s\n", nmea.FormatDMS(m.Latitude))
					fmt.Printf("Longitude GPS: %s\n", nmea.FormatGPS(m.Longitude))
				*/
				//				}

				sentence := string(inputGps)
				//				if nmea.IsValidSentence(sentence) {
				endOfTypeIndex := strings.IndexByte(sentence, ',')
				sentenceType := sentence[1:endOfTypeIndex]
				//println("SentenceType: ", sentenceType)

				if sentenceType == "GPGGA" {
					g, err := nmea.ParseGPGGA(sentence)
					if err == nil {
						println("*Fix:", g.PositionFix, "Sat", g.UsedSatellites, " Long=", g.Longitude, " Lat=", g.Latitude, "Alt=", g.Altitude)
					} else {
						fmt.Println(err)
					}
				}

				if sentenceType == "GPRMC" {
					g, err := nmea.ParseGPRMC(sentence)
					if err == nil {
						//		fmt.Printf("* RMC Lon=%f Lat=%f Spd=%f Hdg=%f", g.Longitude, g.Latitude, g.Speed, g.Heading)
						println("*Status:", g.Status, " Spd=", g.Speed, "Head=%f", g.Heading)
					}
				}

				//				}

				inputGps = nil

			} else if data != 13 {
				inputGps = append(inputGps, data)
			}

		}

		time.Sleep(10 * time.Millisecond)
	}

}

//s, err := nmea.Parse(sentence)

// blink is endless loop to make status led blink
func blink(led machine.Pin, delay time.Duration) {
	for {
		led.Low()
		time.Sleep(delay)
		led.High()
		time.Sleep(delay)
	}
}

/*
 * Lora sx127x initialization
 *
 * sx127x RST => bluepill PA0
 * sx127x NSS => bluepill PA1
 * bluepill SPI: SPI0
 */
func initRadio() {

	// SPI
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 500000,
		Mode:      0},
	)

	// Extra SPI pins RST/CS to sx127x
	rstPin := machine.PA0
	rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	csPin := machine.PA1
	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	loraRadio = sx127x.New(machine.SPI0, csPin, rstPin)

	var err = loraRadio.Configure(loraConfig)
	if err != nil {
		println(err)
		return
	}

}

// main is .... main
func main() {

	// Led
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	go blink(machine.LED, 250*time.Millisecond)

	// UARTs
	uartConsole = machine.UART0
	uartGps = machine.UART1
	uartConsole.Configure(machine.UARTConfig{9600, 1, 0})
	uartGps.Configure(machine.UARTConfig{9600, 1, 0})
	go serial()

	println("GoTiny sx127x Demo")

	initRadio()

	var cycle uint32

	for {
		time.Sleep(5 * time.Second)
		cycle++
	}

}
