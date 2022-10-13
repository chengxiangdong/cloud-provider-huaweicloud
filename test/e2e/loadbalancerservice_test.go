package e2e

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/rand"

	"sigs.k8s.io/cloud-provider-huaweicloud/pkg/helper"
	"sigs.k8s.io/cloud-provider-huaweicloud/test/e2e/framework"
)

var _ = ginkgo.Describe("loadbalancer service testing", func() {
	var deployment *appsv1.Deployment
	var service *corev1.Service

	ginkgo.BeforeEach(func() {
		deploymentName := deploymentNamePrefix + rand.String(RandomStrLength)
		deployment = helper.NewDeployment(testNamespace, deploymentName)
		framework.CreateDeployment(kubeClient, deployment)
	})

	ginkgo.AfterEach(func() {
		framework.RemoveDeployment(kubeClient, deployment.Namespace, deployment.Name)
		if service != nil {
			framework.RemoveService(kubeClient, service.Namespace, service.Name)
			ginkgo.By(fmt.Sprintf("Wait for the Service(%s/%s) to be deleted", testNamespace, service.Name), func() {
				gomega.Eventually(func(g gomega.Gomega) (bool, error) {
					_, err := kubeClient.CoreV1().Services(testNamespace).Get(context.TODO(), service.Name, metav1.GetOptions{})
					if apierrors.IsNotFound(err) {
						return true, nil
					}
					if err != nil {
						return false, err
					}
					return false, nil
				}, pollTimeout, pollInterval).Should(gomega.Equal(true))
			})
		}
	})

	ginkgo.It("service enhanced testing", func() {
		// todo(chengxiangdong): Refactor to use dynamic creation of ELB and VPC later.
		sessionAffinity := "SOURCE_IP"
		elbID := "d9b4d06b-8813-46db-bea2-28c5fcbc16f5"
		subnetID := "82e61d5f-674f-4ef4-a665-d5785a0733e0"
		lbIP := "192.168.0.55"
		serviceName := serviceNamePrefix + rand.String(RandomStrLength)
		service = newLoadbalancerService(testNamespace, serviceName, sessionAffinity, elbID, subnetID, lbIP)
		framework.CreateService(kubeClient, service)

		var ingress string
		ginkgo.By("Check service status", func() {
			gomega.Eventually(func(g gomega.Gomega) (bool, error) {
				svc, err := kubeClient.CoreV1().Services(testNamespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
				g.Expect(err).ShouldNot(gomega.HaveOccurred())

				if len(svc.Status.LoadBalancer.Ingress) > 0 {
					ingress = svc.Status.LoadBalancer.Ingress[0].IP
					g.Expect(ingress).Should(gomega.Equal(lbIP))
					return true, nil
				}

				return false, nil
			}, pollTimeout, pollInterval).Should(gomega.Equal(true))
		})

		ginkgo.By("Check if ELB listener is available", func() {
			url := fmt.Sprintf("http://%s", ingress)
			statusCode, err := helper.DoRequest(url)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			gomega.Expect(statusCode).Should(gomega.Equal(200))
		})
	})
})

// newLoadbalancerService new a loadbalancer type service
func newLoadbalancerService(namespace, name, sessionAffinity, elbID, subnetID, lbIP string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Annotations: map[string]string{
				"kubernetes.io/elb.class":             "union",
				"kubernetes.io/session-affinity-mode": sessionAffinity,
				"kubernetes.io/elb.id":                elbID,
				"kubernetes.io/elb.subnet-id":         subnetID,
			},
			Labels: map[string]string{"app": "nginx"},
		},
		Spec: corev1.ServiceSpec{
			LoadBalancerIP:        lbIP,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeCluster,
			Type:                  corev1.ServiceTypeLoadBalancer,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.IntOrString{IntVal: 80},
				},
			},
			Selector: map[string]string{"app": "nginx"},
		},
	}
}
