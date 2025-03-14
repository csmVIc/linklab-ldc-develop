#!/bin/bash

NAMESPACE=linklab

kubectl delete -f k8s.yaml

kubectl delete configmap edge-client-config -n $NAMESPACE

kubectl delete configmap edge-client-yaml -n $NAMESPACE