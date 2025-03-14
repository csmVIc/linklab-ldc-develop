# helloworld 正确代码测试

curl -v --request POST --form 'parameters={"filehash":"9d3ac3b79d349bc30f0cbce3557911ddbbdc9343", "boardType":"haas100", "compileType":"alios-haas"};type=application/json' --form "file=@bin/helloworld_demo.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v -o bin/helloworld_demo.bin http://10.214.149.214:30822/api/compile/block\?filehash\="9d3ac3b79d349bc30f0cbce3557911ddbbdc9343"\&boardtype\="haas100"\&compiletype\="alios-haas" 

curl -v --request POST --form 'parameters={"filehash":"9d3ac3b79d349bc30f0cbce3557911ddbbdc9343", "boardType":"haas100", "compileType":"alios-haas"};type=application/json' --form "file=@bin/helloworld_demo.zip;type=application/octet-stream" http://kubernetes.tinylink.cn/linklab/compilev2/api/compile

# helloworld 错误代码测试

curl -v --request POST --form 'parameters={"filehash":"933377097f6ff19ebe9dcbdb06436463a298068a", "boardType":"haas100", "compileType":"alios-haas"};type=application/json' --form "file=@bin/helloworld_demo_error.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v http://10.214.149.214:30822/api/compile/block\?filehash\="933377097f6ff19ebe9dcbdb06436463a298068a"\&boardtype\="haas100"\&compiletype\="alios-haas" 

# mqtt 正确代码测试

curl -v --request POST --form 'parameters={"filehash":"d7607cd4f9f273417d5fd20805bac56b7145a258", "boardType":"haas100", "compileType":"alios-haas"};type=application/json' --form "file=@bin/mqtt_demo.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v -o bin/mqtt_demo.bin http://10.214.149.214:30822/api/compile/block\?filehash\="d7607cd4f9f273417d5fd20805bac56b7145a258"\&boardtype\="haas100"\&compiletype\="alios-haas"

# mqtt 错误代码测试

curl -v --request POST --form 'parameters={"filehash":"610e05b6f785b407f75d95c0ea26d8faa8cc26fe", "boardType":"haas100", "compileType":"alios-haas"};type=application/json' --form "file=@bin/mqtt_demo_error.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v http://10.214.149.214:30822/api/compile/block\?filehash\="610e05b6f785b407f75d95c0ea26d8faa8cc26fe"\&boardtype\="haas100"\&compiletype\="alios-haas"

# haas python 代码测试

curl -v --request POST --form 'parameters={"filehash":"3a17463750dc454517113be43fc9461f82ddb91c", "boardType":"haas100", "compileType":"alios-haas-python"};type=application/json' --form "file=@bin/haas100-python.zip;type=application/octet-stream"  http://10.214.149.214:30822/api/compile

curl -v -o bin/haas100-python-result.zip http://10.214.149.214:30822/api/compile/block\?filehash\="3a17463750dc454517113be43fc9461f82ddb91c"\&boardtype\="haas100"\&compiletype\="alios-haas-python" 





