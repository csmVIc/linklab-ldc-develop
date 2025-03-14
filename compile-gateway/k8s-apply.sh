#!/bin/bash

NEEDINGRESS=false

NAMESPACE=linklab

kubectl create configmap compilev2-gateway-config --from-file=config/config.json -n $NAMESPACE

kubectl apply -f k8s.yaml

if $NEEDINGRESS; then
  echo "ingress init"
  kubectl apply -f ingress.yaml
fi
