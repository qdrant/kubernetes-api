package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//goland:noinspection GoUnusedConst
const (
	KindQdrantNodePoolVersion     = "QdrantNodePoolVersion"
	ResourceQdrantNodePoolVersion = "qdrantnodepoolversions"
)

// NodePoolEndOfLifeSeverity indicates how urgently a node pool version should be retired.
type NodePoolEndOfLifeSeverity string

//goland:noinspection GoUnusedConst
const (
	NodePoolEndOfLifeSeverityNormal   NodePoolEndOfLifeSeverity = "Normal"
	NodePoolEndOfLifeSeverityCritical NodePoolEndOfLifeSeverity = "Critical"
)

// NodePoolEndOfLife describes when a node pool version is scheduled for retirement.
type NodePoolEndOfLife struct {
	// Date the version is retired.
	Date metav1.Time `json:"date"`
	// Severity of the end-of-life signal.
	// +kubebuilder:validation:Enum=Normal;Critical
	// +kubebuilder:default=Normal
	// +optional
	Severity NodePoolEndOfLifeSeverity `json:"severity,omitempty"`
}

// QdrantNodePoolVersionSpec defines the desired state of QdrantNodePoolVersion.
type QdrantNodePoolVersionSpec struct {
	// Node pool version, matching the qdrant.io/node-pool-version label and taint on the nodes.
	// +kubebuilder:validation:Pattern=`^v[0-9]+\.[0-9]+\.[0-9]+$`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="spec.version is immutable."
	Version string `json:"version"`
	// Kubernetes version the node pool runs.
	// +kubebuilder:validation:Pattern=`^[0-9]+\.[0-9]+(\.[0-9]+)?$`
	KubernetesVersion string `json:"kubernetesVersion"`
	// Whether new workloads may be scheduled onto this node pool version.
	Available bool `json:"available"`
	// Absent until the version is scheduled for retirement.
	// +optional
	EndOfLife *NodePoolEndOfLife `json:"endOfLife,omitempty"`
}

// QdrantNodePoolVersionStatus defines the observed state of QdrantNodePoolVersion.
// +kubebuilder:pruning:PreserveUnknownFields
type QdrantNodePoolVersionStatus struct {
	// ObservedGeneration is the most recent generation observed by the operator.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=qdrantnodepoolversions,scope=Cluster,singular=qdrantnodepoolversion,shortName=qnpv
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Kubernetes",type=string,JSONPath=`.spec.kubernetesVersion`
// +kubebuilder:printcolumn:name="Available",type=boolean,JSONPath=`.spec.available`
// +kubebuilder:printcolumn:name="EOL",type=string,JSONPath=`.spec.endOfLife.date`
// +kubebuilder:printcolumn:name="Severity",type=string,JSONPath=`.spec.endOfLife.severity`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'node-pool-version-' + self.spec.version",message="metadata.name must be node-pool-version-<spec.version>."

// QdrantNodePoolVersion describes a versioned node pool (qdrant.io/node-pool-version) available in the region.
type QdrantNodePoolVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec   QdrantNodePoolVersionSpec   `json:"spec"`
	Status QdrantNodePoolVersionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// QdrantNodePoolVersionList contains a list of QdrantNodePoolVersion
type QdrantNodePoolVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QdrantNodePoolVersion `json:"items"`
}

func init() {
	registerTypes(&QdrantNodePoolVersion{}, &QdrantNodePoolVersionList{})
}
