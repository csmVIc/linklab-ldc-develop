#!/bin/bash

kubectl delete -f k8s.yaml

kubectl delete configmap device-service-config -n linklab