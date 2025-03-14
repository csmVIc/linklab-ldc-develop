#!/bin/bash

# import hashlib
# import sys
# hashlib.sha1(open('source/alios.zip','rb').read()).hexdigest()

# curl -v --request POST --form 'parameters={"filehash":"aa84941827bb2c430656d040f2dbcf3f11cfb34d", "boardType":"esp32", "compileType":"esp32duino"};type=application/json' --form "file=@source/ESP32DevKitCArduino-MQTT.zip;type=application/octet-stream"  http://kubernetes.tinylink.cn/linklab/compilev2/api/compile

curl -v -o bin/ESP32DevKitCArduino-MQTT.bin http://kubernetes.tinylink.cn/linklab/compilev2/api/compile/block\?filehash\="aa84941827bb2c430656d040f2dbcf3f11cfb34d"\&boardtype\="esp32"\&compiletype\="esp32duino"

# curl -v --request POST --form 'parameters={"filehash":"c6f5073b870aa8ba0e359b74de0df761079599e5", "boardType":"arduino:avr:mega:cpu=atmega2560", "compileType":"arduino"};type=application/json' --form "file=@bin/arduino.zip;type=application/octet-stream"  http://api.testbd.tinylink.cn/linklab/compilev2/api/compile

# curl -v -o bin/result.zip http://api.testbd.tinylink.cn/linklab/compilev2/api/compile/block\?filehash\="c6f5073b870aa8ba0e359b74de0df761079599e5"\&boardtype\="arduino:avr:mega:cpu=atmega2560"\&compiletype\="arduino"