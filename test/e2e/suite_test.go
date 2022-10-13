package e2e

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cloud-provider-huaweicloud/pkg/helper"
	"sigs.k8s.io/cloud-provider-huaweicloud/pkg/util"
)

const (
	// RandomStrLength represents the random string length to combine names.
	RandomStrLength = 5
)

const (
	deploymentNamePrefix         = "deploy-"
	serviceNamePrefix            = "service-"
	podNamePrefix                = "pod-"
	crdNamePrefix                = "cr-"
	jobNamePrefix                = "job-"
	workloadNamePrefix           = "workload-"
	federatedResourceQuotaPrefix = "frq-"
	configMapNamePrefix          = "configmap-"
	secretNamePrefix             = "secret-"
	pvcNamePrefix                = "pvc-"
	saNamePrefix                 = "sa-"
	ingressNamePrefix            = "ingress-"
	daemonSetNamePrefix          = "daemonset-"
	statefulSetNamePrefix        = "statefulset-"
	roleNamePrefix               = "role-"
	clusterRoleNamePrefix        = "clusterrole-"
	roleBindingNamePrefix        = "rolebinding-"
	clusterRoleBindingNamePrefix = "clusterrolebinding-"

	updateDeploymentReplicas  = 6
	updateStatefulSetReplicas = 6
	updateServicePort         = 81
	updatePodImage            = "nginx:latest"
	updateCRnamespace         = "e2e-test"
	updateBackoffLimit        = 3
	updateParallelism         = 3
)

var (
	// pollInterval defines the interval time for a poll operation.
	pollInterval time.Duration
	// pollTimeout defines the time after which the poll operation times out.
	pollTimeout time.Duration
)

var (
	kubeConfig         string
	restConfig         *rest.Config
	kubeClient         kubernetes.Interface
	dynamicClient      dynamic.Interface
	controlPlaneClient client.Client
	testNamespace      string
)

func init() {
	// usage ginkgo -- --poll-interval=5s --pollTimeout=5m
	// eg. ginkgo -v --race --trace --fail-fast -p --randomize-all ./test/e2e/ -- --poll-interval=5s --pollTimeout=5m
	flag.DurationVar(&pollInterval, "poll-interval", 5*time.Second, "poll-interval defines the interval time for a poll operation")
	flag.DurationVar(&pollTimeout, "poll-timeout", 300*time.Second, "poll-timeout defines the time which the poll operation times out")
}

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Skip)
	ginkgo.RunSpecs(t, "E2E Suite")
}

var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	return nil
}, func(bytes []byte) {
	var err error

	kubeConfig = os.Getenv("KUBECONFIG")
	gomega.Expect(kubeConfig).ShouldNot(gomega.BeEmpty())

	restConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	kubeClient, err = kubernetes.NewForConfig(restConfig)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	dynamicClient, err = dynamic.NewForConfig(restConfig)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	controlPlaneClient, err = newForConfigOrDie(restConfig)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	testNamespace = "ccm-test" //fmt.Sprintf("ccmtest-%s", rand.String(RandomStrLength))
	testNamespace = fmt.Sprintf("ccmtest-%s", rand.String(RandomStrLength))
	err = setupTestNamespace(testNamespace, kubeClient)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
})

var _ = ginkgo.AfterSuite(func() {
	// cleanup all namespaces we created both in control plane and member clusters.
	// It will not return error even if there is no such namespace in there that may happen in case setup failed.
	err := cleanupTestNamespace(testNamespace, kubeClient)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
})

// cleanupTestNamespace will remove the namespace we setup before for the whole testing.
func cleanupTestNamespace(namespace string, kubeClient kubernetes.Interface) error {
	err := util.DeleteNamespace(kubeClient, namespace)
	if err != nil {
		return err
	}
	return nil
}

// setupTestNamespace will create a namespace in control plane and all member clusters, most of cases will run against it.
// The reason why we need a separated namespace is it will make it easier to cleanup resources deployed by the testing.
func setupTestNamespace(namespace string, kubeClient kubernetes.Interface) error {
	namespaceObj := helper.NewNamespace(namespace)
	_, err := util.CreateNamespace(kubeClient, namespaceObj)
	if err != nil {
		return err
	}
	return nil
}

func newForConfigOrDie(config *rest.Config) (client.Client, error) {
	c, err := client.New(config, client.Options{
		Scheme: runtime.NewScheme(),
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}
