package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//goland:noinspection GoUnusedConst
const (
	KindQdrantClusterScheduledSnapshot     = "QdrantClusterScheduledSnapshot"
	ResourceQdrantClusterScheduledSnapshot = "qdrantclusterscheduledsnapshots"
)

// QdrantClusterScheduledSnapshotSpec defines the desired state of QdrantCluster
type QdrantClusterScheduledSnapshotSpec struct {
	// Id specifies the unique identifier of the cluster
	ClusterId string `json:"cluster-id"`
	// Specifies short Id which identifies a schedule
	// +kubebuilder:validation:MaxLength=8
	ScheduleShortId string `json:"scheduleShortId"`
	// Cron expression for frequency of creating snapshots, see https://en.wikipedia.org/wiki/Cron.
	// The schedule is specified in UTC.
	// +kubebuilder:validation:Pattern=`^(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|([\d\*]+(\/|-)\d+)|\d+|\*) ?){5,7})$`
	Schedule string `json:"schedule"`
	// Retention of schedule in hours
	// +kubebuilder:validation:Pattern=^[0-9]+h$
	Retention string `json:"retention"`
}

type ScheduledSnapshotPhase string

//goland:noinspection GoUnusedConst
const (
	ScheduleActive   ScheduledSnapshotPhase = "Active"
	ScheduleDisabled ScheduledSnapshotPhase = "Disabled"
)

// QdrantClusterScheduledSnapshotStatus defines the observed state of the snapshot
type QdrantClusterScheduledSnapshotStatus struct {
	// Phase is the current phase of the scheduled snapshot
	// +kubebuilder:validation:Enum=Active;Disabled
	Phase ScheduledSnapshotPhase `json:"phase,omitempty"`
	// The next scheduled time in UTC
	Scheduled metav1.Time `json:"scheduled,omitempty"`
	// Message from the operator in case of failures, like schedule not valid
	// +optional
	Message *string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
// +kubebuilder:resource:path=qdrantclusterscheduledsnapshots,singular=qdrantclusterscheduledsnapshot,shortName=qcssnap;qcssnaps
// +kubebuilder:printcolumn:name="clusterid",type=string,JSONPath=`.spec.cluster-id`
// +kubebuilder:printcolumn:name="scheduleShortId",type=string,JSONPath=`.spec.scheduleShortId`
// +kubebuilder:printcolumn:name="schedule",type=string,JSONPath=`.spec.schedule`
// +kubebuilder:printcolumn:name="retention",type=string,JSONPath=`.spec.retention`
// +kubebuilder:printcolumn:name="scheduled",type=string,JSONPath=`.status.scheduled`
// +kubebuilder:printcolumn:name="age",type=date,JSONPath=`.metadata.creationTimestamp`

// QdrantClusterScheduledSnapshot is the Schema for the qdrantclusterscheduledsnapshots API
type QdrantClusterScheduledSnapshot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QdrantClusterScheduledSnapshotSpec   `json:"spec,omitempty"`
	Status QdrantClusterScheduledSnapshotStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// QdrantClusterScheduledSnapshotList contains a list of QdrantCluster
type QdrantClusterScheduledSnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QdrantClusterScheduledSnapshot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QdrantClusterScheduledSnapshot{}, &QdrantClusterScheduledSnapshotList{})
}
