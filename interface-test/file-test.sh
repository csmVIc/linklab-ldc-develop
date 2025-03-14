#!/bin/bash

# while true; do
  # curl -v --request POST \
  # --header "Authorization: 7ccae70a7bb244dcd9ad8d13eb56d3e5a1fd7f2e27c6848c09b3d9b2fa19892e" \
  # --form 'parameters={"boardname":"VirtualDevice"};type=application/json' \
  # --form "file=@bin/VirtualDevice.bin;type=application/octet-stream" \
  # http://10.214.149.214:32710/api/file
# done

# while true; do
  # curl -v --request POST \
  # --header "Authorization: de0f8c718c1be17e223df00485c61f9d7fe19d5496b912869c63556871f8d59a" \
  # --form 'parameters={"boardname":"ArduinoMega2560"};type=application/json' \
  # --form "file=@bin/ArduinoMega2560.bin;type=application/octet-stream" \
  # http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file
# done

#  curl -v --request POST \
#   --header "Authorization: beb51d947b0358fe9216a410880cdb20299a12788dbaafa425712f026f762a4c" \
#   --form 'parameters={"boardname":"TinySim"};type=application/json' \
#   --form "file=@bin/TinySim.zip;type=application/octet-stream" \
#   http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file

  # curl -v --request POST \
  #   --header "Authorization: 423c0889878290e7fc0369089a7bf0284bc5871b1055b3cb803e4d974a805b13" \
  #   --form 'parameters={"boardname":"Haas100"};type=application/json' \
  #   --form "file=@bin/Haas100.bin;type=application/octet-stream" \
  #   http://10.214.149.214:32710/api/file

  # curl -v --request POST \
  #   --header "Authorization: 423c0889878290e7fc0369089a7bf0284bc5871b1055b3cb803e4d974a805b13" \
  #   --form 'parameters={"boardname":"Haas100"};type=application/json' \
  #   --form "file=@bin/Haas100-MQTT.bin;type=application/octet-stream" \
  #   http://10.214.149.214:32710/api/file
  
  #   curl -v --request POST \
  #   --header "Authorization: f5b93c65a79215d7b7a8964438b7716f8f01804f91ea99a3e05ba390912ead16" \
  #   --form 'parameters={"boardname":"Haas100"};type=application/json' \
  #   --form "file=@bin/Haas100.bin;type=application/octet-stream" \
  #   http://192.168.88.20:32710/api/file

  # curl -v --request POST \
  #   --header "Authorization: f5b93c65a79215d7b7a8964438b7716f8f01804f91ea99a3e05ba390912ead16" \
  #   --form 'parameters={"boardname":"Haas100"};type=application/json' \
  #   --form "file=@bin/Haas100-MQTT.bin;type=application/octet-stream" \
  #   http://192.168.88.20:32710/api/file
  
  # curl -v --request POST \
  #   --header "Authorization: fb70630a5cda4ec40dd62cd669dce0ccb22115ecf9d95d1ec2f0a98feb6b2bbc" \
  #   --form 'parameters={"boardname":"Haas100"};type=application/json' \
  #   --form "file=@bin/Haas100.bin;type=application/octet-stream" \
  #   http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file
  

  # curl -v --request POST \
  #   --header "Authorization: fb70630a5cda4ec40dd62cd669dce0ccb22115ecf9d95d1ec2f0a98feb6b2bbc" \
  #   --form 'parameters={"boardname":"Haas100"};type=application/json' \
  #   --form "file=@bin/Haas100-MQTT.bin;type=application/octet-stream" \
  #   http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file
  
  # curl -v --request POST \
  #   --header "Authorization: 7185426bf458d9b7e9daf785176140f7eb6c9bfb57cbfb5227728921149fd9b3" \
  #   --form 'parameters={"boardname":"Haas100"};type=application/json' \
  #   --form "file=@bin/Haas100-Python.zip;type=application/octet-stream" \
  #   http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file
  
  # curl -v --request POST \
  #   --header "Authorization: fb70630a5cda4ec40dd62cd669dce0ccb22115ecf9d95d1ec2f0a98feb6b2bbc" \
  #   --form 'parameters={"boardname":"ESP32DevKitCArduino"};type=application/json' \
  #   --form "file=@bin/ESP32DevKitCArduino-MQTT.bin;type=application/octet-stream" \
  #   http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file
  
  curl -v --request POST \
    --header "Authorization: fb70630a5cda4ec40dd62cd669dce0ccb22115ecf9d95d1ec2f0a98feb6b2bbc" \
    --form 'parameters={"boardname":"Python3Exec"};type=application/json' \
    --form "file=@bin/Python3Exec-MQTT3.zip;type=application/octet-stream" \
    http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file
  
  #  curl -v --request POST \
  #   --header "Authorization: 9ca4f83e28648468ed335953905d8877bb63f7f22c2cd20e52e3799912c04707" \
  #   --form 'parameters={"boardname":"Python3Exec"};type=application/json' \
  #   --form "file=@bin/Python3Exec-Hello.zip;type=application/octet-stream" \
  #   http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file
  
  # curl -v --request POST \
  #   --header "Authorization: 9ca4f83e28648468ed335953905d8877bb63f7f22c2cd20e52e3799912c04707" \
  #   --form 'parameters={"boardname":"TinySim"};type=application/json' \
  #   --form "file=@bin/TinySim.zip;type=application/octet-stream" \
  #   http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file
  
  # curl -v --request POST \
  #   --header "Authorization: 9ca4f83e28648468ed335953905d8877bb63f7f22c2cd20e52e3799912c04707" \
  #   --form 'parameters={"boardname":"TinySim"};type=application/json' \
  #   --form "file=@bin/TinySim.zip;type=application/octet-stream" \
  #   http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file

# while true; do
  # curl -v --request POST \
  # --header "Authorization: c2f7193b949290ebbc519d3edd7c3b0735595c26cc8256228e58d1046ae6570e" \
  # --form 'parameters={"boardname":"ESP32DevKitC"};type=application/json' \
  # --form "file=@bin/ESP32DevKitC.bin;type=application/octet-stream" \
  # http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file
# done

# while true; do
  # curl -v --request POST \
  # --header "Authorization: c2f7193b949290ebbc519d3edd7c3b0735595c26cc8256228e58d1046ae6570e" \
  # --form 'parameters={"boardname":"TelosB"};type=application/json' \
  # --form "file=@bin/TelosB.bin;type=application/octet-stream" \
  # http://kubernetes.tinylink.cn/linklab/device-control-v2/file-cache/api/file
# done
