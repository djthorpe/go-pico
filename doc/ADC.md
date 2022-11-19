# Analog to Digital Converter (ADC)

```go
// Raw value from ADC
func (*ADC) Get() uint16

// Return voltage given the value of the reference voltage
func (*ADC) GetVoltage(float32) float32

// Return temperature ReadTemperature does a one-shot sample of the internal
// temperature sensor and returns a celsius reading
func (*ADC) GetTemperature() float32
```
