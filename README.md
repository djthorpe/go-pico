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
  * Download **picotool** for flashing the code;
  * You'll need GNU Make to compile the example code in `cmd`.

There are probably lots of other dependencies I haven't listed here....

## Installation

Assuming a working operating system on your Raspberry Pi, create an `/opt`
folder for your installation if it's not already created, and install dependencies:

```bash
sudo apt -y install \
  git wget minicom \
  cmake gcc-arm-none-eabi libnewlib-arm-none-eabi libstdc++-arm-none-eabi-newlib \
  build-essential pkg-config libusb-1.0-0-dev
sudo install -o ${USER} -d /opt && cd /opt
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

To test, log out and back in again, and then check for **golang** and **tinygo** versions:

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
want to:

  1. Connect a reset button to the Pico
  2. Connect the default UART to your Raspberry Pi

