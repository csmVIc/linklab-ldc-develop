#!/bin/bash

NAMESPACE=linklab

kubectl delete -f k8s.yaml

# kubectl delete configmap device-manage-client-config -n $NAMESPACE
