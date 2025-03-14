#!/bin/bash

PROJECT_NAME="user-service"

URL="crpi-ldgaqlsrparac7fl.cn-hangzhou.personal.cr.aliyuncs.com"
NAMESPACE="linklab-ym/csmvic"
# PROJECT="device-control-v2-${PROJECT_NAME}"
# VERSION="v2.1-beta" init版本
VERSION="v5.5" # 用户镜像
ARCH=$(dpkg --print-architecture)

echo ${ARCH}

cd ../base-library
zip -r base-library.zip .
mv base-library.zip ../${PROJECT_NAME}
cd ../${PROJECT_NAME}

docker rmi ${URL}/${NAMESPACE}:${VERSION}
docker build -t ${URL}/${NAMESPACE}:${VERSION}-${ARCH} --file=Dockerfile .

rm base-library.zip

docker push ${URL}/${NAMESPACE}:${VERSION}-${ARCH}

