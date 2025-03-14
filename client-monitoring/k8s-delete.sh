#!/bin/bash

kubectl delete -f k8s.yaml

kubectl delete configmap client-monitoring-config -n linklab