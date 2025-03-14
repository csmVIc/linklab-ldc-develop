#!/bin/bash

PROJECT_NAME="file-cache"

URL="registry.cn-hangzhou.aliyuncs.com"
NAMESPACE="linklab"
PROJECT="device-control-v2-${PROJECT_NAME}"
VERSION="v2.0-beta"
ARCH=""

if [ $(dpkg --print-architecture) == "amd64" ]
then 
  ARCH="amd64"
elif [ $(dpkg --print-architecture) == "armhf" ] 
then 
  ARCH="armhf"
else
  echo "Not support this architecture $(dpkg --print-architecture)"
  exit -1
fi

echo ${ARCH}

cd ../base-library
zip -r base-library.zip .
mv base-library.zip ../${PROJECT_NAME}
cd ../${PROJECT_NAME}

docker rmi ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}
docker build -t ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-${ARCH} --file=Dockerfile .

rm base-library.zip

docker push ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-${ARCH}