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

// CreateSecret create Secret.
func CreateSecret(client kubernetes.Interface, secret *corev1.Secret) {
	ginkgo.By(fmt.Sprintf("Creating Secret(%s/%s)", secret.Namespace, secret.Name), func() {
		_, err := client.CoreV1().Secrets(secret.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
}

// RemoveSecret delete Secret.
func RemoveSecret(client kubernetes.Interface, namespace, name string) {
	ginkgo.By(fmt.Sprintf("Removing Secret(%s/%s)", namespace, name), func() {
		err := client.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
}
