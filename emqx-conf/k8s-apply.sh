#!/bin/bash

NAMESPACE="linklab"
VERSION="4.2.4"

helm install emqx-server emqx/emqx -f config/values.yaml --namespace $NAMESPACE --version $VERSION

# helm install mqtt-for-user-server emqx/emqx -f config/values-for-user.yaml --namespace $NAMESPACE --version $VERSION