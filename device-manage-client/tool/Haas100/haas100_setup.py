#!/usr/bin/env python3
# -*- coding: UTF-8 -*-

import os
import time
import pyboard
import haas100_reboot

MAX_RETRY = 5
FLASH_CMD_TEMPLATE = "python ./tool/Haas100/flash_program.py --serialport={} --bin={} --bin={}#0xB32000"
PYTHON_FILESYSTEM_PATH = "./tool/Haas100/littlefs.bin"
PYTHON_FIRMWARE_PATH = "./tool/Haas100/py_engine_demo@haas100.bin"
PYTHON_TEST_PATH = "./tool/Haas100/haas100_python.py"

# python flash_program.py --bin=./micropython_repl_demo@haas100.bin --bin=./littlefs.bin#0xB32000


def main():

    for devname in os.listdir("/dev"):
        if devname.startswith("Haas100"):
            device = "/dev/" + devname

            cmdstr = FLASH_CMD_TEMPLATE.format(
                device, PYTHON_FIRMWARE_PATH, PYTHON_FILESYSTEM_PATH)

            time.sleep(5)
            for index in range(MAX_RETRY):
                print("TRY COUNT ", index)
                print(cmdstr)

                # REBOOT
                haas100_reboot.reboot(device, 1500000)
                # SLEEP
                time.sleep(5)

                cmdres = os.system(cmdstr)
                if cmdres != 0 and index == MAX_RETRY-1:
                    raise Exception("{} {} error".format(cmdstr, cmdres))
                else:
                    break

                # SLEEP
                time.sleep(5)


if __name__ == "__main__":
    main()
