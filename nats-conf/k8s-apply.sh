#!/bin/bash

NEEDSTAN=false
IS_ALIYUN_K8S=false

echo "nats init"
kubectl create configmap nats-server-config --from-file=conf/k8s-nats.conf -n linklab
kubectl apply -f k8s-nats.yaml

if $NEEDSTAN; then
  echo "stan init"
  kubectl create configmap stan-server-config --from-file=conf/k8s-stan.conf -n linklab

  if [ $IS_ALIYUN_K8S == false ]
  then
    kubectl apply -f pv-0.yaml
    kubectl apply -f pv-1.yaml
    kubectl apply -f pv-2.yaml
  else
    kubectl apply -f csi-pv-0.yaml
    kubectl apply -f csi-pv-1.yaml
    kubectl apply -f csi-pv-2.yaml
    kubectl apply -f csi-pv-3.yaml
    kubectl apply -f csi-pv-4.yaml
  fi

  kubectl apply -f k8s-stan.yaml
fi