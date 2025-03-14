#!/bin/bash

IS_ALIYUN_K8S=false

kubectl create configmap influxdb-server-config --from-file=conf/ -n linklab

kubectl create configmap influxdb-server-env --from-env-file=env/influxdb.env -n linklab

if [ $IS_ALIYUN_K8S == false ]
then
  kubectl apply -f pv.yaml
  kubectl apply -f k8s.yaml
else
  kubectl apply -f csi-pv.yaml
  kubectl apply -f aliyun-k8s.yaml
fi