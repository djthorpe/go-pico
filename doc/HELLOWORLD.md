
# Hello, World

The "Hello, World" code pulses the on-board LED using the PWM controller,
so that it fades up and down.

```go
package main

import (
  . "github.com/djthorpe/go-pico"
)

// Define the pins used
var (
  LED = Pin(25)
)

// Global variables, for the PWM state
var (
  fade     uint16
  going_up bool
)

// Main function
func main() {
  LED.PWM().SetEnabled(true)
  LED.PWM().SetInterrupt(on_pwm_wrap)
  // Wait forever
  select {}
}

// Called when the PWM counter wraps
func on_pwm_wrap(pwm *PWM) {
  if going_up {
    fade = fade + 1
    if fade > 255 {
      fade = 255
      going_up = false
    }
  } else {
    if fade == 0 {
	  going_up = true
	} else {
      fade = fade - 1
    }
  }
  // Square the fade value to make the LED's brightness appear more linear
  // Note this range matches with the wrap value
  pwm.Set(LED, fade*fade)
}
```

It is assumed the LED is connected to GPIO pin 25, which can be connected to PWM module:

```go
LED.PWM().SetEnabled(true)
LED.PWM().SetInterrupt(on_pwm_wrap)
```

The interrupt handler is called when the PWM counter wraps, which cycles through values 0 to 255, in order to provide varying length pulses. In order to compile and run the code, plug in your Pico and ensure it's in BOOTSEL mode. A `Makefile` is used to compile the code, which is then placed in the `build` folder. So,

```shell
cd go-pico
make cmd/helloworld && picotool load -x build/helloworld.uf2
```

Next, you can read about how to use each module:

  * Pulse Width Modulation [PWM](PWM.md)
  * General Purpose IO [GPIO](GPIO.md)
  * Analog to Digital Converter [ADC](ADC.md)


