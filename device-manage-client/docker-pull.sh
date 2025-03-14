#!/bin/bash

docker pull registry.cn-hangzhou.aliyuncs.com/linklab/device-control-v2-device-manage-client:v1.0
docker images -a | grep "none" | awk '{print $3}' | xargs docker rmi