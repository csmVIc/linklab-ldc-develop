#!/usr/bin/env python3
# -*- coding: UTF-8 -*-

import os
import time
import random
import string
import pyboard
import shutil
import secrets
import argparse
import zipfile
import haas100_reboot

RANDOM_EXTRACT_DIR_LEN = 32


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
    cmd_parser.add_argument(
        "--timeout", default=10, help='device run timeout (s)')
    args = cmd_parser.parse_args()

    # 参数完整判断
    if len(args.device) < 1 or len(args.file) < 1:
        print("Error:Invalid params")
        print(
            "Usage: haas100_python_upload.py [-h] [--device /dev/Haas100] [--baudrate 1500000] [--file ./tmp/file.hex] [--timeout 10]")
        raise Exception("error")

    if zipfile.is_zipfile(args.file):
        # PYTHON ZIP
        # EXTRACT
        extractdir = os.path.join(args.extractdir, ''.join(secrets.choice(
            string.ascii_uppercase + string.digits) for _ in range(RANDOM_EXTRACT_DIR_LEN)))
        with zipfile.ZipFile(args.file, "r") as zip_ref:
            zip_ref.extractall(extractdir)
        # 扫描PYTHON文件
        pyfiles = [os.path.join(extractdir, fname) for fname in os.listdir(
            extractdir) if os.path.isfile(os.path.join(extractdir, fname))]
        if len(pyfiles) < 1:
            print("Error:The number of python files is empty")
            raise Exception("error")
        # REBOOT
        # haas100_reboot.reboot(args.device, args.baudrate)
        # SLEEP
        # time.sleep(5)
        # UPLOAD PYTHON
        try:
            pyb = pyboard.Pyboard(args.device, args.baudrate)
            pyb.enter_raw_repl()
            for pyfile in pyfiles:
                pyb.exec_with_timeout(pyfile, args.timeout)

        except Exception as err:
            # 退出时删除解压目录
            shutil.rmtree(extractdir)
            print("Error:Python {} exec error {}".format(
                pyfiles, err))
            raise Exception("error")
        # 退出时删除解压目录
        shutil.rmtree(extractdir)
    else:
        print("Error:Not support file type")
        raise Exception("error")


if __name__ == "__main__":
    main()
