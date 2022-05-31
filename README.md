# go-pico

Experiments with Raspberry Pi Pico (RP2040) Microcontroller. 
The aim of this repository is to test the Pico Microcontroller with
[tinygo](https://tinygo.org/).

The code is abstracted so it can run either on a Raspberry Pi or on
the Pico. You'll use the Raspberry Pi to compile and connect from the Pico
through UART device to display output.

Presently you'll need a Raspberry Pi (and a Pico) and the following sort of
environment:

  * Use `armhf` rather than `arm64` operating system on the Raspberry Pi as that
    is easiest to install the **tinygo** compiler. No luck compiling for `arm64`
    yet;
  * You'll need to install **golang** on your Rasperry Pi;
  * Use **minicom** on your Raspberry Pi for displaying output from the pico;
  * Download the Pico SDK and build **picotool** for flashing the code;
  * You'll need GNU Make to compile the example code in `cmd`.

There are probably lots of other dependencies I haven't listed here....

## Installation

Assuming a working operating system on your Raspberry Pi, create an `/opt`
folder for your installation if it's not already created, and install dependencies:

```bash
sudo install -o ${USER} -d /opt && cd /opt
sudo apt -y install \
  git wget \
  cmake gcc-arm-none-eabi libnewlib-arm-none-eabi libstdc++-arm-none-eabi-newlib \
  build-essential pkg-config libusb-1.0-0-dev
```

Proceed to install a recent version of **golang**:

```bash
GOBUILD="go1.18.1.linux-armv6l"

cd /opt && wget https://redirector.gvt1.com/edgedl/go/${GOBUILD}.tar.gz  
install -d "/opt/${GOBUILD}" && tar -C "/opt/${GOBUILD}" -zxvf "${BUILD}.tar.gz" && rm -f "/opt/${GOBUILD}.tar.gz"  
rm -f /opt/go && cd /opt && ln -s "${GOBUILD}/go" go
```

Proceed to install **tinygo**:

```bash
TINYGOBUILD="0.23.0"

cd /opt && wget https://github.com/tinygo-org/tinygo/releases/download/v${TINYGOBUILD}/tinygo${TINYGOBUILD}.linux-arm.tar.gz
install -d "/opt/tinygo-${TINYGOBUILD}" && tar -C "/opt/tinygo-${TINYGOBUILD}" -zxvf "tinygo${TINYGOBUILD}.linux-arm.tar.gz" && rm -f "/opt/tinygo${TINYGOBUILD}.linux-arm.tar.gz"
rm -f /opt/tinygo && cd /opt && ln -s "tinygo-${TINYGOBUILD}/tinygo" tinygo
```

Install [**picosdk**](https://github.com/raspberrypi/pico-sdk) and [**picotool**](https://github.com/raspberrypi/picotool) so that you can flash your Pico:

```bash
PICOSDK="1.3.1"
PICOTOOL="1.1.0"

cd /opt && install -d pico && cd /opt/pico \
  && git clone -q --branch ${PICOSDK} --single-branch git@github.com:raspberrypi/pico-sdk.git pico-sdk-${PICOSDK} \
  && git clone -q --branch ${PICOTOOL} --single-branch git@github.com:raspberrypi/picotool.git picotool-${PICOTOOL}
rm -f /opt/pico/sdk && cd /opt/pico && ln -s "pico-sdk-${PICOSDK}" sdk
rm -f /opt/pico/picotool && cd /opt/pico && ln -s "picotool-${PICOTOOL}" picotool
cd /opt/pico/picotool && install -d build && cd build && PICO_SDK_PATH=/opt/pico/sdk cmake .. && make && install picotool /opt/pico/bin/picotool 
```

Add **golang** and **tinygo** to your path:

```bash
cat >> "${HOME}/.profile" <<EOF
if [ -x "/opt/go" ] ; then
  export PATH="\${PATH}:/opt/go/bin"
fi
if [ -x "/opt/tinygo" ] ; then
  export PATH="\${PATH}:/opt/tinygo/bin"
fi
if [ -x "/opt/pico" ] ; then
  export PATH="\${PATH}:/opt/pico/bin"
  export PICO_SDK_PATH="/opt/pico/sdk"
fi
EOF
```

## Testing Installation

To test you installation, log out and back in again, and then check for **golang** and **tinygo** versions:

```bash
bash% go version
go version go1.18.1 linux/arm
bash% tinygo version
tinygo version 0.23.0 linux/arm (using go version go1.18.1 and LLVM version 14.0.0)
bash% picotool version
picotool v1.1.0 (Linux 5.15.32-v7+, GNU-10.2.1, Release
```

## Connecting the Pico

The pinouts for the Pico are listed [here](https://datasheets.raspberrypi.com/pico/Pico-R3-A4-Pinout.pdf). You will
want to connect a reset button to the Pico and connect the default UART to your Raspberry Pi. For the Pico, the pinouts are as follows (when orientating the device face-up so the USB port is at the top):

| Wire     | Pin | Pin | Wire     |
|----------|-----|-----|----------|
| UART0 TX |  1  | 40  |          |
| UART0 RX |  2  | 39  |          |
| GND      |  3  | 38  |          |
|          |  4  | 37  |          |
|          |  5  | 36  |          |
|          |  6  | 35  |          |
|          |  7  | 34  |          |
|          |  8  | 33  |          |
|          |  9  | 32  |          |
|          | 10  | 31  |          |
|          | 11  | 30  | RESET    |
|          | 12  | 29  |          |
|          | 13  | 28  | GND      |
|          | 14  | 27  |          |
|          | 15  | 26  |          |
|          | 16  | 25  |          |
|          | 17  | 24  |          |
|          | 18  | 23  |          |
|          | 19  | 22  |          |
|          | 20  | 21  |          |

Connect Pins 28 and 30 to a push button. For the Raspberry Pi, orientating the device face-up so the power port is at the top left, the pinouts are as follows:

| Wire     | Pin | Pin | Wire     |
|----------|-----|-----|----------|
|          |  1  |  2  |          |
|          |  3  |  4  |          |
|          |  5  |  6  | GND      |
|          |  7  |  8  | UART  TX |
|          |  9  | 10  | UART  RX |
|          | 11  | 12  |          |
|          | 13  | 14  | GND      |
|          | 15  | 16  |          |
|          | 17  | 18  |          |
|          | 19  | 20  |          |
|          | 21  | 22  |          |
|          | 23  | 24  |          |
|          | 25  | 26  |          |
|          | 27  | 28  |          |
|          | 29  | 30  |          |
|          | 31  | 32  |          |
|          | 33  | 34  |          |
|          | 35  | 36  |          |
|          | 37  | 38  |          |
|          | 39  | 40  |          |

Connect **TX to RX** and **RX to TX** on the devices, and pair the **GND** connections. Set up **minicom** on your Raspberry Pi:

```bash
sudo usermod -a -G tty ${USER}
sudo apt -y install minicom
sudo cat > "/etc/minicom/minirc.dfl" <<EOF
pu port             /dev/serial0
pu addlinefeed      Yes
pu linewrap         Yes
pu addcarreturn     Yes
EOF
```

Use `raspi-config` to enable the serial port **but do not enable login through the serial port**:

  1. Interface Options
  2. Serial Port
  3. Would you like a login shell to be accessible over serial? **No**
  4. Would you like the serial port hardware to be enabled? **Yes**

You may need to reboot your Raspberry Pi.

## Blink

The purpose of the "Blink" application is to control the GPIO pins to switch an LED (light emitting diode) on and off.
You can blink the Pico's in-built LED or wire up a Raspberry Pi with an LED in order to demonstrate this.

What you'll need:

  * A light emitting diode of any colour or variation;
  * A resistor, probably between 5 and 100 ohms in value.

### Blink for Pico

Download the code and compile the **blink** application for the Pico using the following commands:

```bash
install -d ${HOME}/projects && cd ${HOME}/projects
git clone https://github.com/djthorpe/go-pico.git
cd go-pico && make cmd/pico/blink
```

This will place the application `blink.uf2` in the `build` folder. Then use **picotool** to flash your pico. You may need to `sudo` as the permissions:

```bash
sudo /opt/pico/bin/picotool load -x build/blink.uf2
```

You'll get an error if the Pico isn't in **BOOLSEL** mode. In order to do
that, hold the **BOOLSEL** button down on the Pico whilst pressing your **RESET** button. Release the latter button before the former, and try again.
The `-x` flag forces a device reset, so the application should run
immediately after load.

You should also try to run **minicom** (perhaps in a separate window) to see
debugging output of your application:

```bash
minicom --device /dev/serial0 --baudrate 115200
```

Try pressing the **RESET** button to see the application restart. Use `CTRL A` plus `X` to exit Minicom.

The **blink** application looks like this:

```go
package main

import (
	"os"
	"time"

	// Modules
	gpio "github.com/djthorpe/go-pico/pkg/gpio"
	uart "github.com/djthorpe/go-pico/pkg/uart"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

var (
	UARTConfig = uart.Config{BaudRate: 115200, DataBits: 8, StopBits: 1}
	LEDPin     = Pin(25)
	GPIOConfig = gpio.Config{Out: []Pin{LEDPin}}
)

func main() {
	// Create console
	stdout, err := UARTConfig.New()
	if err != nil {
		panic(err)
	}

	// Create GPIO
	gpio, err := GPIOConfig.New()
	if err != nil {
		stdout.Println(err)
		os.Exit(-1)
	}

	stdout.Println("loaded", gpio)

	// Blink lights
	for {
		gpio.High(LEDPin)
		time.Sleep(time.Millisecond * 800)
		gpio.Low(LEDPin)
		time.Sleep(time.Millisecond * 200)
	}
}
```

Two pico "devices" are used, the UART and the GPIO. The configuration
is defined right at the top (note the LED pin is listed as [GP25 here](https://datasheets.raspberrypi.com/pico/Pico-R3-A4-Pinout.pdf).
Then, within the `main` function, the devices are setup and an endless loop is entered
to switch on the LED for 800ms, then off for 200ms.

### Blink for the Raspberry Pi

Very similar code can also be written for the Raspberry Pi. An alternative
**blink** application can be compiled and run:

```bash
cd ${HOME}/projects/go-pico && make cmd/rpi/blink
./build/blink
```

This code expects an LED on GPIO22 on the Raspberry Pi:

| Wire     | Pin | Pin | Wire     |
|----------|-----|-----|----------|
|          |  1  |  2  |          |
|          |  3  |  4  |          |
|          |  5  |  6  | GND      |
|          |  7  |  8  | UART  TX |
|          |  9  | 10  | UART  RX |
|          | 11  | 12  |          |
|          | 13  | 14  | GND      |
| GPIO22   | 15  | 16  |          |
|          | 17  | 18  |          |
|          | 19  | 20  |          |
|          | 21  | 22  |          |
|          | 23  | 24  |          |
|          | 25  | 26  |          |
|          | 27  | 28  |          |
|          | 29  | 30  |          |
|          | 31  | 32  |          |
|          | 33  | 34  |          |
|          | 35  | 36  |          |
|          | 37  | 38  |          |
|          | 39  | 40  |          |

In order to wire up an LED, connect a resistor and LED in series between the GPIO22 and GND
pins. The orientation of the LED should be that the longer lead (anode) is connected through
the resistor to GPIO22.

The **blink** application is very similar to the Pico-targetted version:

```go
package main

import (
	"fmt"
	"os"
	"time"

	// Modules
	gpio "github.com/djthorpe/go-pico/pkg/gpio"

	// Namespace imports
	. "github.com/djthorpe/go-pico"
)

var (
	LEDPin     = Pin(22) // GPIO22
	GPIOConfig = gpio.Config{
		Out: []Pin{LEDPin},
	}
)

func main() {
	// Create GPIO
	gpio, err := GPIOConfig.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("loaded", gpio)

	// Blink lights
	for {
		gpio.High(LEDPin)
		time.Sleep(time.Millisecond * 800)
		gpio.Low(LEDPin)
		time.Sleep(time.Millisecond * 200)
	}
}
```

## Switch

The purpose of the "Switch" application is to determine if a switch is being clicked or released and perform
an action based on this event. In addition to the components for "blink" you'll need:

  * A PCB-mountable switch, ideally one which activates when pressed and deactivates when
    released.

TODO

## Further Examples

There are further examples for both the Pico and the Raspberry Pi. In
order to compile them, use the following commands:

```bash
make pico
make rpi
```

All the applications are stored in the `build` folder, as before.

