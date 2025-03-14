#!/bin/bash

IS_ALIYUN_K8S=true
NAMESPACE=linklab

if [ $IS_ALIYUN_K8S == false ]
then
  echo "TODO"
else
  kubectl apply -f csi-pv-0.yaml
  kubectl apply -f csi-pv-1.yaml
  kubectl apply -f csi-pv-2.yaml
  kubectl apply -f csi-pv-3.yaml
  kubectl apply -f csi-pv-4.yaml
  kubectl apply -f csi-pv-5.yaml
  kubectl apply -f csi-pv-6.yaml
  kubectl apply -f csi-pv-7.yaml
  kubectl apply -f csi-pv-8.yaml
  kubectl apply -f csi-pv-9.yaml
fi

kubectl create configmap virtual-device-client-config --from-file=config/k8s-config.json -n $NAMESPACE

kubectl apply -f k8s.yaml