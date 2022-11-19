
# Pulse Width Modulation (PWM)

```go
// Enable state
func (*PWM) SetEnabled(bool)
func (*PWM) Enabled() bool

// Counter
func (*PWM) SetCounter(value uint16) {
func (*PWM) Counter() uint16
func (*PWM) Inc()
func (*PWM) Dec()

// Wrapping value
func (*PWM) SetWrap(uint16)
func (*PWM) Wrap() uint16

// Period in nanoseconds
func (*PWM) SetPeriod(uint64) error
func (*PWM) Period() uint64

// Set and clear the interrupt when counter reaches wrap value
func (*PWM) SetInterrupt(func(*PWM))
```
