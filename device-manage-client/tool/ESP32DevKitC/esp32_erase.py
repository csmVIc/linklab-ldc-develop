#!/usr/bin/env python3
# -*- coding: UTF-8 -*-

import os
import time
import esp32_reboot

FLASH_CMD_TEMPLATE = "esptool.py --chip esp32 --port {} erase_flash"

def main():

    for devname in os.listdir("/dev"):
        if devname.startswith("ESP32DevKitC"):
            device = "/dev/" + devname

            # REBOOT
            esp32_reboot.reboot(device, 115200)

            # SLEEP
            time.sleep(1.0)

            # BURN
            cmdstr = FLASH_CMD_TEMPLATE.format(device)
            print(cmdstr)
            cmdres = os.system(cmdstr)
            if cmdres != 0:
                raise Exception("{} {} error".format(cmdstr, cmdres))


if __name__ == "__main__":
    main()
