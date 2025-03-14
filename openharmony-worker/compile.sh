#!/bin/bash
python3 -m pip uninstall ohos-build -y
pip3 install build/lite
echo 'export PATH=/root/.local/bin:$PATH' >>  /etc/profile 
echo 'export PATH=/root/gcc_riscv32/bin:$PATH' >>  /etc/profile 
source /etc/profile 
hb build -f
