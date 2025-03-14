#!/bin/bash

while true; do
  curl -v \
  --header "Authorization: 9ca4f83e28648468ed335953905d8877bb63f7f22c2cd20e52e3799912c04707" \
  http://kubernetes.tinylink.cn/linklab/device-control-v2/user-service/api/device/list?boardname="all"
done


