#!/bin/bash

PROJECT_NAME="edge-allocater"

URL="crpi-ldgaqlsrparac7fl.cn-hangzhou.personal.cr.aliyuncs.com"
NAMESPACE="linklab-ym/csmvic"
# VERSION="v2.0-beta"
# VERSION="v2.1"
VERSION="v14.1"
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
