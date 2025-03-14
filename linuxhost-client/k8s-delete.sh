#!/bin/bash

kubectl delete -f k8s.yaml

kubectl delete configmap linuxhost-client-config -n linklab

kubectl delete pvc writespace-volume-linuxhost-client-0 -n linklab

kubectl delete pvc writespace-volume-linuxhost-client-1 -n linklab

kubectl delete -f pv-0.yaml

kubectl delete -f pv-1.yaml