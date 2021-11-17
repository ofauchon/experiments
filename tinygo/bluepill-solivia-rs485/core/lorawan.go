package core

import (
	"time"

	"github.com/ofauchon/zaza-tracker/libs"
)

//LoraWanInit() Prepares LoraWan stack for communication
func LoraWanInit() {
	LoraStack.Otaa.AppEUI = [8]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	LoraStack.Otaa.DevEUI = [8]uint8{0xA8, 0x40, 0x41, 0x00, 0x01, 0x81, 0xB3, 0x65}
	LoraStack.Otaa.AppKey = [16]uint8{0x2C, 0x44, 0xFC, 0xF8, 0x6C, 0x7B, 0x76, 0x7B, 0x8F, 0xD3, 0x12, 0x4F, 0xCE, 0x7A, 0x32, 0x16}

	r := GetRand32()
	LoraStack.Otaa.DevNonce[0] = uint8(r & 0xFF)
	LoraStack.Otaa.DevNonce[1] = uint8((r >> 8) & 0xFF)

	//	cpu := GetCpuId()
	println("Lorawan Configuration ")
	//	println("lorawan:  CPUID    : ", cpu[2], "/", cpu[1], "/", cpu[0])
	println("lorawan:  DevEUI   : ", libs.BytesToHexString(LoraStack.Otaa.DevEUI[:]))
	println("lorawan:  AppEUI   : ", libs.BytesToHexString(LoraStack.Otaa.AppEUI[:]))
	println("lorawan:  AppKey   : ", libs.BytesToHexString(LoraStack.Otaa.AppKey[:]))
	println("lorawan:  DevNounce: ", libs.BytesToHexString(LoraStack.Otaa.DevNonce[:]))

}

// LoraWanTask() routing deals with the LoraWan
func LoraWanTask() {

	msg := []byte("TinyGoLora")

	for {

		// Send join packet
		println("lorawan: Start JOIN sequence")
		payload, err := LoraStack.GenerateJoinRequest()
		if err != nil {
			println("lorawan: Error generating join request", err)
		}
		println("lorawan: Send JOIN request ", libs.BytesToHexString(payload))

		// Send join
		LoraRadio.TxLora(payload)

		println("lorawan: Wait for JOINACCEPT for 10s")
		// TODO Receive join Accept (Timeout 10s)
		LoraRadio.RxLora()
		// FIXME
		resp := []uint8{1, 2, 3, 4}
		if err != nil {
			println("lorawan: Error loraRx: ", err)
		}
		println("lorawan: Received a frame ")
		err = LoraStack.DecodeJoinAccept(resp)
		if (err) == nil {
			println("lorawan: Valid JOINACCEPT, now connected")

			println("lorawan:   DevAddr: ", libs.BytesToHexString(LoraStack.Session.DevAddr[:]), " (LSB)")
			println("lorawan:   NetID  : ", libs.BytesToHexString(LoraStack.Otaa.NetID[:]))
			println("lorawan:   NwkSKey: ", libs.BytesToHexString(LoraStack.Session.NwkSKey[:]))
			println("lorawan:   AppSKey: ", libs.BytesToHexString(LoraStack.Session.AppSKey[:]))
			// Sent sample message
			payload, err := LoraStack.GenMessage(0, msg)
			if err == nil {
				//println("TX_	UPMSG: --appkey ", libs.BytesToHexString(LoraStack.Session.AppSKey[:]), " --nwkkey ", libs.BytesToHexString(LoraStack.Session.NwkSKey[:]), " --hex", libs.BytesToHexString(payload))
				println("lorawan: Sending payload ", string(msg))
				LoraRadio.TxLora(payload)
			} else {
				println("lorawan: Error building uplink message")
			}

		} else {
			println("lorawan: Cant' decode message (join accept expected) ", err)
		}

		// Wait 60s
		println("SLEEP 60s")
		time.Sleep(time.Second * 60)
	} //for
}
