#!/bin/bash

NUMBER=1

# LDC
# kubectl delete hpa file-cache --namespace=linklab
# kubectl delete hpa user-service --namespace=linklab
# kubectl delete hpa login-authentication --namespace=linklab
# kubectl delete hpa device-service --namespace=linklab
# kubectl delete hpa decision-maker --namespace=linklab
# kubectl delete hpa client-monitoring --namespace=linklab
# kubectl delete hpa log-subscription --namespace=linklab

# Compile
# kubectl delete hpa compilev2-gateway --namespace=linklab
# kubectl delete hpa compilev2-mongo-gateway --namespace=linklab
# kubectl delete hpa compilev2-tinysim-worker --namespace=linklab
# kubectl delete hpa compilev2-arduino-worker --namespace=linklab

# LDC
kubectl scale deployment file-cache --replicas=${NUMBER} --namespace=linklab
kubectl scale deployment user-service --replicas=${NUMBER} --namespace=linklab
kubectl scale deployment login-authentication --replicas=${NUMBER} --namespace=linklab
kubectl scale deployment device-service --replicas=${NUMBER} --namespace=linklab
kubectl scale deployment decision-maker --replicas=${NUMBER} --namespace=linklab
kubectl scale deployment client-monitoring --replicas=${NUMBER} --namespace=linklab
kubectl scale deployment log-subscription --replicas=${NUMBER} --namespace=linklab

# Compile
kubectl scale deployment compilev2-gateway --replicas=${NUMBER} --namespace=linklab
kubectl scale deployment compilev2-mongo-gateway --replicas=${NUMBER} --namespace=linklab
# kubectl scale statefulset compilev2-tinysim-worker --replicas=${NUMBER} --namespace=linklab
# kubectl scale statefulset compilev2-arduino-worker --replicas=${NUMBER} --namespace=linklab
kubectl scale statefulset compilev2-alios-worker --replicas=${NUMBER} --namespace=linklab