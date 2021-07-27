#!/bin/bash

# xoa container tren cloud
# kubectl delete service decode
kubectl delete deployment nginx1-deployment

echo "finish delete deployment 1"
# deploy lai duoi edge
kubectl create -f /home/node5/test_deployment_1.yaml
echo "finish create deployment 1"

