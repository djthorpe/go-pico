//go:build rpi

package gpio

/*
import (
	"os"

	"github.com/djthorpe/gopi"
)

type GPIO struct {
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE
/*

func (this *GPIO) New(gopi.Config) error {
	// Open the /dev/mem and provide offset & size for accessing memory
	if file, base, size, err := gpioOpenDevice(); err != nil {
		return err
	} else {
		defer file.Close()

		// Memory map GPIO registers to byte array
		if mem8, err := syscall.Mmap(int(file.Fd()), int64(base+GPIO_BASE), int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED); err != nil {
			return err
		} else {
			this.mem8 = mem8
		}

		// Convert mapped byte memory to unsafe []uint32 pointer, adjust length as needed
		header := *(*reflect.SliceHeader)(unsafe.Pointer(&this.mem8))
		header.Len /= (32 / 8)
		header.Cap /= (32 / 8)
		this.mem32 = *(*[]uint32)(unsafe.Pointer(&header))
	}

	// Check length of arrays
	if len(this.mem8) == 0 || len(this.mem32) == 0 {
		return gopi.ErrInternalAppError.WithPrefix("New")
	}

	// Set up pin watching
	this.watch = make(map[gopi.GPIOPin]gopi.GPIOState)

	// Return success
	return nil
}

func (this *GPIO) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if err := syscall.Munmap(this.mem8); err != nil {
		return os.NewSyscallError("munmap", err)
	}

	// Release resources
	this.pins = nil
	this.watch = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *GPIO) Run(ctx context.Context) error {
	timer := time.NewTicker(watchDelta)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			this.changeWatchState()
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *GPIO) String() string {
	str := "<gpio.broadcom"
	if p := this.NumberOfPhysicalPins(); p > 0 {
		str += " number_of_physical_pins=" + fmt.Sprint(p)
	}
	if l := this.pins; len(l) > 0 {
		str += " pins=" + fmt.Sprint(l)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PINS

// Return number of physical pins, or 0 if if cannot be returned
// or nothing is known about physical pins
func (this *GPIO) NumberOfPhysicalPins() uint {
	if this.product.Model == rpi.RPI_MODEL_A || this.product.Model == rpi.RPI_MODEL_B {
		return uint(26)
	} else {
		return uint(40)
	}
}

// Return array of available logical pins or nil if nothing is
// known about pins
func (this *GPIO) Pins() []gopi.GPIOPin {
	pins := make([]gopi.GPIOPin, GPIO_MAXPINS)
	for i := 0; i < GPIO_MAXPINS; i++ {
		pins[i] = gopi.GPIOPin(i)
	}
	return pins
}

// Return logical pin for physical pin number. Returns
// GPIO_PIN_NONE where there is no logical pin at that position
// or we don't know about the physical pins
func (this *GPIO) PhysicalPin(pin uint) gopi.GPIOPin {
	// Check for Raspberry Pi Version 1 and fudge things a little
	if this.product.Model == rpi.RPI_MODEL_A || this.product.Model == rpi.RPI_MODEL_B {
		// pin can be 1-28
		if pin < 1 || pin > 28 {
			return gopi.GPIO_PIN_NONE
		}
		if this.product.Revision == rpi.Revision(1) && pin == 3 {
			return gopi.GPIOPin(0)
		}
		if this.product.Revision == rpi.Revision(1) && pin == 5 {
			return gopi.GPIOPin(1)
		}
		if this.product.Revision == rpi.Revision(1) && pin == 13 {
			return gopi.GPIOPin(21)
		}
	}

	// now do things normally...
	if logical_pin, ok := pinmap[pin]; ok == false {
		return gopi.GPIO_PIN_NONE
	} else {
		return logical_pin
	}
}

// Return physical pin number for logical pin. Returns 0 where there
// is no physical pin for this logical pin, or we don't know anything
// about the layout
func (this *GPIO) PhysicalPinForPin(logical gopi.GPIOPin) uint {
	if physical, ok := this.pins[logical]; ok == false {
		return 0
	} else {
		return physical
	}
}

// Read pin state
func (this *GPIO) ReadPin(logical gopi.GPIOPin) gopi.GPIOState {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	var register uint32
	if uint8(logical) <= uint8(31) {
		// GPIO0 - GPIO31
		register = this.mem32[GPIO_GPLVL0>>2]
	} else {
		// GPIO32 - GPIO53
		register = this.mem32[GPIO_GPLVL1>>2]
	}
	if (register & (1 << (uint8(logical) & 31))) != 0 {
		return gopi.GPIO_HIGH
	}
	return gopi.GPIO_LOW
}

// Write pin state
func (this *GPIO) WritePin(logical gopi.GPIOPin, state gopi.GPIOState) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	value := uint32(1 << (uint8(logical) & 31))
	switch state {
	case gopi.GPIO_LOW:
		if uint8(logical) <= uint8(31) {
			this.mem32[GPIO_GPCLR0>>2] = value
		} else {
			this.mem32[GPIO_GPCLR1>>2] = value
		}
	case gopi.GPIO_HIGH:
		if uint8(logical) <= uint8(31) {
			this.mem32[GPIO_GPSET0>>2] = value
		} else {
			this.mem32[GPIO_GPSET1>>2] = value
		}
	}
}

// Get pin mode
func (this *GPIO) GetPinMode(logical gopi.GPIOPin) gopi.GPIOMode {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// return the register and the number of bits to shift to
	// access the current mode
	register, shift := gopiPinToRegister(logical)

	// Retrieve register, shift to the right, and return last three bits
	return gopi.GPIOMode((this.mem32[register>>2] >> shift) & 7)
}

// Set pin mode
func (this *GPIO) SetPinMode(logical gopi.GPIOPin, mode gopi.GPIOMode) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// get register and the number of bits to shift to
	// access the current mode
	register, shift := gopiPinToRegister(logical)

	// Set register
	this.mem32[register>>2] = (this.mem32[register>>2] &^ (7 << shift)) | (uint32(mode) << shift)
}

// Set pull mode to pull down or pull up - will
// return ErrNotImplemented if not supported
func (this *GPIO) SetPullMode(logical gopi.GPIOPin, pull gopi.GPIOPull) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check pin to make sure there is a physical pin mapping
	if this.PhysicalPinForPin(logical) == 0 {
		return gopi.ErrBadParameter.WithPrefix(fmt.Sprint(logical))
	}

	// Set the low two bits of register to 0 (off) 1 (down) or 2 (up)
	switch pull {
	case gopi.GPIO_PULL_UP, gopi.GPIO_PULL_DOWN:
		this.mem32[GPIO_GPPUD] |= uint32(pull)
	case gopi.GPIO_PULL_OFF:
		this.mem32[GPIO_GPPUD] &^= 3
	}

	// Wait for 150 cycles
	time.Sleep(time.Microsecond)

	// Determine clock register
	clockReg := GPIO_GPPUDCLK0
	if logical >= gopi.GPIOPin(32) {
		clockReg = GPIO_GPPUDCLK1
	}

	// Clock it in
	this.mem32[clockReg] = 1 << (logical % 32)

	// Wait for value to clock in
	time.Sleep(time.Microsecond)

	// Write 00 to the register to clear it
	this.mem32[GPIO_GPPUD] &^= 3

	// Wait for value to clock in
	time.Sleep(time.Microsecond)

	// Remove the clock
	this.mem32[clockReg] = 0

	// Return success
	return nil
}

// Start watching for rising and/or falling edge,
// or stop watching when GPIO_EDGE_NONE is passed.
// Will return ErrNotImplemented if not supported
func (this *GPIO) Watch(pin gopi.GPIOPin, edge gopi.GPIOEdge) error {
	// Check pin mode is INPUT
	if mode := this.GetPinMode(pin); mode != gopi.GPIO_INPUT {
		return gopi.ErrOutOfOrder.WithPrefix("Watch", pin)
	}

	// Get existing state of pin
	state := gopi.GPIO_LOW
	if edge != gopi.GPIO_EDGE_NONE {
		state = this.ReadPin(pin)
	}

	// Lock for writing
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Delete watch or set existing state
	if edge == gopi.GPIO_EDGE_NONE {
		delete(this.watch, pin)
		return nil
	} else {
		this.watch[pin] = state
	}

	// Return success
	return nil
}

func (this *GPIO) changeWatchState() {
	for pin, state := range this.watch {
		if newstate := this.ReadPin(pin); newstate == state {
			continue
		} else {
			this.RWMutex.Lock()
			defer this.RWMutex.Unlock()
			this.watch[pin] = newstate
		}
		if this.Publisher != nil {
			edge := gopi.GPIO_EDGE_NONE
			if state == gopi.GPIO_LOW {
				edge = gopi.GPIO_EDGE_RISING
			} else {
				edge = gopi.GPIO_EDGE_FALLING
			}
			this.Publisher.Emit(gpio.NewEvent(fmt.Sprint(pin), pin, edge), true)
		}
	}
}
*/
