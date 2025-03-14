#!/usr/bin/env python3
# -*- coding: UTF-8 -*-

import os
import time
import esp32_reboot

FLASH_CMD_TEMPLATE = "esptool.py --chip esp32 --port {} --baud 921600 --before default_reset --after hard_reset write_flash -z --flash_mode dio --flash_freq 40m --flash_size detect 0x1000 tool/ESP32DevKitC/bin/alios/bootloader.bin 0x8000 tool/ESP32DevKitC/bin/alios/custom_partitions.bin 0x10000 {}"
MQTT_FIRMWARE_PATH = "./tool/ESP32DevKitC/bin/setup/mqtt.bin"
HELLO_FIRMWARE_PATH = "./tool/ESP32DevKitC/bin/setup/hello.bin"

def main():

    for devname in os.listdir("/dev"):
        if devname.startswith("ESP32DevKitC"):
            device = "/dev/" + devname

            # # REBOOT
            # esp32_reboot.reboot(device, 115200)
            
            # # BURN
            # cmdstr = FLASH_CMD_TEMPLATE.format(
            #     device, MQTT_FIRMWARE_PATH)
            # print(cmdstr)
            # cmdres = os.system(cmdstr)
            # if cmdres != 0:
            #     raise Exception("{} {} error".format(cmdstr, cmdres))

            # REBOOT
            esp32_reboot.reboot(device, 115200)

            # SLEEP
            time.sleep(1.0)

            # BURN
            cmdstr = FLASH_CMD_TEMPLATE.format(
                device, HELLO_FIRMWARE_PATH)
            print(cmdstr)
            cmdres = os.system(cmdstr)
            if cmdres != 0:
                raise Exception("{} {} error".format(cmdstr, cmdres))


if __name__ == "__main__":
    main()
