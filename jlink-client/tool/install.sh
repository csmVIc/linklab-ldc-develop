#!/bin/bash

# 用于nRF52840
apt-get install -y python3 python3-pip libusb-1.0 libusb-dev
apt-get install -y libfreetype6 libxrender1 libxrandr2 libxfixes-dev libxcursor1 libfontconfig1
if [ "$(dpkg --print-architecture)" == "amd64" ]
then
  apt install -y /app/tool/nRF52840/JLink_Linux_V758d_x86_64.deb
elif [ "$(dpkg --print-architecture)" == "armhf" ]
then 
  apt install -y /app/tool/nRF52840/JLink_Linux_V760a_arm.deb
else
  echo "Not support this architecture $(dpkg --print-architecture)"
  exit -1
fi 
pip3 install --trusted-host=mirrors.aliyun.com -i https://mirrors.aliyun.com/pypi/simple/ --upgrade pylink-square
cd /app/tool/nRF52840/pylink && python3 setup.py install && cp /opt/SEGGER/JLink/libjlinkarm.so /usr/local/lib/

# 重新加载安装的库
ldconfig