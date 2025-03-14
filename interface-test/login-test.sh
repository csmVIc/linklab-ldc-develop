#!/bin/bash

NUM=0

while [ $NUM -le 0 ]; do

  USERNAME="UserTest-$NUM"

  # curl --header "Content-Type: application/json" \
  #   --request POST \
  #   --data "{\"id\":\"${USERNAME}\",\"password\":\"6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918\"}" \
  #   http://kubernetes.tinylink.cn/linklab/device-control-v2/login-authentication/user/login

  # curl --header "Content-Type: application/json" \
  #   --request POST \
  #   --data "{\"id\":\"${USERNAME}\",\"password\":\"6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918\"}" \
  #   http://10.214.149.214:31285/user/login
  
  # curl --header "Content-Type: application/json" \
  #   --request POST \
  #   --data "{\"id\":\"Nancie-09\",\"password\":\"6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918\"}" \
  #   http://10.214.149.214:31285/user/login
  
  curl --header "Content-Type: application/json" \
    --request POST \
    --data "{\"id\":\"yangg\",\"password\":\"6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918\"}" \
    http://kubernetes.tinylink.cn/linklab/device-control-v2/login-authentication/user/login

  # curl --header "Content-Type: application/json" \
  #   --request POST \
  #   --data "{\"id\":\"${USERNAME}\",\"password\":\"6b51d431df5d7f141cbececcf79edf3dd861c3b4069f0b11661a3eefacbba918\"}" \
  #   http://192.168.88.20:31285/user/login

  NUM=$(( $NUM + 1 ))
done