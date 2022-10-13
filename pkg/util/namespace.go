/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
)

// IsNamespaceExist tells if specific already exists.
func IsNamespaceExist(client kubeclient.Interface, namespace string) (bool, error) {
	_, err := client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CreateNamespace just try to create the namespace.
func CreateNamespace(client kubeclient.Interface, namespaceObj *corev1.Namespace) (*corev1.Namespace, error) {
	_, err := client.CoreV1().Namespaces().Create(context.TODO(), namespaceObj, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return namespaceObj, nil
		}

		return nil, err
	}

	return namespaceObj, nil
}

// DeleteNamespace just try to delete the namespace.
func DeleteNamespace(client kubeclient.Interface, namespace string) error {
	err := client.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}

// EnsureNamespaceExist makes sure that the specific namespace exist in cluster.
// If namespace not exit, just create it.
func EnsureNamespaceExist(client kubeclient.Interface, namespace string, dryRun bool) (*corev1.Namespace, error) {
	namespaceObj := &corev1.Namespace{}
	namespaceObj.ObjectMeta.Name = namespace

	if dryRun {
		return namespaceObj, nil
	}

	exist, err := IsNamespaceExist(client, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to check if namespace exist. namespace: %s, error: %v", namespace, err)
	}
	if exist {
		return namespaceObj, nil
	}

	createdObj, err := CreateNamespace(client, namespaceObj)
	if err != nil {
		return nil, fmt.Errorf("ensure namespace failed due to create failed. namespace: %s, error: %v", namespace, err)
	}

	return createdObj, nil
}
