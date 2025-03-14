#!/bin/bash

# 系统限制
sudo cp system.conf /etc/systemd/system.conf
sudo cp user.conf /etc/systemd/user.conf

# 用户限制
sudo cp limits.conf /etc/security/limits.conf

sudo reboot