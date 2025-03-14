#!/bin/bash

NEEDINGRESS=false

NAMESPACE=linklab

kubectl create configmap user-service-config --from-file=config/k8s-config.json -n $NAMESPACE

kubectl apply -f k8s.yaml

if $NEEDINGRESS; then
  echo "ingress init"
  kubectl apply -f ingress.yaml
fi
