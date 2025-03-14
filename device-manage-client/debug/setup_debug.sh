#!/bin/bash

# 密码
PASSWORD=raspberry

# 绝对目录
cd "$(dirname "$0")"
pwd

# 串口调试
echo ${PASSWORD} | sudo -S apt-get install -y minicom
echo ${PASSWORD} | sudo -S cp minicom/minirc.dfl /etc/minicom/minirc.dfl

# ESP32DevKitC
echo ${PASSWORD} | sudo -S apt-get install -y python3 python3-pip
pip3 install --trusted-host=mirrors.aliyun.com -i https://mirrors.aliyun.com/pypi/simple/ --upgrade esptool

# Haas100
echo ${PASSWORD} | sudo -S apt-get install -y python3 python3-pip
pip3 install --trusted-host=mirrors.aliyun.com -i https://mirrors.aliyun.com/pypi/simple/ --upgrade pyserial
echo ${PASSWORD} | sudo -S apt-get install -y python python-pip
pip install --trusted-host=mirrors.aliyun.com -i https://mirrors.aliyun.com/pypi/simple/ --upgrade pyserial