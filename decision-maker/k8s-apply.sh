#!/bin/bash

NAMESPACE=linklab

kubectl create configmap decision-maker-config --from-file=config/k8s-config.json -n $NAMESPACE

kubectl apply -f k8s.yaml