#!/bin/bash

DIFF=21
RANDOM=$$
R=$(($((RANDOM%DIFF))+20))

kubectl get deploy nginx-deployment -n default -o yaml > update.yaml
cat update.yaml | awk '/terminationGracePeriodSeconds/' > terminationGracePeriodSecond.txt
sed -i 's/terminationGracePeriodSeconds://g' terminationGracePeriodSecond.txt
SECONDS=$(awk '{ print $0}' terminationGracePeriodSecond.txt)

KUBE_EDITOR="sed -i 's/terminationGracePeriodSeconds: $SECONDS/terminationGracePeriodSeconds: $R/g'" kubectl edit deploy nginx-deployment -n default
