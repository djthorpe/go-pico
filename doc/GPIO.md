# General Purpose Input/Output (GPIO)

```go
// Mode
func (Pin) Mode() Mode
func (Pin) SetMode(Mode)

// State
func (Pin) Set(bool)
func (Pin) Get() bool

// Set mode and return module
func (Pin) PWM() *PWM
func (Pin) ADC() *ADC

// Set Interrupt
func (Pin) SetInterrupt(callback Pin_callback_t)
```

The pin modes are as follows:

|----------------------------|----------------------------------|
| Mode                       | Description                      |
|----------------------------|----------------------------------|
|	`ModeOutput`             | Output mode                      |
|	`ModeInput`              | Input mode                       |
|	`ModeInputPulldown`      | Input mode with pulldown         |
|	`ModeInputPullup`        | Input mode with pullup           |
|	`ModeUART`               | UART mode                        |
|	`ModePWM`                | PWM mode                         |
|	`ModeI2C`                | I2C mode                         |
|	`ModeSPI`                | SPI mode                         |
|	`ModeOff`	             | Pin is off                       |
|----------------------------|----------------------------------|

