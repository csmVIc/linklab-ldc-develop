#!/bin/bash

PROJECT_NAME="edge-build"

URL="registry.cn-hangzhou.aliyuncs.com"
NAMESPACE="linklab"
PROJECT="${PROJECT_NAME}"
VERSION="v2.0-beta"

rm -r ~/.docker/manifests

# 创建一个新的manifest --amend 标志表示如果manifest已存在则更新它,而不是报错
docker manifest create --amend ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-armhf \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-arm64
     
docker manifest annotate ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-armhf \
      --os linux --arch arm --variant v7

docker manifest annotate ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-arm64 \
      --os linux --arch arm64 --variant v8

docker manifest push ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}
rm -r ~/.docker/manifests
docker manifest inspect ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}