#!/bin/zsh

# You will want to install libusb

# Things you may want to change
PREFIX="/opt"
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
if [ "$(id -u)" != "0" ]; then
  echo "This script must be run as root"
  exit 1
fi

# Install the prefix
install -d -m 0755 "${PREFIX}" || exit -1

# Download tinygo, install
TINYGODEST="tinygo-${TINYGOBUILD}"
TINYGOSRC="${TINYGODEST}.${OS:l}-${ARCH:l}.tar.gz"
if [ -d "/opt/tinygo-${TINYGOBUILD}" ] ; then
  echo "${TINYGODEST} installed"
else
    echo "Downloading ${TINYGOSRC}"
    curl --silent --location --output "${TINYGOSRC}" --output-dir "${TEMP}" "https://github.com/tinygo-org/tinygo/releases/download/v${TINYGOBUILD}/${TINYGOSRC}" || exit -1
    install -d "/opt/${TINYGODEST}" || exit -1
    tar -C "/opt/${TINYGODEST}" -zxf "${TEMP}/${TINYGOSRC}" || exit -1
fi
if [ -d "${PREFIX}/${TINYGODEST}/tinygo" ]; then
  rm -f "${PREFIX}/tinygo" || exit -1
  pushd && cd "${PREFIX}" && ln -s "${PREFIX}/${TINYGODEST}/tinygo" && popd || exit -1
fi

# Download pico-sdk
PICOSDKDEST="pico-sdk-${PICOSDK}"
if [ -d "/opt/${PICOSDKDEST}" ] ; then
  echo "${PICOSDKDEST} installed"
else
  pushd && cd "${PREFIX}" && git clone -q -c advice.detachedHead=false --branch "${PICOSDK}" --single-branch https://github.com/raspberrypi/pico-sdk.git "${PICOSDKDEST}" && popd || exit -1
fi
if [ -d "${PREFIX}/${PICOSDKDEST}" ]; then
  rm -f "${PREFIX}/pico-sdk" || exit -1
  pushd && cd "${PREFIX}" && ln -s "${PREFIX}/${PICOSDKDEST}" pico-sdk && popd || exit -1
fi

# Download picotool
PICOTOOLDEST="picotool-${PICOTOOL}"
if [ -d "/opt/${PICOTOOLDEST}" ] ; then
  echo "${PICOTOOLDEST} installed"
else
  pushd && cd "${PREFIX}" && git clone -q -c advice.detachedHead=false --branch "${PICOTOOL}" --single-branch https://github.com/raspberrypi/picotool.git "${PICOTOOLDEST}" && popd || exit -1
fi
if [ -d "${PREFIX}/${PICOTOOLDEST}" ]; then
  rm -f "${PREFIX}/picotool" || exit -1
  pushd && cd "${PREFIX}" && ln -s "${PREFIX}/${PICOTOOLDEST}" picotool && popd || exit -1
fi

# Compile picotool
if [ -d "${PREFIX}/${PICOTOOLDEST}" ]; then
  pushd
  cd "${PREFIX}/${PICOTOOLDEST}"
  install -d build && cd build || exit -1
  PICO_SDK_PATH=/opt/pico-sdk cmake .. || exit -1
  make || exit -1
  popd
fi

# Install binaries
#/opt/picotool/build/picotool
#cd /opt/pico/picotool && install -d build && cd build &&  && make && install picotool /opt/pico/bin/picotool 

