
```shell
curl -v --request POST \
  --header "Authorization: bfda4d2b099812e8733a9f704b318ed3e9a460d579f5d34252064d9cf2bb66b3" \
  --form 'parameters={"boardname":"ESP32DevKitC"};type=application/json' \
  --form "file=@test.bin;type=application/octet-stream" \
  http://localhost:8083/api/file
```

```shell
curl -v --request GET \
  --header "Authorization: 06448cec9faf15ae3d083a26f47b26f8de94be06a0f350c9340fc21c91637cd6" \
  http://localhost:8083/api/file\?filehash\="3698042c9054781fc72c3a89bb0dd881bc3e96ba93ae3f8f571c170d11f5d0e6"\&boardname\="ESP32DevKitC"
```

curl -v --request POST \
  --header "Authorization: faff67ef90c9f2181a58f4bd983fe4dbfd38e8dc32035aad2396ea2ef98b21c3" \
  --form 'parameters={"boardname":"ESP32DevKitC"};type=application/json' \
  --form "file=@bin/ESP32DevKitC.bin;type=application/octet-stream" \
  http://localhost:8083/api/file

6ec0d4238b7164784a62f4b163c712c480b45b2a931ed6c2c6b00e4c66890ca1

curl -v --request POST \
  --header "Authorization: faff67ef90c9f2181a58f4bd983fe4dbfd38e8dc32035aad2396ea2ef98b21c3" \
  --form 'parameters={"boardname":"ArduinoMega2560"};type=application/json' \
  --form "file=@bin/ArduinoMega2560.bin;type=application/octet-stream" \
  http://localhost:8083/api/file

eb3920b037e505b19c9a0ce0d8f28ae56f5ed28d9f70830ed22e11fd07d01c82

curl -v --request POST \
  --header "Authorization: d6b804c0169fdbc0952dc8ef54a2a147d059438e70ec03eee05762913801fd9d" \
  --form 'parameters={"boardname":"ESP32DevKitC"};type=application/json' \
  --form "file=@bin/ESP32DevKitC.bin;type=application/octet-stream" \
  http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file