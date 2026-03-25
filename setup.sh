#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
PLATFORM=${1:-"linux"}
ARCH=${2:-"amd64"}
DOCKER_VERSION="v29.3.0"
mkdir -p dist/

/usr/bin/env bash ./build/download_docker_binary.sh "$PLATFORM" "$ARCH" "$DOCKER_VERSION"
