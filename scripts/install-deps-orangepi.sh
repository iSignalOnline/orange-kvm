#!/bin/bash
# Install build dependencies for the Orange Pi Zero cross-compilation image.
# Unlike install-deps.sh (which pulls the Rockchip jetkvm-native-buildkit
# toolchain), this script installs only the standard Debian arm-linux-gnueabihf
# cross-compiler. The Orange Pi binary is built with CGO_ENABLED=0, so no
# hardware-vendor toolchain is required.

SUDO_PATH=$(which sudo)
function sudo() {
  if [ "$UID" -eq 0 ]; then
    "$@"
  else
    ${SUDO_PATH} "$@"
  fi
}

set -ex

export DEBIAN_FRONTEND=noninteractive

APT_PACKAGES=(
  build-essential
  gcc-arm-linux-gnueabihf
  g++-arm-linux-gnueabihf
  binutils-arm-linux-gnueabihf
  wget
  git
  ca-certificates
)

sudo apt-get update && \
    sudo apt-get install -y --no-install-recommends "${APT_PACKAGES[@]}" && \
    sudo rm -rf /var/lib/apt/lists/*
