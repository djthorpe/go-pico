#!/bin/zsh

# You will want to install libusb
#
# mac: 
#  brew install cmake libusb
#  brew tap ArmMbed/homebrew-formulae
#  brew install arm-none-eabi-gcc

# Things you may want to change
PREFIX="${HOME}/opt"
TINYGOBUILD="0.23.0"
PICOSDK="1.3.1"
PICOTOOL="1.1.0"

# Determine OS, Arch and temp directory
OS="$(uname -s)"
ARCH="$(uname -m)"
TEMP=`mktemp -d`

# Adjust arch
if [ "${ARCH}" = "x86_64" ]; then
  ARCH="amd64"
elif [ "${ARCH}" = "aarch64" ]; then
  ARCH="arm64"
elif [ "${ARCH}" = "armv7l" ]; then
  ARCH="arm"
fi

# Fudge for Mac M1/M2
if [ "${ARCH}" = "arm64" ] && [ "${OS:l}" = "darwin" ]; then
  ARCH="amd64"
fi

# Check for curl and git
if [ ! -x "$(command -v curl)" ]; then
  echo "curl is not installed. Please install curl before continuing."
  exit 1
fi
if [ ! -x "$(command -v git)" ]; then
  echo "git is not installed. Please install git before continuing."
  exit 1
fi
if [ ! -x "$(command -v make)" ]; then
  echo "make is not installed. Please install make before continuing."
  exit 1
fi
if [ ! -x "$(command -v cmake)" ]; then
  echo "cmake is not installed. Please install cmake before continuing."
  exit 1
fi

# Make temporary location, cleanup on exit
trap cleanup EXIT
cleanup() {
    rm -fr ${TEMP}
}

# Print out the variables
echo "prefix: ${PREFIX}"
echo "tinygo: ${TINYGOBUILD}"
echo "os: ${OS:l}"
echo "arch: ${ARCH}"
echo

# Check for root
#if [ "$(id -u)" != "0" ]; then
#  echo "This script must be run as root"
#  exit 1
#fi

# Install the prefix and bin folders
install -d -m 0755 "${PREFIX}/bin" || exit -1

# Download tinygo, install
TINYGODEST="tinygo-${TINYGOBUILD}"
TINYGOSRC="${TINYGODEST}.${OS:l}-${ARCH:l}.tar.gz"
if [ -d "${PREFIX}/tinygo-${TINYGOBUILD}" ] ; then
  echo "${TINYGODEST} installed"
else
    echo "Downloading ${TINYGOSRC}"
    curl --silent --location --output "${TINYGOSRC}" --output-dir "${TEMP}" "https://github.com/tinygo-org/tinygo/releases/download/v${TINYGOBUILD}/${TINYGOSRC}" || exit -1
    install -d "${PREFIX}/${TINYGODEST}" || exit -1
    tar -C "${PREFIX}/${TINYGODEST}" -zxf "${TEMP}/${TINYGOSRC}" || exit -1
fi
if [ -d "${PREFIX}/${TINYGODEST}/tinygo" ]; then
  rm -f "${PREFIX}/tinygo" || exit -1
  pushd && cd "${PREFIX}" && ln -s "${PREFIX}/${TINYGODEST}/tinygo" && popd || exit -1
fi

# Download pico-sdk
PICOSDK_SRC="https://github.com/raspberrypi/pico-sdk.git"
PICOSDK_DEST="pico-sdk-${PICOSDK}"
if [ -d "${PREFIX}/${PICOSDK_DEST}" ] ; then
  echo "${PICOSDK_DEST} installed"
else
  pushd && cd "${PREFIX}" && git clone -q -c advice.detachedHead=false --branch "${PICOSDK}" --single-branch "${PICOSDK_SRC}"  "${PICOSDK_DEST}" && popd || exit -1
fi
if [ -d "${PREFIX}/${PICOSDK_DEST}" ]; then
  rm -f "${PREFIX}/pico-sdk" || exit -1
  pushd && cd "${PREFIX}" && ln -s "${PREFIX}/${PICOSDK_DEST}" pico-sdk && popd || exit -1
fi

# Download picotool
PICOTOOL_SRC="https://github.com/raspberrypi/picotool.git"
PICOTOOL_DEST="picotool-${PICOTOOL}"
if [ -d "${PREFIX}/${PICOTOOL_DEST}" ] ; then
  echo "${PICOTOOL_DEST} installed"
else
  pushd && cd "${PREFIX}" && git clone -q -c advice.detachedHead=false --branch "${PICOTOOL}" --single-branch "${PICOTOOL_SRC}" "${PICOTOOL_DEST}" && popd || exit -1
fi
if [ -d "${PREFIX}/${PICOTOOL_DEST}" ]; then
  rm -f "${PREFIX}/picotool" || exit -1
  pushd && cd "${PREFIX}" && ln -s "${PREFIX}/${PICOTOOL_DEST}" picotool && popd || exit -1
fi


# Compile pico-sdk
if [ -d "${PREFIX}/${PICOSDK_DEST}" ]; then
  pushd
  cd "${PREFIX}/${PICOSDK_DEST}"
  git submodule update --init || exit -1
  install -d build && cd build || exit -1
  PICO_SDK_PATH="${PREFIX}/pico-sdk" cmake .. || exit -1
  make || exit -1
  popd
fi


# Compile picotool
if [ -d "${PREFIX}/${PICOTOOL_DEST}" ]; then
  pushd
  cd "${PREFIX}/${PICOTOOL_DEST}"
  git submodule update --init || exit -1
  install -d build && cd build || exit -1
  PICO_SDK_PATH="${PREFIX}/pico-sdk" cmake .. || exit -1
  make || exit -1
  popd
fi

# Install binaries
if [ -d "${PREFIX}/bin" ]; then
  pushd
  cd "${PREFIX}/bin"
  install "${PREFIX}/${PICOSDK_DEST}/build/elf2uf2/elf2uf2" elf2uf2 || exit -1
  install "${PREFIX}/${PICOTOOL_DEST}/build/picotool" picotool || exit -1
  popd
fi

