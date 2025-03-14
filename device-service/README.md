# LinkLab Device Control V2 Device Service

## 测试

### 客户端

成功登录

```shell
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"clientid":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","hash":"pesOfc4A0ciy0dAJ69qTKHczLc4YTK/Cty0fdx2AoZk="}' \
  http://localhost:8081/api/client/login
```

密码错误

用户未注册

json参数绑定错误

### 设备状态更新

```shell
curl --header "Content-Type: application/json" \
  --cookie "clientid=8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92;token=8b8ad19199bd34826092ccd4329f657721069eba0908a206a50d899faf53c238" \
  --request POST \
  --data '{"add":["ESP32DevKitC-1.1.2.1"],"heartbeat":["ArduinoMega2560-1.3.4.4"],"sub":["ESP32DevKitC-1.1.2.4.2"]}' \
  http://localhost:8081/api/device/status
```

### 设备烧写结果

```shell
curl --header "Content-Type: application/json" \
  --cookie "clientid=8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92;token=841c456a782ff8abb8ca96b39110406d8da3601d102ba81cf91b578865ed8040" \
  --request POST \
  --data '{"groupid":"fc2620d393b66f6df6f7cdfb2dea6a89e4e246fbf5b50dc7c64ad30c73589199","userid":"yangg","deviceid":"ESP32DevKitC-1.1.2.1","taskindex":1,"success":1,"msg":"burn success"}' \
  http://localhost:8081/api/device/burn

curl --header "Content-Type: application/json" \
  --cookie "clientid=8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92;token=841c456a782ff8abb8ca96b39110406d8da3601d102ba81cf91b578865ed8040" \
  --request POST \
  --data '{"groupid":"fc2620d393b66f6df6f7cdfb2dea6a89e4e246fbf5b50dc7c64ad30c73589199","userid":"yangg","deviceid":"ESP32DevKitC-1.1.2.1","taskindex":1,"success":-1,"msg":"burn failed"}' \
  http://localhost:8081/api/device/burn
```

### 设备运行日志


curl --header "Content-Type: application/json" \
  --cookie "clientid=8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92;token=3c0833ea0a60e10daae63d20b71d42a986daff077c0aff7ddc7093a3fc60c9f0" \
  --request POST \
  --data '{"groupid":"b00f49de50bc5c83ac63170c50759cad39fa64a126926bc7f47969f76ddbfdee","userid":"yangg","deviceid":"ESP32DevKitC-1.1.2.1","taskindex":1,"msg":"hello world!"}' \
  http://localhost:8081/api/device/log

### 设备结束运行

curl --header "Content-Type: application/json" \
  --cookie "clientid=8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92;token=1fa91138c39c43ee3e356ab80f4c0c70dfe298a52c0c920aae96a6f74b6b4c4b" \
  --request POST \
  --data '{"groupid":"fc2620d393b66f6df6f7cdfb2dea6a89e4e246fbf5b50dc7c64ad30c73589199","userid":"yangg","deviceid":"ESP32DevKitC-1.1.2.1","taskindex":1,"timeout":1}' \
  http://localhost:8081/api/device/end


$SYS/brokers/+/clients/+/connected
$SYS/brokers/+/clients/+/disconnected


