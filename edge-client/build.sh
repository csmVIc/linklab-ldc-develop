#!/bin/bash

PROJECT_NAME="edge-client"

URL="crpi-ldgaqlsrparac7fl.cn-hangzhou.personal.cr.aliyuncs.com"
NAMESPACE="linklab-ym/csmvic"
# PROJECT="device-control-v2-${PROJECT_NAME}"
# VERSION="v2.0-beta"
# VERSION="v6.15"
# VERSION="v6.21"   //初始版本 有podapply有Debugf日志的版本
# VERSION="v6.22"   //更新linkingress debugf日志
# VERSION="v6.25"   # 更新了linkingress中关于service的创建的内容，包括创建的service是NodePort类型的
# VERSION="v6.28"  # 更新了Port为99，linkpodingress函数返回的port值是pod暴露的端口号，能够返回nodeport了
# VERSION="v6.31" # 更新了podlog中关于websocket超时时间和缓存区大小的处理，设置了debugf日志，解决问题：能够正常输出日志
VERSION="v6.33" # 更新了podapply中pod的创建过程中，会检查env参数，当enableEdgeSocket为true时，会增加args参数，
                # 并且增加了GetPodNamesByNamespace函数，会从redis中根据namespace获取podname。
VERSION="v6.34" # # 修改linkingress中bug,当容器没有端口号时，不创建svc和ingress
VERSION="v6.36" # 	enableEdgeSocket = true
VERSION="v6.37" # 	取消设置enableEdgeSocket = true,取消了6.34关于端口号的规则判断，并增加了打印端口号len的日志
VERSION="v6.38" # 	在6.37的基础上，增加了p.CreateIngress = true 
                # 总结38:设置了CreateIngress = true其他的都是默认值，貌似模型的协同连接和svc ingress的创建是有关系的
VERSION="v6.44" # 	只修改了podapply.go文件，创建edge的时候会先创建cloud
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