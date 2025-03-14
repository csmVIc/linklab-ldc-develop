# 正确代码测试

57ca639b42afdc378b19c5ac9e5e2b09975d603d

curl -v --request POST --form 'parameters={"filehash":"57ca639b42afdc378b19c5ac9e5e2b09975d603d", "boardType":"sky", "compileType":"contiki-ng"};type=application/json' --form "file=@bin/helloworld.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v -o bin/helloworld_demo.bin http://10.214.149.214:30822/api/compile/block\?filehash\="57ca639b42afdc378b19c5ac9e5e2b09975d603d"\&boardtype\="sky"\&compiletype\="contiki-ng"

# 错误代码测试

c2f6669293208125d8c82c7b367a0209bae79ab4

curl -v --request POST --form 'parameters={"filehash":"c2f6669293208125d8c82c7b367a0209bae79ab4", "boardType":"sky", "compileType":"contiki-ng"};type=application/json' --form "file=@bin/helloworld-error.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v http://10.214.149.214:30822/api/compile/block\?filehash\="c2f6669293208125d8c82c7b367a0209bae79ab4"\&boardtype\="sky"\&compiletype\="contiki-ng"
