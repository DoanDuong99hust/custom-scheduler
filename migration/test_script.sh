#!/bin/bash

# xoa container tren cloud
# kubectl delete service decode
# kubectl delete deployment nginx-deployment

echo "finish delete deployment"
# deploy lai duoi edge
kubectl create -f /home/node5/test_deployment.yaml
echo "finish create deployment"

