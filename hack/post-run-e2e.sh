#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

kubectl delete -n kube-system daemonset huawei-cloud-controller-manager
kubectl delete serviceaccount cloud-controller-manager
kubectl delete clusterrole ccm-secret-role
kubectl delete clusterrolebinding ccm-secret-binding
kubectl delete clusterrolebinding extension-apiserver-auth-reader-binding
