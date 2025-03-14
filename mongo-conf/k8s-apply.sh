#!/bin/bash

IS_ALIYUN_K8S=false
NAMESPACE="linklab"
VERSION="10.4.0"

if [ $IS_ALIYUN_K8S == false ]
then
  kubectl apply -f pv-0.yaml
  # kubectl apply -f pv-1.yaml
  # kubectl apply -f pv-2.yaml
else
  kubectl apply -f csi-pv-0.yaml
  kubectl apply -f csi-pv-1.yaml
  kubectl apply -f csi-pv-2.yaml
fi

kubectl create configmap mongo-server-init-shell --from-file=scripts/ -n $NAMESPACE

if [ $IS_ALIYUN_K8S == false ]
then
  helm install mongo-server -f conf/values.yaml bitnami/mongodb --namespace $NAMESPACE --version $VERSION
else
  helm install mongo-server -f conf/aliyun-values.yaml bitnami/mongodb --namespace $NAMESPACE --version $VERSION
fi

if [ $IS_ALIYUN_K8S == false ]
then
  kubectl apply -f nodeport-0.yaml
  # kubectl apply -f nodeport-1.yaml
  # kubectl apply -f nodeport-2.yaml
fi

