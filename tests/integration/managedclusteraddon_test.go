// +build integration

/*
Copyright 2021 Red Hat OpenShift Data Foundation.

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

package integration_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tokenexchange "github.com/red-hat-storage/odf-multicluster-orchestrator/addons/token-exchange"
	multiclusterv1alpha1 "github.com/red-hat-storage/odf-multicluster-orchestrator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	addonapiv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
)

var (
	mirrorPeer1 = multiclusterv1alpha1.MirrorPeer{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mirrorpeer1",
		},
	}
	mirrorPeer1LookupKey = types.NamespacedName{Namespace: mirrorPeer1.Namespace, Name: mirrorPeer1.Name}
	mcAddOn1LookupKey    = types.NamespacedName{Namespace: "cluster1", Name: tokenexchange.TokenExchangeName}
	mcAddOn2LookupKey    = types.NamespacedName{Namespace: "cluster2", Name: tokenexchange.TokenExchangeName}
)

var _ = Describe("ManagedClusterAddOn creation, updation and deletion", func() {
	When("creating or updating ManagedClusterAddOn", func() {
		It("should not return validation error", func() {
			newMirrorPeer := mirrorPeer1.DeepCopy()
			newMirrorPeer.Spec = multiclusterv1alpha1.MirrorPeerSpec{
				Items: []multiclusterv1alpha1.PeerRef{
					{
						ClusterName: "cluster1",
						StorageClusterRef: multiclusterv1alpha1.StorageClusterRef{
							Name:      "test-storagecluster",
							Namespace: "test-namespace",
						},
					},
					{
						ClusterName: "cluster2",
						StorageClusterRef: multiclusterv1alpha1.StorageClusterRef{
							Name:      "test-storagecluster",
							Namespace: "test-namespace",
						},
					},
				},
			}
			var mcAddOn1 addonapiv1alpha1.ManagedClusterAddOn
			var mcAddOn2 addonapiv1alpha1.ManagedClusterAddOn
			By("creating MirrorPeer. Also, this should automatically create ManagedClusterAddOn", func() {
				err := k8sClient.Create(context.TODO(), newMirrorPeer, &client.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
			})
			By("polling for the created ManagedClusterAddOn", func() {
				Eventually(func() error {
					err := k8sClient.Get(context.TODO(), mcAddOn1LookupKey, &mcAddOn1)
					if err != nil {
						return err
					}
					return fmt.Errorf("Waiting on ManagedClusterAddOn %s/%s to get created", mcAddOn1LookupKey.Namespace, mcAddOn1LookupKey.Name)
				}, 60*time.Second, 2*time.Second).ShouldNot(HaveOccurred())

				Eventually(func() error {
					err := k8sClient.Get(context.TODO(), mcAddOn2LookupKey, &mcAddOn2)
					if err != nil {
						return err
					}
					return fmt.Errorf("Waiting on ManagedClusterAddOn %s/%s to get created", mcAddOn2LookupKey.Namespace, mcAddOn2LookupKey.Name)
				}, 60*time.Second, 2*time.Second).ShouldNot(HaveOccurred())
			})
			By("updating ManagedClusterAddOn", func() {
				mcAddOn1.Spec.InstallNamespace = "new-test-namespace"
				err := k8sClient.Update(context.TODO(), &mcAddOn1, &client.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())

				err = k8sClient.Delete(context.TODO(), &mcAddOn1, &client.DeleteOptions{})
				Expect(err).NotTo(HaveOccurred())
				err = k8sClient.Delete(context.TODO(), &mcAddOn2, &client.DeleteOptions{})
				Expect(err).NotTo(HaveOccurred())
				err = k8sClient.Delete(context.TODO(), newMirrorPeer, &client.DeleteOptions{})
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
