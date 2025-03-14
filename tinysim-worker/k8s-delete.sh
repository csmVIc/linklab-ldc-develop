#!/bin/bash

IS_ALIYUN_K8S=false
NAMESPACE=linklab

kubectl delete -f k8s.yaml

if [ $IS_ALIYUN_K8S == false ]
then
  kubectl delete -f init-nfs-pv.yaml
  kubectl get pvc -n $NAMESPACE | grep -a "compilev2-tinysim-worker" | awk '{print $1}' | xargs kubectl delete pvc -n $NAMESPACE
  kubectl delete -f writespace-host-pv.0.yaml
  kubectl delete -f writespace-host-pv.1.yaml
else
  # TODO
  echo "TODO"
fi

kubectl delete configmap tinysim-worker-config -n $NAMESPACE