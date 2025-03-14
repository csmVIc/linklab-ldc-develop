#!/bin/bash

IS_ALIYUN_K8S=false
NAMESPACE="linklab"
VERSION="12.2.2"

if [ $IS_ALIYUN_K8S == false ]
then
  kubectl apply -f pv-0.yaml
  kubectl apply -f pv-1.yaml
  kubectl apply -f pv-2.yaml
else
  kubectl apply -f csi-pv-0.yaml
  kubectl apply -f csi-pv-1.yaml
  kubectl apply -f csi-pv-2.yaml
fi

if [ $IS_ALIYUN_K8S == false ]
then
  helm install redis-server bitnami/redis -f conf/values.yaml --namespace $NAMESPACE --version $VERSION
else
  helm install redis-server bitnami/redis -f conf/aliyun-values.yaml --namespace $NAMESPACE --version $VERSION
fi



