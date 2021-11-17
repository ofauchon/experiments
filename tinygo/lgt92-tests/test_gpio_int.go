package main

import (
	"device/stm32"
	"machine"
	"runtime/interrupt"
	"time"
)

func intHandler(intr interrupt.Interrupt) {
	stm32.EXTI.PR.SetBits(1)
	println("Interrupt")
}

// start here at main function
func main() {

	println("Test")
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.PA0.Configure(machine.PinConfig{Mode: machine.PinInputModeFloating})

	stm32.RCC.APB2ENR.SetBits(stm32.RCC_APB2ENR_AFIOEN)        // Enable AFIO
	stm32.AFIO.EXTICR1.ClearBits(stm32.AFIO_EXTICR1_EXTI0_Msk) // Clear all EXTI0 bits to enable PORTA only
	stm32.EXTI.RTSR.SetBits(stm32.EXTI_RTSR_TR0)               // Detect Rising Edge of EXTI0 Line
	stm32.EXTI.FTSR.SetBits(stm32.EXTI_FTSR_TR0)               // Detect Falling Edge of EXTI0 Line
	stm32.EXTI.IMR.SetBits(stm32.EXTI_IMR_MR0)                 // Enable EXTI0 line

	intr := interrupt.New(stm32.IRQ_EXTI0, intHandler)
	intr.SetPriority(0xc0)
	intr.Enable()

	state := true
	for {
		println("Tick PA0:", machine.PA0.Get())
		time.Sleep(1 * time.Second)
		if state {
			machine.LED.High()
		} else {
			machine.LED.Low()

		}
		state = !state
	}

}
