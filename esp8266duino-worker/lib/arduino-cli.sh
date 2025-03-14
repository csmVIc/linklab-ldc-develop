#!/bin/bash

ARCH=""

if [ $(uname -m) == "x86_64" ]
then 
  ARCH="x86_64"
elif [ $(uname -m) == "armv7l" ] 
then 
  ARCH="armv7l"
else
  echo "Not support this architecture $(uname -m)"
  exit -1
fi

echo ${ARCH}
mkdir -p /usr/local/bin
ln ./arduino-cli/${ARCH}/arduino-cli /usr/local/bin/arduino-cli


