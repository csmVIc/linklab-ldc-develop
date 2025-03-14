#!/bin/bash

NAMESPACE=linklab

kubectl create configmap edge-client-config --from-file=config/config.json -n $NAMESPACE

kubectl create configmap edge-client-yaml --from-file=yaml/build.yaml -n $NAMESPACE

kubectl apply -f k8s.yaml