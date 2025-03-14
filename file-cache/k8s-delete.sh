#!/bin/bash

NEEDINGRESS=false

kubectl delete -f k8s.yaml

kubectl delete configmap file-cache-config -n linklab

if $NEEDINGRESS; then
  echo "ingress delete"
  kubectl delete -f ingress.yaml
fi