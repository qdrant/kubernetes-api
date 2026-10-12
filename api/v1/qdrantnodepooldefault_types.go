package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//goland:noinspection GoUnusedConst
const (
	KindQdrantNodePoolDefault     = "QdrantNodePoolDefault"
	ResourceQdrantNodePoolDefault = "qdrantnodepooldefaults"
	// QdrantNodePoolDefaultName is the required metadata.name of the singleton QdrantNodePoolDefault.
	QdrantNodePoolDefaultName = "default"
)

//goland:noinspection GoUnusedConst
const (
	// QdrantNodePoolLabel is the label and taint key on nodes for a versioned node pool.
	// The operator also uses this key for the pod nodeSelector and matching toleration.
	QdrantNodePoolLabel = "qdrant.io/node-pool-version"
	// QdrantNodePoolAnnotation is the annotation key on pods recording which node pool
	// version the pod was rendered for. The value is the target version string (e.g. "v1.0.0").
	QdrantNodePoolAnnotation = "qdrant.io/node-pool-version"
)

// QdrantNodePoolDefaultSpec defines the desired state of QdrantNodePoolDefault.
type QdrantNodePoolDefaultSpec struct {
	// Node pool version that new and recreated workloads target. Not checked against the QdrantNodePoolVersion CRs yet.
	// +kubebuilder:validation:Pattern=`^v[0-9]+\.[0-9]+\.[0-9]+$`
	Version string `json:"version"`
}

// QdrantNodePoolDefaultStatus defines the observed state of QdrantNodePoolDefault.
// +kubebuilder:pruning:PreserveUnknownFields
type QdrantNodePoolDefaultStatus struct {
	// ObservedGeneration is the most recent generation observed by the operator.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=qdrantnodepooldefaults,scope=Cluster,singular=qdrantnodepooldefault,shortName=qnpd
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'default'",message="QdrantNodePoolDefault is a singleton and must be named \"default\"."

// QdrantNodePoolDefault is a singleton (named "default") selecting the node pool version new workloads are scheduled onto.
type QdrantNodePoolDefault struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec   QdrantNodePoolDefaultSpec   `json:"spec"`
	Status QdrantNodePoolDefaultStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// QdrantNodePoolDefaultList contains a list of QdrantNodePoolDefault
type QdrantNodePoolDefaultList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QdrantNodePoolDefault `json:"items"`
}

func init() {
	registerTypes(&QdrantNodePoolDefault{}, &QdrantNodePoolDefaultList{})
}
