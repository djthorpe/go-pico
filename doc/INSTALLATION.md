
# Installation

The repository installation has been tested on macOS (arm64 and aarch64) and Linux (arm).
For macOS, it is assumed [homebrew](https://brew.sh/) is installed. 

Firstly, install the dependencies:

```zsh
# Installation of dependencies for macOS
brew install cmake libusb git wget
brew tap ArmMbed/homebrew-formulae
brew install arm-none-eabi-gcc
```

```zsh
# Installation of dependencies for Fedora Linux
sudo dnf install gcc-arm-linux-gnu \
  arm-none-eabi-gcc-cs-c++ arm-none-eabi-gcc-cs \
  arm-none-eabi-binutils arm-none-eabi-newlib \
  git wget build-essential cmake pkg-config libusb-1.0-0-dev 
```

## Install

There is then a script ypu can use to install the remaining dependencies:

```zsh
GOPICO="0.0.1"
cd ${HOME} && wget https://github.com/djthorpe/go-pico/archive/refs/tags/${GOPICO}.zip && unzip ${GOPICO}.zip
go-pico-${GOPICO}/scripts/install-tinygo.sh
```

Then, add the following to your .profile

```zsh
if [ -d "/opt/tinygo" ] ; then
  export PATH="${PATH}:/opt/tinygo/bin"
fi
if [ -d "/opt/bin" ] ; then
  export PATH="${PATH}:/opt/bin"
fi
if [ -d "/opt/pico-sdk" ] ; then
  export PICO_SDK_PATH="/opt/pico-sdk"
fi
```

## Test

To test you installation, log out and back in again, and then check for **tinygo** and **picotool** versions:

```bash
bash% tinygo version
tinygo version 0.24.0 darwin/amd64 (using go version go1.18.1 and LLVM version 14.0.0)
bash% picotool version
picotool v1.1.0 (Darwin 21.5.0, AppleClang-13.1.6.13160021, Release)
```
