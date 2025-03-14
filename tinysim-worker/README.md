




# 正确代码测试
curl -v --request POST --form 'parameters={"filehash":"a653d89d3e679d1f73605af7f2c9dee3e0b70d27", "boardType":"linuxhost", "compileType":"tinysim"};type=application/json' --form "file=@bin/tinysim.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v -o bin/result.zip http://10.214.149.214:30822/api/compile/block\?filehash\="a653d89d3e679d1f73605af7f2c9dee3e0b70d27"\&boardtype\="linuxhost"\&compiletype\="tinysim"

# 错误代码测试

curl -v --request POST --form 'parameters={"filehash":"b35b0ea2bff288ed2710ec241ddd3d79b0aa9c68", "boardType":"linuxhost", "compileType":"tinysim"};type=application/json' --form "file=@bin/tinysim-error.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v http://10.214.149.214:30822/api/compile/block\?filehash\="b35b0ea2bff288ed2710ec241ddd3d79b0aa9c68"\&boardtype\="linuxhost"\&compiletype\="tinysim"
