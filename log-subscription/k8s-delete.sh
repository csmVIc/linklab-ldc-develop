#!/bin/bash

kubectl delete -f k8s.yaml

kubectl delete configmap log-subscription-config -n linklab