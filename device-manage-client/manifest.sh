#!/bin/bash

PROJECT_NAME="device-manage-client"

URL="registry.cn-hangzhou.aliyuncs.com"
NAMESPACE="linklab"
PROJECT="device-control-v2-${PROJECT_NAME}"
VERSION="v1.0"

rm -r ~/.docker/manifests

docker manifest create --amend ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-armhf \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-amd64
      # ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-arm64

docker manifest annotate ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-armhf \
      --os linux --arch arm --variant v7

# docker manifest annotate ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
#       ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-arm64 \
#       --os linux --arch arm64

docker manifest annotate ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-amd64 \
      --os linux --arch amd64

docker manifest push ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}
rm -r ~/.docker/manifests
docker manifest inspect ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}