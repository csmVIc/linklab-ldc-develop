#!/bin/bash

IS_ALIYUN_K8S=false
NAMESPACE="linklab"

helm list -n $NAMESPACE | grep -a "redis-server" | awk '{print $1}' | xargs helm uninstall -n $NAMESPACE

kubectl get pvc -n $NAMESPACE | grep -a "redis-server" | awk '{print $1}' | xargs kubectl delete pvc -n $NAMESPACE

if [ $IS_ALIYUN_K8S == false ]
then
  kubectl delete -f pv-0.yaml
  kubectl delete -f pv-1.yaml
  kubectl delete -f pv-2.yaml
else
  kubectl delete -f csi-pv-0.yaml
  kubectl delete -f csi-pv-1.yaml
  kubectl delete -f csi-pv-2.yaml
fi 