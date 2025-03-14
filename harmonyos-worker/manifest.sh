#!/bin/bash

PROJECT_NAME="harmonyos-worker"

URL="registry.cn-hangzhou.aliyuncs.com"
NAMESPACE="linklab"
PROJECT="${PROJECT_NAME}"
VERSION="v2.1-beta"
rm -r ~/.docker/manifests
docker manifest create --amend ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-amd64
docker manifest annotate ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-amd64 \
      --os linux --arch amd64
docker manifest push ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}
rm -r ~/.docker/manifests
docker manifest inspect ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}