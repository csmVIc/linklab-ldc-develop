#!/bin/bash
#  创建和推送 Docker 镜像清单
PROJECT_NAME="edge-client"

URL="registry.cn-hangzhou.aliyuncs.com"
NAMESPACE="linklab"
PROJECT="device-control-v2-${PROJECT_NAME}"
VERSION="v2.0-beta"
# 清除旧的manifest清单
rm -r ~/.docker/manifests
# 创建docker manifest
docker manifest create --amend ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-amd64 \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-armhf \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-arm64
# 添加架构注解 - 用于多架构镜像 - armhf
docker manifest annotate ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-armhf \
      --os linux --arch arm --variant v7
# 添加架构注解 - 用于多架构镜像 - amd64
docker manifest annotate ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-amd64 \
      --os linux --arch amd64
# 添加架构注解 - 用于多架构镜像 - arm64
docker manifest annotate ${URL}/${NAMESPACE}/${PROJECT}:${VERSION} \
      ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}-arm64 \
      --os linux --arch arm64 --variant v8
# 推送docker manifest
docker manifest push ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}
# 删除本地manifest
rm -r ~/.docker/manifests
# 查看manifest
docker manifest inspect ${URL}/${NAMESPACE}/${PROJECT}:${VERSION}