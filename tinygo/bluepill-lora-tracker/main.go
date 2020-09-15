package main

import (
	"errors"
	"image/color"
	"machine"
	"strconv"
	"strings"
	"time"

	"./tools"

	"tinygo.org/x/drivers/gps"
	"tinygo.org/x/drivers/lora/sx127x"
	"tinygo.org/x/drivers/pcd8544"
)

var lcd *pcd8544.Device

var loraConfig = sx127x.Config{
	Frequency:       868000000,
	SpreadingFactor: 7,
	Bandwidth:       125000,
	CodingRate:      6,
	TxPower:         17,
}

var uartConsole machine.UART
var uartGps *machine.UART

var loraRadio sx127x.Device

var send_data = string("")
var send_delay = int(0)

type status struct {
	fix *gps.Fix
}

var st status

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
func gpsTask(pGps gps.Device, pParser gps.Parser) {
	var fix gps.Fix
	println("Start gpsTask")
	for {
		s, err := pGps.NextSentence()
		if err != nil {
			//println(err)
			continue
		}

		fix, err = pParser.Parse(s)
		if err != nil {
			//println(err)
			continue
		}
		print("*")

		if fix.Valid {
			st.fix = &fix
		}

		time.Sleep(500 * time.Millisecond)

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
func blinkTask(led machine.Pin, delay time.Duration) {
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

func initLCD() {
	dcPin := machine.PB12
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rstPin := machine.PB13
	rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	scePin := machine.PB14
	scePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.SPI0.Configure(machine.SPIConfig{})

	lcd = pcd8544.New(machine.SPI0, dcPin, rstPin, scePin)
	lcd.Configure(pcd8544.Config{})
}

func printSomething(msg string, x int16, y int16) {

	msg2 := []byte(msg)
	var c byte
	var col color.RGBA

	for k := 0; k < len(msg); k++ {
		c = msg2[k]

		for i := 0; i < 8; i++ {
			pix8 := tools.FontCP437[c][i]

			for j := 0; j < 8; j++ {
				col = color.RGBA{0, 0, 0, 255}

				if pix8&(1<<j) > 0 {
					col = color.RGBA{255, 255, 255, 255}
				}
				lcd.SetPixel(x+int16(k)*8+int16(i), y+int16(j), col)
			}
		}
	}
	lcd.Display()

}

func updateLCD() {
	if st.fix != nil {
		printSomething("St:"+strconv.Itoa(int(st.fix.Satellites)), 0, 0)
		printSomething("Al:"+strconv.Itoa(int(st.fix.Altitude)), 35, 0)
		printSomething("Sp:"+strconv.Itoa(int(st.fix.Speed)), 0, 9)
		printSomething("Hd:"+strconv.Itoa(int(st.fix.Heading)), 35, 9)
		lcd.Display()
	}

}

// main is .... main
func main() {

	// Led
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	go blinkTask(machine.LED, 250*time.Millisecond)

	// UARTs
	uartConsole = machine.UART0
	uartConsole.Configure(machine.UARTConfig{9600, 1, 0})
	println("GoTiny sx127x Demo")
	go consoleTask()

	// LCD
	initLCD()

	// GPS
	uartGps = machine.UART1
	uartGps.Configure(machine.UARTConfig{BaudRate: 9600})
	gps1 := gps.NewUART(uartGps)
	parser1 := gps.NewParser()
	go gpsTask(gps1, parser1)

	// LORA
	initRadio()

	var cycle uint32

	for {
		updateLCD()
		time.Sleep(5 * time.Second)
		cycle++
	}

}
