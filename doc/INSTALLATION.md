
# Installation

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
