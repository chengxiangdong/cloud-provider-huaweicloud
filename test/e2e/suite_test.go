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

package e2e

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	ginkgov2 "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"sigs.k8s.io/cloud-provider-huaweicloud/pkg/helper"
	"sigs.k8s.io/cloud-provider-huaweicloud/pkg/util"
)

const (
	// RandomStrLength represents the random string length to combine names.
	RandomStrLength = 5
)

const (
	deploymentNamePrefix = "deploy-"
	serviceNamePrefix    = "service-"
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
	testNamespace      string
)

func init() {
	// usage ginkgo -- --poll-interval=5s --pollTimeout=5m
	// eg. ginkgo -v --race --trace --fail-fast -p --randomize-all ./test/e2e/ -- --poll-interval=5s --pollTimeout=5m
	flag.DurationVar(&pollInterval, "poll-interval", 5*time.Second, "poll-interval defines the interval time for a poll operation")
	flag.DurationVar(&pollTimeout, "poll-timeout", 300*time.Second, "poll-timeout defines the time which the poll operation times out")
}

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgov2.Skip)
	ginkgov2.RunSpecs(t, "E2E Suite")
}

var _ = ginkgov2.SynchronizedBeforeSuite(func() []byte {
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

	testNamespace = fmt.Sprintf("ccmtest-%s", rand.String(RandomStrLength))
	err = setupTestNamespace(testNamespace, kubeClient)
	gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
})

var _ = ginkgov2.AfterSuite(func() {
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
