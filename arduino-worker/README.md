
# LED 正确测试

curl -v --request POST --form 'parameters={"filehash":"c6f5073b870aa8ba0e359b74de0df761079599e5", "boardType":"arduino:avr:mega:cpu=atmega2560", "compileType":"arduino"};type=application/json' --form "file=@bin/arduino.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v -o bin/result.bin http://10.214.149.214:30822/api/compile/block\?filehash\="c6f5073b870aa8ba0e359b74de0df761079599e5"\&boardtype\="arduino:avr:mega:cpu=atmega2560"\&compiletype\="arduino"

# LED 错误测试

curl -v --request POST --form 'parameters={"filehash":"932cbc30b0444e78ca79cc972e2407608f847ba9", "boardType":"arduino:avr:mega:cpu=atmega2560", "compileType":"arduino"};type=application/json' --form "file=@bin/arduino-error.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v http://10.214.149.214:30822/api/compile/block\?filehash\="932cbc30b0444e78ca79cc972e2407608f847ba9"\&boardtype\="arduino:avr:mega:cpu=atmega2560"\&compiletype\="arduino"



