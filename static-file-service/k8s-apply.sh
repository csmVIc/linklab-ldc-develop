#!/bin/bash

NAMESPACE=linklab

kubectl create configmap static-file-service-config --from-file=config/nginx.conf -n $NAMESPACE

kubectl apply -f data-nfs-pv.yaml

kubectl apply -f k8s.yaml