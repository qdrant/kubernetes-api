package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	KindQdrantClusterSnapshot     = "QdrantClusterSnapshot"
	ResourceQdrantClusterSnapshot = "qdrantclustersnapshots"
)

type QdrantClusterSnapshotSpec struct {
	// The cluster ID for which a Snapshot need to be taken
	// The cluster should be in the same namespace as this QdrantClusterSnapshot is located
	ClusterId string `json:"cluster-id"`
	// The CreationTimestamp of the backup (expressed in Unix epoch format)
	// +optional
	CreationTimestamp int64 `json:"creation-timestamp,omitempty"`
	// Specifies the short Id which identifies a schedule, if any.
	// This field should not be set if the backup is made manually.
	// +kubebuilder:validation:MaxLength=8
	// +optional
	ScheduleShortId *string `json:"scheduleShortId,omitempty"`
	// The retention period of this snapshot in hours, if any.
	// If not set, the backup doesn't have a retention period, meaning it will not be removed.
	// +kubebuilder:validation:Pattern=^[0-9]+h$
	// +optional
	Retention *string `json:"retention,omitempty"`
}

type QdrantClusterSnapshotPhase string

const (
	SnapshotRunning   QdrantClusterSnapshotPhase = "Running"
	SnapshotSkipped   QdrantClusterSnapshotPhase = "Skipped"
	SnapshotFailed    QdrantClusterSnapshotPhase = "Failed"
	SnapshotSucceeded QdrantClusterSnapshotPhase = "Succeeded"
)

// +kubebuilder:pruning:PreserveUnknownFields
type QdrantClusterSnapshotStatus struct {
	// +kubebuilder:validation:Enum=Running;Skipped;Failed;Succeeded
	// +optional
	Phase QdrantClusterSnapshotPhase `json:"phase,omitempty"`
	// VolumeSnapshots is the list of volume snapshots that were created
	// +optional
	VolumeSnapshots []VolumeSnapshotInfo `json:"volumeSnapshots,omitempty"`
	// The calculated time (in UTC) this snapshot will be deleted, if so.
	// +optional
	RetainUntil *metav1.Time `json:"retainUntil,omitempty"`
	// CompletionTime specifies how long it took for the snapshot to complete
	// When serialized, it is a Duration in string format which follows "DDdHHhMMmSSs" format
	// For example: "1d3h5m10s", "3h5m10s", "5m10s", "10s" etc.
	// +optional
	CompletionTime *metav1.Duration `json:"completionTime,omitempty"`
}

type VolumeSnapshotInfo struct {
	// VolumeSnapshotName is the name of the volume snapshot
	VolumeSnapshotName string `json:"volumeSnapshotName"`
	// VolumeName is the name of the volume that was backed up
	VolumeName string `json:"volumeName"`
	// ReadyToUse indicates if the volume snapshot is ready to use
	// +optional
	ReadyToUse bool `json:"readyToUse"`
	// SnapshotHandle is the identifier of the volume snapshot in the respective cloud provider
	// +optional
	SnapshotHandle string `json:"snapshotHandle,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=qdrantclustersnapshots,singular=qdrantclustersnapshot,shortName=qcsnap;qcsnaps
// +kubebuilder:printcolumn:name="clusterid",type=string,JSONPath=`.spec.cluster-id`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="retainUntil",type=string,JSONPath=`.status.retainUntil`
// +kubebuilder:printcolumn:name="completion-time",type=string,JSONPath=`.status.completionTime`
// +kubebuilder:printcolumn:name="age",type=date,JSONPath=`.metadata.creationTimestamp`

// QdrantClusterSnapshot is the Schema for the qdrantclustersnapshots API
type QdrantClusterSnapshot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QdrantClusterSnapshotSpec   `json:"spec,omitempty"`
	Status QdrantClusterSnapshotStatus `json:"status,omitempty"`
}

func (qcs *QdrantClusterSnapshot) IsCompleted() bool {
	return qcs.Status.Phase == SnapshotSucceeded || qcs.Status.Phase == SnapshotFailed || qcs.Status.Phase == SnapshotSkipped
}

//+kubebuilder:object:root=true

// QdrantClusterSnapshotList contains a list of QdrantClusterSnapshot
type QdrantClusterSnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QdrantClusterSnapshot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QdrantClusterSnapshot{}, &QdrantClusterSnapshotList{})
}
