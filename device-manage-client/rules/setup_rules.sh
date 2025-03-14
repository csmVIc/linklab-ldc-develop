#!/bin/bash

# 密码
BOARDNAME=EDURISCV64

# 绝对目录
cd "$(dirname "$0")"
pwd

# UDEV配置文件修改
sudo cp ./${BOARDNAME}.rules /etc/udev/rules.d

# UDEV配置文件生效
sudo udevadm control --reload-rules
sudo udevadm trigger
