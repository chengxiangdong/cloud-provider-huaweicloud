#!/usr/bin/env bash

# Copyright 2022 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

count=`kubectl get -n kube-system secret cloud-config | grep cloud-config | wc -l`
if [[ "$count" -ne 1 ]]; then
  echo "Please create a secret with the name: cloud-config."
  exit 1
fi

cat << EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:cloud-controller-manager
rules:
  - resources:
      - tokenreviews
    verbs:
      - get
      - list
      - watch
      - create
      - update
      - patch
    apiGroups:
      - authentication.k8s.io
  - resources:
      - configmaps
      - endpoints
      - pods
      - services
      - secrets
      - serviceaccounts
    verbs:
      - get
      - list
      - watch
      - create
      - update
      - patch
    apiGroups:
      - ''
  - resources:
      - nodes
    verbs:
      - get
      - list
      - watch
      - delete
      - patch
      - update
    apiGroups:
      - ''
  - resources:
      - services/status
      - pods/status
    verbs:
      - update
      - patch
    apiGroups:
      - ''
  - resources:
      - nodes/status
    verbs:
      - patch
      - update
    apiGroups:
      - ''
  - resources:
      - events
      - endpoints
    verbs:
      - create
      - patch
      - update
    apiGroups:
      - ''
  - resources:
      - leases
    verbs:
      - get
      - update
      - create
      - delete
    apiGroups:
      - coordination.k8s.io
  - resources:
      - customresourcedefinitions
    verbs:
      - get
      - update
      - create
      - delete
    apiGroups:
      - apiextensions.k8s.io
  - resources:
      - ingresses
    verbs:
      - get
      - list
      - watch
      - update
      - create
      - patch
      - delete
    apiGroups:
      - networking.k8s.io
  - resources:
      - ingresses/status
    verbs:
      - update
      - patch
    apiGroups:
      - networking.k8s.io
  - resources:
      - endpointslices
    verbs:
      - get
      - list
      - watch
    apiGroups:
      - discovery.k8s.io
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cloud-controller-manager
  namespace: kube-system
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:cloud-controller-manager
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:cloud-controller-manager
subjects:
  - kind: ServiceAccount
    name: cloud-controller-manager
    namespace: kube-system
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: huawei-cloud-controller-manager
  namespace: kube-system
  labels:
    k8s-app: huawei-cloud-controller-manager
spec:
  selector:
    matchLabels:
      k8s-app: huawei-cloud-controller-manager
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        k8s-app: huawei-cloud-controller-manager
    spec:
      nodeSelector:
        node-role.kubernetes.io/master: ""
      securityContext:
        runAsUser: 1001
      tolerations:
        - key: node.cloudprovider.kubernetes.io/uninitialized
          value: "true"
          effect: NoSchedule
        - key: node-role.kubernetes.io/master
          effect: NoSchedule
      serviceAccountName: cloud-controller-manager
      containers:
        - name: huawei-cloud-controller-manager
          imagePullPolicy: Never
          image: swr.ap-southeast-1.myhuaweicloud.com/k8scloudcontrollermanager/huawei-cloud-controller-manager:latest
          args:
            - /bin/huawei-cloud-controller-manager
            - --v=5
            - --cloud-config=/etc/config/cloud.conf
            - --cloud-provider=huaweicloud
            - --use-service-account-credentials=true
            - --bind-address=127.0.0.1
          volumeMounts:
            - mountPath: /etc/kubernetes
              name: k8s-certs
              readOnly: true
            - mountPath: /etc/ssl/certs
              name: ca-certs
              readOnly: true
            - mountPath: /etc/config
              name: cloud-config-volume
              readOnly: true
            - mountPath: /usr/libexec/kubernetes/kubelet-plugins/volume/exec
              name: flexvolume-dir
          resources:
            requests:
              cpu: 200m
      hostNetwork: true
      volumes:
      - hostPath:
          path: /usr/libexec/kubernetes/kubelet-plugins/volume/exec
          type: DirectoryOrCreate
        name: flexvolume-dir
      - hostPath:
          path: /etc/kubernetes
          type: DirectoryOrCreate
        name: k8s-certs
      - hostPath:
          path: /etc/ssl/certs
          type: DirectoryOrCreate
        name: ca-certs
      - name: cloud-config-volume
        secret:
          secretName: cloud-config
EOF

testRes="false"
for (( i=0; i<10; i++));
do
    lines=`kubectl get pod -n kube-system | grep huawei-cloud-controller-manager | grep Running | wc -l`
    if [ "$lines" = "1" ]; then
        testRes="true"
        break
    fi
    sleep 3
done

if [ "$testRes" = "true" ]; then
    echo -e "Deploy huawei-cloud-controller-manager success\n"
    exit 0
else
    echo -e "FAIL: deploy huawei-cloud-controller-manager failed\n"
    exit 1
fi
