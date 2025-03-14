#!/bin/bash

NAMESPACE=linklab

kubectl create configmap linuxhost-client-config --from-file=config/k8s-config.json -n $NAMESPACE

kubectl apply -f pv-0.yaml

kubectl apply -f pv-1.yaml

kubectl apply -f k8s.yaml