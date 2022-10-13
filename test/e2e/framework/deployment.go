package framework

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
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
