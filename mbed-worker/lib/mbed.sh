#!/bin/bash

apt-get install -y gcc-arm-none-eabi
apt-get install -y python3 python3-pip git mercurial
mkdir ~/.pip
cp pip.conf ~/.pip/pip.conf
python3 -m pip install mbed-cli
pip3 install -r requirements.txt
mbed config -G GCC_ARM_PATH "/usr/bin"
mbed config -G MBED_OS_DIR "/app/workspace/mbed-os"


