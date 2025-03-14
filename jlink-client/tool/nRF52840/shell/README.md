# jlink 扫描
pylink emulator -l usb

# jlink 烧写
python3 mass_flash.py -f RIOT-1-32.bin -s 3
python3 mass_flash.py -f RIOT-Shell.bin -s 3

# jlink 日志
python3 log_reader.py -s 3 -v -r

# jlink 命令
python3 rtt_view.py -s 3 -t -v
