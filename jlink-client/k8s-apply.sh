#!/bin/bash

NAMESPACE=linklab

# kubectl create configmap jlink-client-config --from-file=config/config.json -n $NAMESPACE

kubectl apply -f k8s.yaml