
# Installation

The repository installation has been tested on Macintosh (arm64 and aarch64) and Linux (arm).
For Macintosh, it is assumed [homebrew](https://brew.sh/) is installed. 

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
