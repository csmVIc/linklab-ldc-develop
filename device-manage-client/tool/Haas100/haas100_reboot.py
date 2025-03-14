#!/usr/bin/env python3
# -*- coding: UTF-8 -*-

import time
import serial


def reboot(device, baudrate):
    handler = serial.Serial(port=device, baudrate=baudrate, timeout=1.0)
    handler.write(b"\r\x02") 
    handler.write(b"\r\x04")
    handler.write("reboot\r\n".encode("utf-8"))
    handler.flush()
    # handler.reset_input_buffer()
    handler.close()
