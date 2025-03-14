# 正确测试

6ecb76d564ad7634e325d4330d17411d45f7f927

curl -v --request POST --form 'parameters={"filehash":"6ecb76d564ad7634e325d4330d17411d45f7f927", "boardType":"esp32", "compileType":"esp32duino"};type=application/json' --form "file=@bin/helloworld.zip;type=application/octet-stream"  http://kubernetes.tinylink.cn/linklab/compilev2/api/compile

curl -v -o bin/result.bin http://kubernetes.tinylink.cn/linklab/compilev2/api/compile/block\?filehash\="6ecb76d564ad7634e325d4330d17411d45f7f927"\&boardtype\="esp32"\&compiletype\="esp32duino"

# 错误测试

8252963ad6bd2abb132677b7674bbee742211b13

curl -v --request POST --form 'parameters={"filehash":"8252963ad6bd2abb132677b7674bbee742211b13", "boardType":"esp32", "compileType":"esp32duino"};type=application/json' --form "file=@bin/helloworld-error.zip;type=application/octet-stream"  http://kubernetes.tinylink.cn/linklab/compilev2/api/compile

curl -v http://10.214.149.214:30822/api/compile/block\?filehash\="8252963ad6bd2abb132677b7674bbee742211b13"\&boardtype\="esp32"\&compiletype\="esp32duino"

# mqtt测试
794b14e02a19877aa1d528a8d41207a812af4abd

curl -v --request POST --form 'parameters={"filehash":"794b14e02a19877aa1d528a8d41207a812af4abd", "boardType":"esp32", "compileType":"esp32duino"};type=application/json' --form "file=@bin/mqtt.zip;type=application/octet-stream"  http://kubernetes.tinylink.cn/linklab/compilev2/api/compile

curl -v -o bin/mqtt.bin http://kubernetes.tinylink.cn/linklab/compilev2/api/compile/block\?filehash\="794b14e02a19877aa1d528a8d41207a812af4abd"\&boardtype\="esp32"\&compiletype\="esp32duino"

7669f5eca74c92fadd85143de12a34ccceb49685

curl -v --request POST --form 'parameters={"filehash":"7669f5eca74c92fadd85143de12a34ccceb49685", "boardType":"esp32", "compileType":"esp32duino"};type=application/json' --form "file=@bin/mqtt.zip;type=application/octet-stream"  http://kubernetes.tinylink.cn/linklab/compilev2/api/compile

curl -v -o bin/mqtt.bin http://kubernetes.tinylink.cn/linklab/compilev2/api/compile/block\?filehash\="7669f5eca74c92fadd85143de12a34ccceb49685"\&boardtype\="esp32"\&compiletype\="esp32duino"

