#!/bin/bash

# 用于Arduino烧写
apt-get install -y libusb-dev

# 用于ESP32DevKitC烧写
apt-get install -y python3 python3-pip
pip3 install --trusted-host=mirrors.aliyun.com -i https://mirrors.aliyun.com/pypi/simple/ --upgrade esptool

# 用于TelosB烧写
apt-get install -y python
if [ "$(dpkg --print-architecture)" == "amd64" ]
then
  apt install -y /app/tool/binutils-msp430_2.24_ti+exp0_amd64.deb
elif [ "$(dpkg --print-architecture)" == "armhf" ] || [ "$(dpkg --print-architecture)" == "armel" ]
then 
  apt install -y /app/tool/binutils-msp430_2.24_ti+exp0_armhf.deb
else
  echo "Not support this architecture $(dpkg --print-architecture)"
  exit -1
fi 
unzip /app/tool/TelosB.zip -d /app/tool/TelosB
rm /app/tool/TelosB.zip
rm /app/tool/binutils-msp430_2.24_ti+exp0_armhf.deb
rm /app/tool/binutils-msp430_2.24_ti+exp0_amd64.deb

# 用于DeveloperKit烧写
apt-get install -y build-essential cmake libusb-1.0
unzip /app/tool/stlink-1.5.1.zip -d /app/tool
cd /app/tool/stlink-1.5.1 && mkdir build && cd build && cmake -DCMAKE_BUILD_TYPE=Release .. && make install
rm /app/tool/stlink-1.5.1.zip

# 用于Haas100烧写
apt-get install -y python3 python3-pip
pip3 install --trusted-host=mirrors.aliyun.com -i https://mirrors.aliyun.com/pypi/simple/ --upgrade pyserial
apt-get install -y python python-pip
pip install --trusted-host=mirrors.aliyun.com -i https://mirrors.aliyun.com/pypi/simple/ --upgrade pyserial
apt-get install -y tftp
 
# 用于STM32F103C8烧写
unzip /app/tool/STM32F103C8/stm32isp.zip -d /app/tool/STM32F103C8/stm32isp
cd /app/tool/STM32F103C8/stm32isp && make
cp /app/tool/STM32F103C8/stm32isp/stm32isp /app/tool/STM32F103C8/tmp
rm -r /app/tool/STM32F103C8/stm32isp 
mv /app/tool/STM32F103C8/tmp /app/tool/STM32F103C8/stm32isp

# 用于HaasEDUK1烧写

# 重新加载安装的库
ldconfig