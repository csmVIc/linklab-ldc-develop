# 正确代码测试

mbed

NUMAKER_IOT_M487

c237287216bfae784a13b23692ae470131778901

curl -v --request POST --form 'parameters={"filehash":"c237287216bfae784a13b23692ae470131778901", "boardType":"NUMAKER_IOT_M487", "compileType":"mbed"};type=application/json' --form "file=@bin/helloworld.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v -o bin/helloworld.bin http://10.214.149.214:30822/api/compile/block\?filehash\="c237287216bfae784a13b23692ae470131778901"\&boardtype\="NUMAKER_IOT_M487"\&compiletype\="mbed" 

mbed compile -m NUMAKER_IOT_M487 -t GCC_ARM --source . --source ../mbed-os


# 错误代码测试

a8a579f223fdc13d479f6d921aaed350aedd5635

curl -v --request POST --form 'parameters={"filehash":"a8a579f223fdc13d479f6d921aaed350aedd5635", "boardType":"NUMAKER_IOT_M487", "compileType":"mbed"};type=application/json' --form "file=@bin/helloworld-error.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile









