package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//goland:noinspection GoUnusedConst
const (
	KindQdrantClusterRestore     = "QdrantClusterRestore"
	ResourceQdrantClusterRestore = "qdrantclusterrestores"
)

// QdrantClusterRestoreSpec defines the desired state of QdrantClusterRestore
type QdrantClusterRestoreSpec struct {
	// Source defines the source snapshot from which the restore will be done
	Source RestoreSource `json:"source"`
	// Destination defines the destination cluster where the source data will end up
	Destination RestoreDestination `json:"destination"`
}

type RestoreSource struct {
	// SnapshotName is the name of the snapshot from which we wish to restore
	SnapshotName string `json:"snapshotName"`
	// Namespace of the snapshot
	Namespace string `json:"namespace"`
}

type RestoreDestination struct {
	// Name of the destination cluster
	Name string `json:"name"`
	// Namespace of the destination cluster
	Namespace string `json:"namespace"`
	// Create when set to true indicates that
	// a new cluster with the specified name should be created.
	// Otherwise, if set to false, the existing cluster is going to be restored
	// to the specified state.
	// +optional
	Create bool `json:"create"`
}

type RestorePhase string

//goland:noinspection GoUnusedConst
const (
	RestoreRunning   RestorePhase = "Running"
	RestoreSkipped   RestorePhase = "Skipped"
	RestoreFailed    RestorePhase = "Failed"
	RestoreSucceeded RestorePhase = "Succeeded"
	RestorePending   RestorePhase = "Pending"
)

// QdrantClusterRestoreStatus defines the observed state of QdrantClusterRestore
// +kubebuilder:pruning:PreserveUnknownFields
type QdrantClusterRestoreStatus struct {
	// Phase is the current phase of the restore
	// +kubebuilder:validation:Enum=Running;Skipped;Failed;Succeeded;Pending
	Phase RestorePhase `json:"phase,omitempty"`
	// Message from the operator in case of failures, like snapshot not found
	// +optional
	Message *string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=qdrantclusterrestores,singular=qdrantclusterrestore,shortName=qcrs;qcr
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// QdrantClusterRestore is the Schema for the qdrantclusterrestores API
type QdrantClusterRestore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QdrantClusterRestoreSpec   `json:"spec,omitempty"`
	Status QdrantClusterRestoreStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// QdrantClusterRestoreList contains a list of QdrantClusterRestore objects
type QdrantClusterRestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QdrantClusterRestore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QdrantClusterRestore{}, &QdrantClusterRestoreList{})
}
