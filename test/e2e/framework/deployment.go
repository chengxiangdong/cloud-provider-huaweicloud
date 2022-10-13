package framework

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateDeployment create Deployment.
func CreateDeployment(client kubernetes.Interface, deployment *appsv1.Deployment) {
	ginkgo.By(fmt.Sprintf("Creating Deployment(%s/%s)", deployment.Namespace, deployment.Name), func() {
		_, err := client.AppsV1().Deployments(deployment.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
}

// RemoveDeployment delete Deployment.
func RemoveDeployment(client kubernetes.Interface, namespace, name string) {
	ginkgo.By(fmt.Sprintf("Removing Deployment(%s/%s)", namespace, name), func() {
		err := client.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
}

// UpdateDeploymentReplicas update deployment's replicas.
func UpdateDeploymentReplicas(client kubernetes.Interface, deployment *appsv1.Deployment, replicas int32) {
	ginkgo.By(fmt.Sprintf("Updating Deployment(%s/%s)'s replicas to %d", deployment.Namespace, deployment.Name, replicas), func() {
		deployment.Spec.Replicas = &replicas
		gomega.Eventually(func() error {
			_, err := client.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
			return err
		}, pollTimeout, pollInterval).ShouldNot(gomega.HaveOccurred())
	})
}

// UpdateDeploymentAnnotations update deployment's annotations.
func UpdateDeploymentAnnotations(client kubernetes.Interface, deployment *appsv1.Deployment, annotations map[string]string) {
	ginkgo.By(fmt.Sprintf("Updating Deployment(%s/%s)'s annotations to %v", deployment.Namespace, deployment.Name, annotations), func() {
		deployment.Annotations = annotations
		gomega.Eventually(func() error {
			_, err := client.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
			return err
		}, pollTimeout, pollInterval).ShouldNot(gomega.HaveOccurred())
	})
}

// UpdateDeploymentVolumes update Deployment's volumes.
func UpdateDeploymentVolumes(client kubernetes.Interface, deployment *appsv1.Deployment, volumes []corev1.Volume) {
	ginkgo.By(fmt.Sprintf("Updating Deployment(%s/%s)'s volumes", deployment.Namespace, deployment.Name), func() {
		deployment.Spec.Template.Spec.Volumes = volumes
		gomega.Eventually(func() error {
			_, err := client.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
			return err
		}, pollTimeout, pollInterval).ShouldNot(gomega.HaveOccurred())
	})
}

// UpdateDeploymentServiceAccountName update Deployment's serviceAccountName.
func UpdateDeploymentServiceAccountName(client kubernetes.Interface, deployment *appsv1.Deployment, serviceAccountName string) {
	ginkgo.By(fmt.Sprintf("Updating Deployment(%s/%s)'s serviceAccountName", deployment.Namespace, deployment.Name), func() {
		deployment.Spec.Template.Spec.ServiceAccountName = serviceAccountName
		gomega.Eventually(func() error {
			_, err := client.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
			return err
		}, pollTimeout, pollInterval).ShouldNot(gomega.HaveOccurred())
	})
}
