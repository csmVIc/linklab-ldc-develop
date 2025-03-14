#!/bin/bash

NEEDINGRESS=false

NAMESPACE=linklab

kubectl delete -f k8s.yaml

kubectl delete configmap compilev2-gateway-config -n $NAMESPACE

if $NEEDINGRESS; then
  echo "ingress delete"
  kubectl delete -f ingress.yaml
fi
