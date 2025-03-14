# 正确测试

6ecb76d564ad7634e325d4330d17411d45f7f927

curl -v --request POST --form 'parameters={"filehash":"6ecb76d564ad7634e325d4330d17411d45f7f927", "boardType":"esp8266", "compileType":"esp8266duino"};type=application/json' --form "file=@bin/helloworld.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v -o bin/result.bin http://10.214.149.214:30822/api/compile/block\?filehash\="6ecb76d564ad7634e325d4330d17411d45f7f927"\&boardtype\="esp8266"\&compiletype\="esp8266duino"

# 错误测试

8252963ad6bd2abb132677b7674bbee742211b13

curl -v --request POST --form 'parameters={"filehash":"8252963ad6bd2abb132677b7674bbee742211b13", "boardType":"esp8266", "compileType":"esp8266duino"};type=application/json' --form "file=@bin/helloworld-error.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v http://10.214.149.214:30822/api/compile/block\?filehash\="8252963ad6bd2abb132677b7674bbee742211b13"\&boardtype\="esp8266"\&compiletype\="esp8266duino"

