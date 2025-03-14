#!/bin/bash

kubectl delete -f k8s.yaml

kubectl delete configmap decision-maker-config -n linklab