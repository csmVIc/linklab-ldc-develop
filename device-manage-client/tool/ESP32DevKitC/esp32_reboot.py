#!/usr/bin/env python3
# -*- coding: UTF-8 -*-
import time
import serial

def reboot(device, baudrate):
  handler = serial.Serial(port=device, baudrate=baudrate, timeout=1.0)
  handler.setRTS(True)
  handler.setDTR(False)
  time.sleep(0.1)
  handler.setDTR(True)
  handler.close()