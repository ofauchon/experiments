package core

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGenCommand(t *testing.T) {
	var dec SoliviaDecoder
	b := dec.GenCommand(0x01, [2]uint8{0x10, 0x01})
	t.Log(b)
}

func TestDecode(t *testing.T) {
	h := "DFDFFFDFDFDFDFDFFFDFDFDF020601956001454F4534363031303231323131333231323031313031353030303836303130313530310A010A020A000A000A0000D7000C03E8002103E8000B00EA00F013830020091913845B65138300065B831383000600740077000B00E700F000F01381138C0002A1E400007393008001540AAE000C03E800000001DFDFDFDFDFDFDF"
	d, _ := hex.DecodeString(h)
	t.Log("Sample length", len(d))

	//	dec := NewSoliviaDecoder()

	var dec SoliviaDecoder

	s, _ := dec.SoliviaParseInfoMsg(d)
	fmt.Println(s.Id)
}
