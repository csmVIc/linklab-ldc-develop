#!/bin/bash

NAMESPACE=linklab

kubectl delete -f k8s.yaml

# kubectl delete configmap jlink-client-config -n $NAMESPACE
