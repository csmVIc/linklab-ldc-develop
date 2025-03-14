#!/bin/bash

NEEDSTAN=true
IS_ALIYUN_K8S=false

kubectl delete -f k8s-nats.yaml
kubectl delete configmap nats-server-config -n linklab

if $NEEDSTAN; then
  kubectl delete -f k8s-stan.yaml

  kubectl get pvc -n linklab | grep -a "stan" | awk '{print $1}' | xargs kubectl delete pvc -n linklab

  if [ $IS_ALIYUN_K8S == false ]
  then
    kubectl delete -f pv-0.yaml
    kubectl delete -f pv-1.yaml
    kubectl delete -f pv-2.yaml
  else
    kubectl delete -f csi-pv-0.yaml
    kubectl delete -f csi-pv-1.yaml
    kubectl delete -f csi-pv-2.yaml
    kubectl delete -f csi-pv-3.yaml
    kubectl delete -f csi-pv-4.yaml
  fi

  kubectl delete configmap stan-server-config -n linklab
fi