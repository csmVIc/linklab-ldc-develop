#!/bin/bash

NAMESPACE="linklab"

helm list -n $NAMESPACE | grep -a "emqx-server" | awk '{print $1}' | xargs helm uninstall -n $NAMESPACE

# helm list -n $NAMESPACE | grep -a "mqtt-for-user-server" | awk '{print $1}' | xargs helm uninstall -n $NAMESPACE