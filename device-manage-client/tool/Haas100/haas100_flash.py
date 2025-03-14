#!/usr/bin/env python3
# -*- coding: UTF-8 -*-

import os
import sys
import time
import secrets
import string
import shutil
import random
import zipfile
import pyboard
import argparse
import haas100_reboot

PREFIX_COUNT = 8
RANDOM_EXTRACT_DIR_LEN = 32
BIN_PREFIX = b'\xff\xff\xff\xff\x01\x00\x00\x00'
FLASH_CMD_TEMPLATE = "python ./tool/Haas100/flash_program.py --serialport={} --bin={}"
PYTHON_FIRMWARE_PATH = "./tool/Haas100/mpy_cli_demo@haas100.bin"


def main():

    # 随机数
    random.seed(time.time())

    # 参数解析
    cmd_parser = argparse.ArgumentParser()
    cmd_parser.add_argument(
        '--file', default='', help='input file to be execed')
    cmd_parser.add_argument(
        '--device', default='', help='the serial device')
    cmd_parser.add_argument(
        '--baudrate', default=1500000, help='the baud rate of the serial device')
    cmd_parser.add_argument(
        "--extractdir", default='./tmp', help='python zip extract directory')
    args = cmd_parser.parse_args()

    # 参数完整判断
    if len(args.device) < 1 or len(args.file) < 1:
        print("Error:Invalid params")
        print(
            "Usage: haas100_flash.py [-h] [--device /dev/Haas100] [--baudrate 1500000] [--file ./tmp/file.hex] [--extractdir ./tmp]")
        raise Exception("error")

    # 烧写
    with open(args.file, "rb") as f:
        prefix_data = f.read(8)
        # 语言判断
        if prefix_data == BIN_PREFIX:
            # C
            # REBOOT
            haas100_reboot.reboot(args.device, int(args.baudrate))
            # SLEEP
            time.sleep(1)
            # BURN
            cmdstr = FLASH_CMD_TEMPLATE.format(args.device, args.file)
            print(cmdstr)
            cmdres = os.system(cmdstr)
            if cmdres != 0:
                raise Exception("error")
        else:
            print("Error:Not support file type")
            raise Exception("error")

if __name__ == "__main__":
    main()
