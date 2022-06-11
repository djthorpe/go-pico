package main

import (
	"device/rp"
	"machine"
	"runtime/interrupt"
)

type inter struct {
	Pin machine.Pin
}

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	button := machine.Pin(23) // BOOTSEL button
	button.Configure(machine.PinConfig{Mode: machine.PinInput})

	i := NewInterrupt(button, machine.PinFalling|machine.PinRising, func() {
		if button.Get() {
			machine.LED.Low()
		} else {
			machine.LED.High()
		}
	})

	for {
	}
}

func NewInter(p machine.Pin, c machine.PinChange, f func()) {
	i := new(inter)
	i.Pin = p
	i.Change = c
	i.intr = interrupt.New(rp.IRQ_IO_IRQ_BANK0, handleInterrupt)
}

func SetInterrupt(p machine.Pin, c machine.PinChange, f func()) {
	intr := interrupt.New(rp.IRQ_IO_IRQ_BANK0, handleInterrupt)
}

func handleInhandleInterrupt(interrupt.Interrupt) {
	64	t.acknowledgeChange()
	65	t.action()
	66}