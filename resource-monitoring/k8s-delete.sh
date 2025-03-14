#!/bin/bash

kubectl delete -f k8s.yaml

kubectl delete configmap resource-monitoring-config -n linklab