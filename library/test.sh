#!/bin/bash

set -e

cd $(dirname $0)/
export WORKINGDIR=$(pwd)

# build C library
cd $WORKINGDIR/c_library
mkdir -p build && cd build
if [[ -e "build.ninja" ]]; then
    meson setup .. --reconfigure --buildtype=debugoptimized
else
    meson setup .. --buildtype=debugoptimized
fi

ninja

export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:$(pwd)

# Run go test
cd ../../
go test -v -cover --count=1 .
