#!/bin/bash

ARCH=""

if [ $(uname -m) == "x86_64" ]
then 
  docker load -i image/amd64/metrics-server.tar
  docker load -i image/amd64/bitnami-mongodb.tar
  docker load -i image/amd64/bitnami-redis-sentinel.tar
  docker load -i image/amd64/bitnami-redis.tar
  docker load -i image/amd64/nginx-ingress-controller.tar
  ARCH="amd64"
elif [ $(uname -m) == "armv7l" ] 
then 
  ARCH="arm"
else
  echo "Not support this architecture $(uname -m)"
  exit -1
fi

echo ${ARCH}