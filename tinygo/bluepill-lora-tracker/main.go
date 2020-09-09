package main

import (
	"errors"
	"machine"
	"strconv"
	"strings"
	"time"

	//nmea "./tools"

	"tinygo.org/x/drivers/gps"
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

// gpsTask handle communication with GPS Module
func gpsTask(parser gps.GPSParser) {
	var fix gps.Fix
	println("Start gpsTask")
	for {
		println("gpsTask tick")
		fix = parser.NextFix()
		/*
		   		println("aa ", fix.Valid, "alt", fix.Altitude)

		   		if fix.Valid {
		   			print(fix.Time.Format("15:04:05"))
		   			print(", lat=", fix.Latitude)
		   			print(", long=", fix.Longitude)
		   			print(", altitude:=")
		   			print(", satellites=", fix.Satellites)
		   			println()
		   		} else {
		   			println("No fix")
		   	println("Start gpsTask")

		   }
		*/

		time.Sleep(10 * time.Second)

	}
}

// serialTask handles interaction with the serial console
func consoleTask() string {
	inputConsole := make([]byte, 128) // serial port buffer

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
	uartConsole.Configure(machine.UARTConfig{9600, 1, 0})
	go consoleTask()

	// GPS
	uartGps = machine.UART1
	uartGps.Configure(machine.UARTConfig{BaudRate: 9600})
	ublox := gps.NewUART(&uartGps)
	parser := gps.Parser(ublox)
	go gpsTask(parser)

	println("GoTiny sx127x Demo")

	initRadio()

	var cycle uint32

	for {
		time.Sleep(5 * time.Second)
		cycle++
	}

}
