package main

import (
	"fmt"
	"machine"
	"time"
)

var (
	BATTERY_ADC = machine.ADC{machine.ADC3}
)

func main() {
	fmt.Println("=>adc")
	machine.InitADC()
	BATTERY_ADC.Configure(machine.ADCConfig{})

	for {
		battery_s := BATTERY_ADC.Get() >> 4
		fmt.Println("=>got")
		battery_v := float32(battery_s) * 3 * float32(3.3) / float32(1<<12)
		fmt.Println("battery_s=", battery_s, "battery_v=", battery_v)
		time.Sleep(time.Second)
	}
}
