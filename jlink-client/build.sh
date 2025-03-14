#!/bin/bash

PROJECT_NAME="jlink-client"

URL="registry.cn-hangzhou.aliyuncs.com"
NAMESPACE="linklab"
PROJECT="device-control-v2-${PROJECT_NAME}"
VERSION="v1.0"
ARCH="$(dpkg --print-architecture)"
echo ${ARCH}

cd ../base-library
zip -r base-library.zip .
mv base-library.zip ../${PROJECT_NAME}
cd ../${PROJECT_NAME}

if [ "$(dpkg --print-architecture)" == "amd64" ]
then
  docker pull debian:buster-slim --platform linux/amd64
  docker pull golang:1.15.3-buster --platform linux/amd64
elif [ "$(dpkg --print-architecture)" == "armhf" ]
then
  docker pull debian:buster-slim --platform linux/arm/v7
  docker pull golang:1.15.3-buster --platform linux/arm/v7
else
  echo "Not support this architecture $(dpkg --print-architecture)"
  exit -1
fi 

sudo rm -f tmp/*

docker rmi ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}
docker build -t ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-${ARCH} --file=Dockerfile .

rm base-library.zip

docker push ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-${ARCH}