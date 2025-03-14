#!/bin/bash

kubectl delete -f k8s.yaml

kubectl delete configmap edge-allocater-config -n linklab