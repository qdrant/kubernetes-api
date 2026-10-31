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
		It("should accept disabled as a rebalanceStrategy", func() {
			qc := QdrantCluster{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespaceName,
					Name:      "test-cluster-rebalance-disabled",
				},
				Spec: QdrantClusterSpec{
					Id:                "test-cluster-rebalance-disabled",
					Size:              1,
					RebalanceStrategy: NewPointer(Disabled),
				},
			}
			err := k8sClient.Create(ctx, &qc)
			Expect(err).To(Not(HaveOccurred()))

			created := &QdrantCluster{}
			err = k8sClient.Get(ctx, client.ObjectKeyFromObject(&qc), created)
			Expect(err).To(Not(HaveOccurred()))
			Expect(DerefPointer(created.Spec.RebalanceStrategy)).To(Equal(Disabled))
		})
		It("should reject an invalid rebalanceStrategy", func() {
			qc := QdrantCluster{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespaceName,
					Name:      "test-cluster-rebalance-invalid",
				},
				Spec: QdrantClusterSpec{
					Id:                "test-cluster-rebalance-invalid",
					Size:              1,
					RebalanceStrategy: NewPointer(RebalanceStrategy("bogus")),
				},
			}
			err := k8sClient.Create(ctx, &qc)
			Expect(err).To(HaveOccurred())
		})
		snapshotsCluster := func(name string, snapshots *SnapshotsConfig) *QdrantCluster {
			return &QdrantCluster{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespaceName,
					Name:      name,
				},
				Spec: QdrantClusterSpec{
					Id:   name,
					Size: 1,
					Config: &QdrantConfiguration{
						Storage: &StorageConfig{Snapshots: snapshots},
					},
				},
			}
		}
		secretKey := func(key string) corev1.SecretKeySelector {
			return corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: "s3-keys"},
				Key:                  key,
			}
		}
		It("should accept S3 snapshot storage with a Secret reference", func() {
			qc := snapshotsCluster("test-cluster-snapshots-s3", &SnapshotsConfig{
				S3: &SnapshotsS3Config{
					Bucket:      "snapshots",
					EndpointURL: NewPointer("http://rustfs:9000"),
					Credentials: &SnapshotsS3Credentials{
						SecretRef: &SnapshotsS3SecretRef{
							AccessKeyID:     secretKey("access-key-id"),
							SecretAccessKey: secretKey("secret-access-key"),
						},
					},
				},
			})
			Expect(k8sClient.Create(ctx, qc)).To(Succeed())
		})
		It("should reject snapshot storage without a backend", func() {
			qc := snapshotsCluster("test-cluster-snapshots-empty", &SnapshotsConfig{})
			err := k8sClient.Create(ctx, qc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("exactly one snapshot storage backend must be set"))
		})
		It("should reject S3 snapshot storage with two kinds of credentials", func() {
			qc := snapshotsCluster("test-cluster-snapshots-two-creds", &SnapshotsConfig{
				S3: &SnapshotsS3Config{
					Bucket: "snapshots",
					Credentials: &SnapshotsS3Credentials{
						SecretRef: &SnapshotsS3SecretRef{
							AccessKeyID:     secretKey("access-key-id"),
							SecretAccessKey: secretKey("secret-access-key"),
						},
						ContainerCredentials: &SnapshotsContainerCredentials{
							URI:            "http://crosp/v1/credentials",
							TokenSecretRef: corev1.LocalObjectReference{Name: "token"},
						},
					},
				},
			})
			err := k8sClient.Create(ctx, qc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at most one of secretRef and containerCredentials may be set"))
		})
		It("should reject S3 snapshot storage with an empty bucket", func() {
			qc := snapshotsCluster("test-cluster-snapshots-no-bucket", &SnapshotsConfig{
				S3: &SnapshotsS3Config{},
			})
			Expect(k8sClient.Create(ctx, qc)).NotTo(Succeed())
		})
		It("should accept S3 retry settings", func() {
			qc := snapshotsCluster("test-cluster-snapshots-retries", &SnapshotsConfig{
				S3: &SnapshotsS3Config{
					Bucket:            "snapshots",
					TimeoutSec:        NewPointer(uint(60)),
					ConnectTimeoutSec: NewPointer(uint(10)),
					MaxRetries:        NewPointer(uint(0)),
					RetryTimeoutSec:   NewPointer(uint(600)),
				},
			})
			Expect(k8sClient.Create(ctx, qc)).To(Succeed())
		})
		It("should reject an S3 request timeout of zero", func() {
			qc := snapshotsCluster("test-cluster-snapshots-zero-timeout", &SnapshotsConfig{
				S3: &SnapshotsS3Config{Bucket: "snapshots", TimeoutSec: NewPointer(uint(0))},
			})
			err := k8sClient.Create(ctx, qc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("timeout_sec"))
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
