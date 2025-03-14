




# 正确代码测试

curl -v --request POST --form 'parameters={"filehash":"98bca5e26f43055315c81dc79cda22d29950f3d2", "boardType":"esp32devkitc", "compileType":"alios"};type=application/json' --form "file=@bin/alios.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v -o bin/alios.bin http://10.214.149.214:30822/api/compile/block\?filehash\="98bca5e26f43055315c81dc79cda22d29950f3d2"\&boardtype\="esp32devkitc"\&compiletype\="alios" 

# 正确代码测试

curl -v --request POST --form 'parameters={"filehash":"98bca5e26f43055315c81dc79cda22d29950f3d2", "boardType":"developerkit", "compileType":"alios"};type=application/json' --form "file=@bin/alios.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v -o bin/alios.bin http://10.214.149.214:30822/api/compile/block\?filehash\="98bca5e26f43055315c81dc79cda22d29950f3d2"\&boardtype\="developerkit"\&compiletype\="alios" 

# 错误代码测试

curl -v --request POST --form 'parameters={"filehash":"e6f662951378bb1c7d9a79964e07279e6f6e98ad", "boardType":"esp32devkitc", "compileType":"alios"};type=application/json' --form "file=@bin/alios-error.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v http://10.214.149.214:30822/api/compile/block\?filehash\="e6f662951378bb1c7d9a79964e07279e6f6e98ad"\&boardtype\="esp32devkitc"\&compiletype\="alios" 
