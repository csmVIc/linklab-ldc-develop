#!/bin/bash

# curl -v --request POST \
#  --header "Authorization: 94acd418c37c3f47ba5c731a750f0b0b9ec4d1ff0e87cf1c5afc7ca9790a9a7d" \
#  --data '{"cmd":"netmgr -t wifi -c Xiaomi_C690_HA eagle402\r\n", "deviceid":"/dev/Haas100-0", "clientid":"ClientTest"}' \
#  http://10.214.149.214:31958/api/device/cmd

# curl -v --request POST \
#  --header "Authorization: 94acd418c37c3f47ba5c731a750f0b0b9ec4d1ff0e87cf1c5afc7ca9790a9a7d" \
#  --data '{"cmd":"reboot\r\n", "deviceid":"/dev/Haas100-0", "clientid":"ClientTest"}' \
#  http://10.214.149.214:31958/api/device/cmd

# curl -v --request POST \
#  --header "Authorization: 5b9bf255425362a6bd5bd5a5488f3034cca8921def8ec5ce6d72fb501998e89d" \
#  --data '{"cmd":"netmgr -t wifi -c linklab-wifi-1 eagle402\r\n", "deviceid":"/dev/Haas100-0", "clientid":"ClientTest"}' \
#  http://192.168.88.20:31958/api/device/cmd

# curl -v --request POST \
#  --header "Authorization: 5b9bf255425362a6bd5bd5a5488f3034cca8921def8ec5ce6d72fb501998e89d" \
#  --data '{"cmd":"reboot\r\n", "deviceid":"/dev/Haas100-0", "clientid":"ClientTest"}' \
#  http://192.168.88.20:31958/api/device/cmd

curl -v --request POST \
 --header "Authorization: 7185426bf458d9b7e9daf785176140f7eb6c9bfb57cbfb5227728921149fd9b3" \
 --data '{"cmd":"reboot\r\n", "deviceid":"/dev/Haas100-2", "clientid":"ClientTest-25"}' \
 http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/cmd