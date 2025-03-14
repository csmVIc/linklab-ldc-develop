#!/bin/bash

kubectl create configmap resource-monitoring-config --from-file=config/config.json -n linklab

kubectl apply -f k8s.yaml