

# /api/device/listuser

# curl -v \
#   --header "Authorization: 768af315f8fb6bcb9b8951f7754750487a396ef8520739c2206c91ced1f3ac1a" \
#   http://10.214.149.214:31958/api/board/list


# curl -v \
#   --header "Authorization: 768af315f8fb6bcb9b8951f7754750487a396ef8520739c2206c91ced1f3ac1a" \
#   http://10.214.149.214:31958/api/device/list?boardname="all"

# curl -v \
#   --header "Authorization: 5b9bf255425362a6bd5bd5a5488f3034cca8921def8ec5ce6d72fb501998e89d" \
#   http://10.214.149.214:31958/api/device/listuser

# curl -v \
#   --header "Authorization: 5b9bf255425362a6bd5bd5a5488f3034cca8921def8ec5ce6d72fb501998e89d" \
#   http://192.168.88.20:31958/api/device/listuser

curl -v \
  --header "Authorization: deaa2efc37923534cd52d70e782ff4d7a84583fb1c0f6784d4b88007862a92cf" \
  http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/listuserdevice
