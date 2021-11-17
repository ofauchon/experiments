package core

import (
	"errors"

	"github.com/snksoft/crc"
)

/*
"10 01" "IDC1"  4 => 0,4A (Solar current)
"10 02" "VDC1"  191 => 191v (solar voltage)
"10 03" "PDC1"  98 => (Solar power ? )
"10 04" "IDC2" => 0
"10 05" "VDC2" => 0
"10 06" "PDC2" => 0
"10 07" "IAC"  => 6 => 6A
"10 08" "VAC" => E7 = 213
"10 09" "PAC"
"10 0A" "FAC"
"11 01" "IDC1AVG"
"11 02" "VDC1AVG"
"11 03" "PDC1AVG"
"11 04" "IDC2AVG"
"11 05" "VDC2AVG"
"11 06" "PDC2AVG"
"11 07" "IACAVG"
"11 08" "VACAVG"
"11 09" "PACAVC"
"11 0A" "FACAVG"
"20 05" "ACTemp"
"21 08" "DCTemp"
"41 01" "DCHWFAIL"
"41 02" "ACHWFAIL"
"13 03" "EToday"
"13 04" "TToday"
"17 03" "ETotal"
"17 04" "TTotal"
"00 00" "InvType"
"00 01" "SerNoShort"
"00 07" "DataCode"
"00 08" "PartNo"
"00 09" "SerNo"
"00 40" "FWVer"

* STX , ENQ: 02 05
 * ID: 01
 * LEN: 02
 * CMD: 60 01
 * CRC: 85 FC
 * ETX: 03
 *
 *
 *
 *  DF DF FF DF DF DF DF DF FF DF DF DF
 * (pos=0) 02 06 (Command ACK)
 * (pos=2) 01 : (Device ID = 01)
 * (pos=3) 95 (Length of payload)
 * (pos=4) 60 01 (Command)
 * (pos=6)  45 4F 45 34 36 30 31 30 32 31 32  (Part# EOE46010212)
 * (pos=18) 31 31 33 32 31 32 30 31 31 30 31 35 30 30 30 38 36 30 (Serial# 132120110150008601)
 * (pos=37) 31 30 31 35 30 31 (101501) Date ?
 * (pos=43) 0A 01 (PWR MGMT Firmware Version ? 10.1)
 * 0A 02
 * 0A 00
 * 0A 00
 * 0A 00
 * (pos=53) 00 D7 (Voltage 215?)
 * (pos=55) 00 0C Current 12A
 * (pos=57) 03 E8 Power ? 1000W
 * (pos=59) 00 21 = 33 ?
 * (pos=61) 03 E8 = 1000 ?
 * (pos=63) 00 0B = 11 ?
 * (pos=65) 00 EA (VOltage = 234V?)
 * (pos= 67) 00 F0 (POwer = 240W)
 * (pos=69) 13 83 (Heztz 4995 mHz?)
 * 00 20 09 19 13 84 5B 65 13 83 00 06 5B 83 13 83 00 06 00 74 00 77 00 0B 00 E7 00 F0 00 F0 13 81 13 8C 00 02 A1 E4 00 00 73 93 00 80 01 54 0A AE 00 0C 03 E8 00 00 00 01
 DF DF DF DF DF DF DF

 LORA PAYLOAD:
 XX XX AC Current
 XX XX AC Inverter Voltage
 XX XX AC Inverter Current
 XX XX AC Inverter Power
 XX XX DC Voltage String 1
 XX XX DC Current String 1
 XX XX DC Power String 1
 XX XX Temperature1
 XX XX Temperature2

*/

type SoliviaDecoder struct {
}

type SoliviaInfos struct {
	Id           uint8
	PartNo       string
	SerialNo     string
	FwPowerVer   uint16
	FwStsVer     uint16
	FwDisplayVer uint16
	LastPacket   []uint8
	DateCode     string

	// AC
	ACFreq uint16
	ACVolt uint16
}

func NewSoliviaDecoder() *SoliviaDecoder {
	return &SoliviaDecoder{}
}

func (*SoliviaDecoder) GenCommand(id uint8, cmd [2]uint8) []uint8 {
	buf := make([]uint8, 0)
	buf = append(buf, 0x02, 0x05, id, 0x02, cmd[0], cmd[1])
	c := crc.CalculateCRC(crc.CRC16, []byte(buf[1:]))
	buf = append(buf, uint8(c&0xFF), uint8((c>>8)&0xFF), 0x03)
	//fmt.Println(hex.EncodeToString(buf))
	return buf
}

func (*SoliviaDecoder) SoliviaParseInfoMsg(buf []uint8) (*SoliviaInfos, error) {
	// Discard garbage/sync bytes
	for len(buf) > 2 {
		if buf[0] != 0x02 && buf[1] != 0x06 {
			buf = buf[1:]
		} else {
			break
		}
	}

	if len(buf) < 4 {
		return nil, errors.New("Packet too short, no headers")
	}

	l := buf[4] + 3

	// Following bytes
	if len(buf) < int(l) {
		return nil, errors.New("Packet length don't match header")
	}

	ret := &SoliviaInfos{}
	ret.LastPacket = buf
	ret.Id = buf[2]
	ret.PartNo = string(buf[6:17])
	ret.SerialNo = string(buf[18:36])
	ret.DateCode = string(buf[37:43])
	ret.ACFreq = uint16(buf[61])<<8 + uint16(buf[62])
	ret.ACVolt = uint16(buf[61])<<8 + uint16(buf[62])

	return ret, nil

}
