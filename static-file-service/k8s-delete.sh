#!/bin/bash

NAMESPACE=linklab

kubectl delete -f k8s.yaml

kubectl delete -f data-nfs-pv.yaml

kubectl delete configmap static-file-service-config -n $NAMESPACE