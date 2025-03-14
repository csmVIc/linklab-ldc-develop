#!/bin/bash

IS_ALIYUN_K8S=true
NAMESPACE=linklab

kubectl delete -f k8s.yaml

if [ $IS_ALIYUN_K8S == false ]
then
  echo "TODO"
else
  kubectl get pvc -n $NAMESPACE | grep -a "virtual-device-client" | awk '{print $1}' | xargs kubectl delete pvc -n $NAMESPACE
  kubectl delete -f csi-pv-0.yaml
  kubectl delete -f csi-pv-1.yaml
  kubectl delete -f csi-pv-2.yaml
  kubectl delete -f csi-pv-3.yaml
  kubectl delete -f csi-pv-4.yaml
  kubectl delete -f csi-pv-5.yaml
  kubectl delete -f csi-pv-6.yaml
  kubectl delete -f csi-pv-7.yaml
  kubectl delete -f csi-pv-8.yaml
  kubectl delete -f csi-pv-9.yaml
fi

kubectl delete configmap virtual-device-client-config -n $NAMESPACE