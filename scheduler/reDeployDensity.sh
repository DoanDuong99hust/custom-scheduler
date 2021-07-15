#!/bin/bash
path_decode_deployment="/home/node5/camera_deployment/deployment_script/decode.yaml";
path_density_deployment="/home/node5/camera_deployment/deployment_script/density.yaml";
action1="delete";
action2="apply";
node=`kubectl get pod -n default -o wide | awk 'BEGIN {getline;getline;getline;getline;name=$1;node=$7;print node}'`
if [ $node == "server" ]
then
        kubectl $action1 -f $path_density_deployment
        kubectl $action2 -f $path_density_deployment
fi

