#!/bin/bash

IS_ALIYUN_K8S=false

if [ $IS_ALIYUN_K8S == false ]
then
  kubectl delete -f k8s.yaml
  kubectl delete -f pv.yaml
else
  kubectl delete -f aliyun-k8s.yaml
  kubectl delete -f csi-pv.yaml
fi

kubectl delete configmap influxdb-server-config -n linklab
kubectl delete configmap influxdb-server-env -n linklab