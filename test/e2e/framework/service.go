package framework

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateService create Service.
func CreateService(client kubernetes.Interface, service *corev1.Service) {
	ginkgo.By(fmt.Sprintf("Creating Service(%s/%s)", service.Namespace, service.Name), func() {
		_, err := client.CoreV1().Services(service.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
}

// RemoveService delete Service.
func RemoveService(client kubernetes.Interface, namespace, name string) {
	ginkgo.By(fmt.Sprintf("Removing Service(%s/%s)", namespace, name), func() {
		err := client.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
}
