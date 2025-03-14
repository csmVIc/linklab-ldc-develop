#!/usr/bin/env python
import os, sys, re, codecs, time, json
from ymodem import YModem

try:
    import serial
    from serial.tools import miniterm
    from serial.tools.list_ports import comports
except:
    print("\n\nNot found pyserial, please install: \nsudo pip install pyserial")
    sys.exit(0)

def get_bin_file():
    """ get binary file from sys.argv --bin=/xxx/yyy/zzz.bin """ 
    bin_files = []
    pattern = re.compile(r'--(.*)=(.*)')
    for arg in sys.argv[1:]:
        if arg.startswith("--"):
            match = pattern.match(arg)
            if match:
                key = match.group(1)
                value = match.group(2)
                if key == 'bin':
                    bin_files.append(value)
    return bin_files

def get_serial_port():
    serialport = ""
    pattern = re.compile(r'--(.*)=(.*)')
    for arg in sys.argv[1:]:
        if arg.startswith("--"):
            match = pattern.match(arg)
            if match:
                key = match.group(1)
                value = match.group(2)
                if key == 'serialport':
                    serialport = value
                    return serialport
    return serialport

def read_json(json_file):
    data = None
    if os.path.isfile(json_file):
        with open(json_file, 'r') as f:
            data = json.load(f)
    return data

def write_json(json_file, data):
    with open(json_file, 'w') as f:
        f.write(json.dumps(data, indent=4, separators=(',', ': ')))

def get_config():
    """ get configuration from .config_burn file, if it is not existed, 
        generate default configuration of chip_haas1000 """ 
    configs = {}
    config_file = os.path.join(os.getcwd(), '.config_burn')
    if os.path.isfile(config_file):
        configs = read_json(config_file)
        if not configs:
            configs = {}
    if 'chip_haas1000' not in configs:
        configs['chip_haas1000'] = {}
    if 'serialport' not in configs['chip_haas1000']:
        configs['chip_haas1000']['serialport'] = ""
    if 'baudrate' not in configs['chip_haas1000']:
        configs['chip_haas1000']['baudrate'] = "1500000"
    if 'binfile' not in configs['chip_haas1000']:
        configs['chip_haas1000']['binfile'] = []

    return configs['chip_haas1000']

def save_config(config):
    """ save configuration to .config_burn file, only update chip_haas1000 portion """ 
    if config:
        configs = {}
        config_file = os.path.join(os.getcwd(), '.config_burn')
        if os.path.isfile(config_file):
            configs = read_json(config_file)
            if not configs:
                configs = {}
        configs['chip_haas1000'] = config
        write_json(config_file, configs)

def send_cmd_check_recv_data(serialport, command, pattern, timeout):
    """ receive serial data, and check it with pattern """
    if command:
        serialport.write(command)
    matcher = re.compile(pattern)
    tic     = time.time()
    buff    = serialport.read(128)
    while (time.time() - tic) < timeout:
        buff += serialport.read(128)
        if matcher.search(buff):
            return True, buff
        else:
            if command:
                serialport.write(command)
    return False, buff

def burn_bin_file(serialport, bin_file):
    filename = bin_file
    address = "0"
    if "#" in bin_file:
        filename = bin_file.split("#", 1)[0]
        address = bin_file.split("#", 1)[1]
    for i in range(3):
        # ymodem update
        serialport.write(b'1')
        time.sleep(0.1)

        # get flash address
        bmatched, buff = send_cmd_check_recv_data(serialport, b'', b'Please input flash addr:', 2)
        print(buff)
        if bmatched:
            break
    if address == "0":
        if bmatched:
            pattern = re.compile(b'Backup part addr:([0-9a-fxA-F]*)')
            match = pattern.search(buff)
            if match:
                address = match.group(1)
            else:
                print("can not get flash address")
                return False
        else:
            print("can not get flash address")
            return False
    else:
        address = address.encode()

    # set flash address
    serialport.write(address)
    serialport.write(b'\r\n')
    time.sleep(0.1)
    bmatched, buff = send_cmd_check_recv_data(serialport, b'', b'CCCCC', 5)
    print(buff)
    if not bmatched:
        print("can not enter into ymodem mode")
        return False

    # send binary file
    def sender_getc(size):
        return serialport.read(size) or None

    def sender_putc(data, timeout=15):
        return serialport.write(data)

    sender = YModem(sender_getc, sender_putc)
    sent = sender.send_file(filename)
    return True


def burn_bin_files(portnum, baudrate, bin_files):
    # open serial port
    serialport = serial.Serial()
    serialport.port = portnum
    serialport.baudrate = baudrate
    serialport.parity = "N"
    serialport.bytesize = 8
    serialport.stopbits = 1
    serialport.timeout = 0.05

    try:
        serialport.open()
    except Exception as e:
        raise Exception("Failed to open serial port: %s!" % portnum)

    # reboot the board, and enter into 2nd boot mode
    bmatched = False
    for i in range(300):
        serialport.write(b'\nreboot\n\n')
        time.sleep(0.1)
        bmatched, buff = send_cmd_check_recv_data(serialport, b'w', b'aos boot#', 2)
        print(buff)
        if bmatched:
            break
        if i > 3:
            print("Please reboot the board manually.")

    if not bmatched:
        print("Please reboot the board manually, and try it again.")
        serialport.close()
        return

    for bin_file in bin_files:
        if not burn_bin_file(serialport, bin_file):
            print("Download file %s failed." % bin_file)
            serialport.close()
            return

    # switch partition
    print("Swap AB partition");
    serialport.write(b'3')
    time.sleep(0.1)
    serialport.write(b'4')
    time.sleep(0.1)
    serialport.write(b'3')
    time.sleep(0.1)
    serialport.write(b'2')
    time.sleep(0.1)
    bmatched, buff = send_cmd_check_recv_data(serialport, b'', b'2ndboot cli menu in 100ms', 2)
    print(buff)
    if bmatched:
        print("Burn \"%s\" success." % bin_files)

    # close serial port
    serialport.close()

def main():
    # step 1: get binary file
    needsave = False
    myconfig = get_config()
    bin_files = get_bin_file()
    if bin_files:
        myconfig["binfile"] = bin_files
        needsave = True
    if not myconfig["binfile"]:
        print("no specified binary file")
        return
    print("binary file is %s" % myconfig["binfile"])

    # step 2: get serial port
    myconfig["serialport"] = get_serial_port()
    if not myconfig["serialport"]:
        myconfig["serialport"] = miniterm.ask_for_port()
        if not myconfig["serialport"]:
            print("no specified serial port")
            return
        else:
            needsave = True

    print("serial port is %s" % myconfig["serialport"])
    print("the settings were restored in the file %s" % os.path.join(os.getcwd(), '.config_burn'))

    # step 3: burn binary file into flash
    burn_bin_files(myconfig["serialport"], myconfig['baudrate'], myconfig["binfile"])
    if needsave:
        save_config(myconfig)

if __name__ == "__main__":
    main()
