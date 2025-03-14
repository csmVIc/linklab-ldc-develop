#!/bin/bash

PROJECT_NAME="alios-edu-worker"

URL="registry.cn-hangzhou.aliyuncs.com"
NAMESPACE="linklab"
PROJECT="${PROJECT_NAME}"
VERSION="v2.0-beta"
ARCH=""

if [ $(uname -m) == "x86_64" ]
then 
  ARCH="amd64"
elif [ $(uname -m) == "armv7l" ] 
then 
  ARCH="arm"
else
  echo "Not support this architecture $(uname -m)"
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