# LinkLab Device Control V2 User Service

## 部署

### docker-compose

docker network create device_control_v2

## 测试

### 用户

#### 登录

成功登录

```shell
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"userid":"yangg","hash":"pesOfc4A0ciy0dAJ69qTKHczLc4YTK/Cty0fdx2AoZk="}' \
  http://localhost:8080/api/user/login
```

密码错误

```shell
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"userid":"yangg","hash":"3nCpU2Cbx91ghKhtrWhkbqHuUlVulW/lPbgcpOJ6+Y="}' \
  http://localhost:8080/api/user/login
```

用户未注册

```shell
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"userid":"yanggggggg","hash":"P3nCpU2Cbx91ghKhtrWhkbqHuUlVulW/lPbgcpOJ6+Y="}' \
  http://localhost:8080/api/user/login
```

json参数绑定错误

```shell
curl -v --header "Content-Type: application/json" \
  --request POST \
  --data '{"hash":"P3nCpU2Cbx91ghKhtrWhkbqHuUlVulW/lPbgcpOJ6+Y="}' \
  http://localhost:8080/api/user/login
```

### 设备

#### 文件上传

```shell
curl -v --request POST \
  --cookie "userid=yangg;token=3f167783af9c6b6414f2265f596bc55b154c27d5357f9d9499930c326a8f2293" \
  --form 'parameters={"boardtype":"ESP32DevKitC"};type=application/json' \
  --form "file=@test.bin;type=application/octet-stream" \
  http://localhost:8080/api/device/file
```

curl -v --header "Content-Type: application/json" \
  --request POST \
  --data '{"userid":"yangg","hash":"pesOfc4A0ciy0dAJ69qTKHczLc4YTK/Cty0fdx2AoZk="}' \
  http://localhost:8080/api/user/login


curl -v --request POST \
  --cookie "userid=yangg;token=dfa4be56bd841ca5ec342a363df2ffea11cfbfe211fbd63ab41463505a173c7a" \
  --data '{"tasks":[{"boardname":"ESP32DevKitC","deviceid":"ESP32DevKitC-1.1.2.1","runtime":30,"filehash":"e634b5ff7c980d68206f40a8d785af58bb0d901ce71d7fa5ee7791fc5d0f42e8","clientid":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","taskindex":1}]}' \
  http://localhost:8080/api/device/burn


curl -v --request POST \
  --cookie "userid=yangg;token=3f167783af9c6b6414f2265f596bc55b154c27d5357f9d9499930c326a8f2293" \
  --data '{"tasks":[{"boardname":"ESP32DevKitC","deviceid":"ArduinoMega2560-1.3.4.4","runtime":30,"filehash":"e634b5ff7c980d68206f40a8d785af58bb0d901ce71d7fa5ee7791fc5d0f42e8","clientid":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","taskindex":1}]}' \  
  http://localhost:8080/api/device/burn


curl -v --request POST \
  --cookie "userid=yangg;token=3f167783af9c6b6414f2265f596bc55b154c27d5357f9d9499930c326a8f2293" \
  --data '{"tasks":[{"boardname":"ESP32DevKitC","deviceid":"ArduinoMega2560-1.3.4.4","runtime":30,"filehash":"e634b5ff7c980d68206f40a8d785af58bb0d901ce71d7fa5ee7791fc5d0f42e8","clientid":"","taskindex":1}]}' \
  http://localhost:8080/api/device/burn

curl -v --request POST \
  --cookie "userid=yangg;token=7509100a6313d36643c36d1a45c7288b30f4a59255af44890497b5e7861d7d82" \
  --data '{"tasks":[{"boardname":"ESP32DevKitC","deviceid":"","runtime":30,"filehash":"032adb045a46ed3297c90fdf4aadd57aabbcb31f44839dcc8705784db770124c","clientid":"","taskindex":1}]}' \
  http://localhost:8080/api/device/burn




```c
{
    "tasks": [
        {
            "boardname": "ESP32DevKitC",
            "deviceid": "ESP32DevKitC-1.1.2.1",
            "runtime": 30,
            "filehash": "e634b5ff7c980d68206f40a8d785af58bb0d901ce71d7fa5ee7791fc5d0f42e8",
            "clientid": "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
            "taskindex": 1
        }
    ]
}
```



ef55bf4f6eb3673e9c0c7b08eed5823407270177fe1cee1bae38c310ef8b448b


curl -v --request POST \
  --header "Authorization: ef55bf4f6eb3673e9c0c7b08eed5823407270177fe1cee1bae38c310ef8b448b" \
  --form 'parameters={"boardname":"ESP32DevKitC"};type=application/json' \
  --form "file=@test.bin;type=application/octet-stream" \
  http://localhost:8080/api/device/file


curl -v \
  --header "Authorization: dd8c0e7e7f5f01fab212cf35b9e8e9c4f0404b7d7beb27fa51c6b2e8a00b023a" \
  http://localhost:8080/api/device/list?boardname="all"

  curl -v \
  --header "Authorization: faff67ef90c9f2181a58f4bd983fe4dbfd38e8dc32035aad2396ea2ef98b21c3" \
  http://localhost:8080/api/device/list?boardname="all"



curl -v \
  --header "Authorization: faff67ef90c9f2181a58f4bd983fe4dbfd38e8dc32035aad2396ea2ef98b21c3" \
  http://localhost:8080/api/board/list


curl -v --request POST \
  --header "Authorization: 70c91a5eae4be9886124d47e77af774d8517a8d804430f2c09233e01cd1c157b" \
  --data '{"tasks":[{"boardname":"ESP32DevKitC","deviceid":"","runtime":30,"filehash":"6ec0d4238b7164784a62f4b163c712c480b45b2a931ed6c2c6b00e4c66890ca1","clientid":"","taskindex":1}]}' \
  http://localhost:8080/api/device/burn


curl -v \
  --header "Authorization: d6b804c0169fdbc0952dc8ef54a2a147d059438e70ec03eee05762913801fd9d" \
  http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/board/list

curl -v \
  --header "Authorization: d6b804c0169fdbc0952dc8ef54a2a147d059438e70ec03eee05762913801fd9d" \
  http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/list?boardname="all"


curl -v --request POST \
  --header "Authorization: faff67ef90c9f2181a58f4bd983fe4dbfd38e8dc32035aad2396ea2ef98b21c3" \
  --data '{"tasks":[{"boardname":"ArduinoMega2560","deviceid":"/dev/ArduinoMega2560-8","runtime":30,"filehash":"eb3920b037e505b19c9a0ce0d8f28ae56f5ed28d9f70830ed22e11fd07d01c82","clientid":"ClientTest","taskindex":1},{"boardname":"ESP32DevKitC","deviceid":"/dev/ESP32DevKitC-0","runtime":30,"filehash":"6ec0d4238b7164784a62f4b163c712c480b45b2a931ed6c2c6b00e4c66890ca1","clientid":"ClientTest","taskindex":2}]}' \
  http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/burn