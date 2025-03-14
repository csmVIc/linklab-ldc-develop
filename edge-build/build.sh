#!/bin/bash

PROJECT_NAME="edge-build"

URL="registry.cn-hangzhou.aliyuncs.com"
NAMESPACE="linklab"
PROJECT="${PROJECT_NAME}"
VERSION="v2.0-beta"
ARCH=$(dpkg --print-architecture)

echo ${ARCH}

cd ../base-library
zip -r base-library.zip .
mv base-library.zip ../${PROJECT_NAME}
cd ../${PROJECT_NAME}

docker rmi ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}
docker build -t ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-${ARCH} --file=Dockerfile .

rm base-library.zip

docker push ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-${ARCH}