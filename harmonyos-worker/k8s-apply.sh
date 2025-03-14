#!/bin/bash

IS_ALIYUN_K8S=false
NAMESPACE=linklab

if [ $IS_ALIYUN_K8S == false ]
then
  kubectl apply -f init-nfs-pv.yaml
  kubectl apply -f writespace-host-pv.0.yaml
  kubectl apply -f writespace-host-pv.1.yaml
else
  # TODO
  echo "TODO"
fi

kubectl create configmap harmonyos-worker-config --from-file=config/config.json -n $NAMESPACE

kubectl apply -f k8s.yaml