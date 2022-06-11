package main

import (
	"fmt"
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-pico/pkg/sdk"
)

func main() {
	// Initialise all pins
	//GPIO_init_mask(0xFFFFFFFF)

	// Set direction input for GP23 (BOOTSEL) annd GP25
	GPIO_init(GPIO_pin(23))
	GPIO_set_dir(GPIO_pin(23), GPIO_DIR_IN)

	// Set direction output for GP25 (LED) and switch on
	GPIO_init(GPIO_pin(25))
	GPIO_set_dir(GPIO_pin(25), GPIO_DIR_OUT)
	GPIO_put(GPIO_pin(25), true)

	// Report pins
	for pin := GPIO_pin(0); pin < NUM_BANK0_GPIOS; pin++ {
		fn := GPIO_get_function(pin)
		fmt.Println(pin, "=>", fn)
		if fn == GPIO_FUNC_SIO {
			if GPIO_get_dir(pin) == GPIO_DIR_IN {
				fmt.Println("  input => ", GPIO_get(pin))
				if GPIO_is_pulled_up(pin) {
					fmt.Println("    pulled up")
				} else if GPIO_is_pulled_down(pin) {
					fmt.Println("   pulled down")
				}
			} else {
				fmt.Println("  output => ", GPIO_get_out_level(pin))
			}
		}
	}

	// Return all outputs
	for {
		time.Sleep(time.Second)
		fmt.Printf("  input => %030b\n", GPIO_get_all())
		GPIO_put(GPIO_pin(25), !GPIO_get(GPIO_pin(25)))
	}
}
