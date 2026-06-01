package v1

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("QdrantCluster", func() {
	Context("API integration tests", func() {
		const namespaceName = "test-namespace"
		ctx := context.Background()
		BeforeEach(func() {
			By("Creating the Namespace to perform the tests")
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespaceName,
				},
			}
			err := k8sClient.Create(ctx, namespace)
			Expect(client.IgnoreAlreadyExists(err)).To(Not(HaveOccurred()))
		})
		It("should not flip ServicePerNode value on update", func() {
			qc := QdrantCluster{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespaceName,
					Name:      "test-cluster",
				},
				Spec: QdrantClusterSpec{
					Id:             "test-cluster",
					Size:           1,
					ServicePerNode: NewPointer(false),
				},
			}
			err := k8sClient.Create(ctx, &qc)
			Expect(err).To(Not(HaveOccurred()))
			Expect(DerefPointer(qc.Spec.ServicePerNode)).To(BeFalse())

			qc.Spec.Size = 2
			err = k8sClient.Update(ctx, &qc)
			Expect(err).To(Not(HaveOccurred()))
			Expect(DerefPointer(qc.Spec.ServicePerNode)).To(BeFalse())
		})
		It("should default OnDemandReplication to Off when omitted", func() {
			qc := QdrantCluster{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespaceName,
					Name:      "test-cluster-on-demand-replication",
				},
				Spec: QdrantClusterSpec{
					Id:      "test-cluster-on-demand-replication",
					Version: "v1.18.0",
					Size:    1,
				},
			}
			err := k8sClient.Create(ctx, &qc)
			Expect(err).To(Not(HaveOccurred()))

			created := &QdrantCluster{}
			err = k8sClient.Get(ctx, client.ObjectKeyFromObject(&qc), created)
			Expect(err).To(Not(HaveOccurred()))
			Expect(created.Spec.OnDemandReplication).To(Equal(OnDemandReplicationOff))
		})
	})

})

// NewPointer is a generic function to create a pointer to any type.
func NewPointer[T any](value T) *T {
	return &value
}

// DerefPointer is a generic function to dereference a pointer with a default-value fallback.
func DerefPointer[T any](ptr *T, defaults ...T) T {
	if ptr != nil {
		return *ptr
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	var empty T
	return empty
}
