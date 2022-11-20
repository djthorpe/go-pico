//go:build debug

package sdk

import "fmt"

//////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v *PWM_config) String() string {
	str := "<pwm_config"
	str += fmt.Sprintf(" csr=0x%08X", v.csr)
	str += fmt.Sprintf(" div=0x%08X", v.div)
	str += fmt.Sprintf(" top=0x%04X", v.top)
	return str + ">"
}

func (v PWM_clkdiv_mode) String() string {
	switch v {
	case PWM_DIV_FREE_RUNNING:
		return "PWM_DIV_FREE_RUNNING"
	case PWM_DIV_B_HIGH:
		return "PWM_DIV_B_HIGH"
	case PWM_DIV_B_RISING:
		return "PWM_DIV_B_RISING"
	case PWM_DIV_B_FALLING:
		return "PWM_DIV_B_FALLING"
	default:
		return fmt.Sprintf("PWM_clkdiv_mode(0x%02X)", uint(v))

	}
}

func (v PWM_chan) String() string {
	switch v {
	case PWM_CHAN_A:
		return "PWM_CHAN_A"
	case PWM_CHAN_B:
		return "PWM_CHAN_B"
	default:
		return fmt.Sprintf("PWM_chan(0x%02X)", uint(v))
	}
}
