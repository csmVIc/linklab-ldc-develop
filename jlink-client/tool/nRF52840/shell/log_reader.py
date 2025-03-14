#!/usr/bin/env python3
# -*- coding: UTF-8 -*-

import Jlink
import argparse
import time
import sys
import threading

# output_flush 输出缓存刷新
def output_flush():
    while True:
        sys.stdout.flush()
        sys.stderr.flush()
        time.sleep(1)

if __name__ == '__main__':

    # 输出缓存刷新线程
    ot = threading.Thread(target=output_flush)
    ot.start()

    print("Log reader application")
    parser = argparse.ArgumentParser(prog="Log reader app", description="use for read RTT log")
    parser.add_argument("-v", "--verbose", dest="verbose", action="store_true")
    parser.add_argument("-s", "--sn", dest="sn", nargs='*', required=True)
    parser.add_argument("-p", "--prefix", dest="prefix")
    parser.add_argument("-r", "--reset", dest="reset", action="store_true")
    args = parser.parse_args()
    print(args)
    if args.sn[0] == "ALL":
        Jlink.log_start(args.verbose, args.prefix, None, all=True, sn_tag=False, reset=args.reset)
    else:
        try:
            sn_list = [int(s) for s in args.sn]
            print(sn_list)
            Jlink.log_start(args.verbose, args.prefix, sn_list, all=False, sn_tag=False, reset=args.reset)
        except BaseException:
            print("illegal sn!")
    count = 0
    while True:  # 保持主线程开启 才能保证后台线程的运行
        try:
            time.sleep(1)
            if count == 12:
                count = 0
                # print("Still awake")
            count += 1
        except KeyboardInterrupt:
            Jlink.kill_all()
            exit()
