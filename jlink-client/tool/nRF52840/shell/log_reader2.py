#!/usr/bin/env python3
# -*- coding: UTF-8 -*-

import sys
import time
import pylink
import argparse
import threading
from signal import signal, SIGINT

# jlink 全局变量
jlinkhandler = pylink.JLink(pylink.library.Library('/usr/local/lib/libjlinkarm.so'))

# output_flush 输出缓存刷新
def output_flush():
  while True:
    sys.stdout.flush()
    sys.stderr.flush()
    time.sleep(0.1)

# signal_handler 信号处理
def signal_handler(signal_received, frame):
    jlinkhandler.rtt_stop()
    if jlinkhandler.connected():
      jlinkhandler.close()
    sys.exit("close")

# input_handler 输入处理
def input_handler():
  canread = True
  while canread:
    try:
      cmdstr = input()
      # print("INPUT [{}]".format(cmdstr))
      writedata = list(bytearray(cmdstr, "utf-8") + b'\r\n')
      writeindex = 0
      while writeindex < len(writedata):
        bytes_written = jlinkhandler.rtt_write(0, writedata[writeindex:])
        writeindex = writeindex + bytes_written
        time.sleep(0.01)
    except EOFError as err:
      print("catch EOFError {}".format(err))
      canread = False

if __name__ == "__main__":

  # 输出缓存刷新线程
  ot = threading.Thread(target=output_flush)
  ot.start()

  # 解析参数
  parser = argparse.ArgumentParser(prog="Log reader app", description="use for read RTT log")
  parser.add_argument("-s", "--sn", dest="sn", type=int, required=True)
  args = parser.parse_args()
  print(args)

  # 初始化参数
  jlinkhandler.open(args.sn)
  jlinkhandler.set_tif(pylink.enums.JLinkInterfaces.SWD)
  jlinkhandler.connect("nRF52840_xxAA")
  if jlinkhandler.connected():
    print("SN = ", args.sn, " Connected")
  else:
    print("SN = ", args.sn, " Connection failed")
    sys.exit("error")
  jlinkhandler.reset(0, False)
  jlinkhandler.rtt_start()
  
  # 信号处理
  signal(SIGINT, signal_handler)

  # 输入处理
  it = threading.Thread(target=input_handler)
  it.start()

  # 主线程
  while jlinkhandler.connected():
    read_data = ''.join([chr(c) for c in jlinkhandler.rtt_read(0, 100)])
    if len(read_data) != 0:
      print(read_data, end='')
    else:
      time.sleep(0.1)