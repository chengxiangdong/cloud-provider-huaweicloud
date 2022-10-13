#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

count=`kubectl get -n kube-system secret cloud-config | grep cloud-config | wc -l`
if [[ "$count" -ne 1 ]]; then
  echo "Please create a secret with the name: cloud-config."
  exit 1
fi

cat << EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cloud-controller-manager
  namespace: kube-system
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: kube-system
  name: ccm-secret-role
rules:
  - apiGroups: [""]
    resources: ["secrets", "endpoints", "serviceaccounts", "services", "nodes"]
    verbs: ["get", "list", "create", "update"]
  - apiGroups:
      - coordination.k8s.io
    resources:
      - leases
    verbs:
      - create
      - get
      - list
      - update
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: ccm-secret-binding
  namespace: kube-system
subjects:
  - kind: ServiceAccount
    name: cloud-controller-manager
    namespace: kube-system
roleRef:
  kind: ClusterRole
  name: ccm-secret-role
  apiGroup: rbac.authorization.k8s.io
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: extension-apiserver-auth-reader-binding
  namespace: kube-system
subjects:
  - kind: ServiceAccount
    name: cloud-controller-manager
    namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: extension-apiserver-authentication-reader
EOF

cat << EOF | kubectl apply -f -
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
          image: docker.io/chengxiangdong/huawei-cloud-controller-manager:latest
          args:
            - /bin/huawei-cloud-controller-manager
            - --v=7
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
