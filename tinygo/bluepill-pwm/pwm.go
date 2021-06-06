package main

import (
	"machine"
	"time"
)

func testPin(m string, t time.Duration, pin0, pin1 machine.Pin) {

	println("Pin Delay Tests (3 x ", m, ")")
	for i := 0; i < 3; i++ {
		pin0.High()
		time.Sleep(t)
		pin0.Low()
	}

	pin1.High()
	time.Sleep(time.Second * 1)
	pin1.Low()
	println("")

}

func setColor(pwm *machine.TIM, channels, color [3]uint8) {
	top := pwm.Top() - 1

	for i := 0; i < 3; i++ {
		if color[i] < 0 {
			color[i] = 0
		} else if color[i] > 100 {
			color[i] = 100
		}
	}
	pwm.Set(channels[0], uint32(color[0])*top/100)
	pwm.Set(channels[1], uint32(color[1])*top/100)
	pwm.Set(channels[2], uint32(color[2])*top/100)

	//	println(color[0], color[1], color[2])

}

func main() {

	println("PWM DEMO")

	pwm := &machine.TIM2
	pinA := machine.PA0
	pinB := machine.PA1
	pinC := machine.PA2

	var err error

	// Configure the PWM with the given period.
	err = pwm.Configure(machine.PWMConfig{
		Period: 16384e3, // 16.384ms
	})
	if err != nil {
		println("failed to configure PWM")
		return
	}

	// The top value is the highest value that can be passed to PWMChannel.Set.
	// It is usually an even number.
	println("top:", pwm.Top())

	ch := [3]uint8{0, 0, 0}
	// Configure the two channels we'll use as outputs.
	ch[0], err = pwm.Channel(pinA)
	if err != nil {
		println("failed to configure channel A")
		return
	}
	ch[1], err = pwm.Channel(pinB)
	if err != nil {
		println("failed to configure channel B")
		return
	}

	ch[2], err = pwm.Channel(pinC)
	if err != nil {
		println("failed to configure channel C")
		return
	}

	valR := int8(0)
	valG := int8(0)
	valB := int8(0)

	stepR := int8(2)
	stepG := int8(2)
	stepB := int8(2)

	for {

		valR += stepR
		valG += stepG
		valB += stepB

		setColor(pwm, ch, [3]uint8{uint8(valR), uint8(valG), uint8(valB)})
		time.Sleep(time.Millisecond * 20)

		if valR > 50 || (pwm.Count()%40) == 0 {
			stepR = -stepR
			valR = 0
		}
		if valG > 50 || (pwm.Count()%40) == 0 {
			stepG = -stepG
			valG = 0
		}
		if valB > 50 || (pwm.Count()%40) == 0 {
			stepB = -stepB
			valB = 0
		}

	}

}
