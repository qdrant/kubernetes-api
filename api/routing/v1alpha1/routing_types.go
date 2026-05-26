package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:object:root=true
// +kubebuilder:resource:path=qdrantclusterroutes,singular=qdrantclusterroute,shortName=qcrt;qcrts
// +kubebuilder:printcolumn:name="Enabled",type=boolean,JSONPath=`.spec.enabled`
// +kubebuilder:printcolumn:name="Shared",type=boolean,JSONPath=`.spec.shared`
// +kubebuilder:printcolumn:name="Dedicated",type=boolean,JSONPath=`.spec.dedicated`
// +kubebuilder:printcolumn:name="Bootstrapped",type=string,JSONPath=`.status.bootstrapped`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:subresource:status

// QdrantClusterRouting is the Schema for the routing towards Qdrant clusters API
type QdrantClusterRouting struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QdrantClusterRoutingSpec   `json:"spec"`
	Status QdrantClusterRoutingStatus `json:"status,omitempty"`
}

// QdrantClusterRoutingSpec describes the configuration for routing towards Qdrant clusters.
type QdrantClusterRoutingSpec struct {
	// ClusterId specifies the unique identifier of the cluster.
	ClusterId string `json:"clusterId"`
	// The fully qualified domain name (also know as host).
	// For shared routing this will be used for SNI resolving.
	FQDN string `json:"fqdn"`
	// Enabled specifies whether to enable ingress for the cluster or not.
	// +kubebuilder:default=true
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// Set if the cluster uses (at least one) shared loadbalancer.
	// Note that this doesn't mean it doesn't have a dedicated loadbalancer as well (e.g. during a migration from one to the other).
	// +optional
	Shared *bool `json:"shared,omitempty"`
	// Set if the cluster uses (at least one) dedicated loadbalancer.
	// Note that this doesn't mean it doesn't have a shared loadbalancer as well (e.g. during a migration from one to the other).
	// +optional
	Dedicated *bool `json:"dedicated,omitempty"`
	// TLS specifies whether tls is enabled or not at qdrant level.
	// +optional
	TLS *bool `json:"tls,omitempty"`
	// ServicePerNode specifies whether the cluster should have a dedicated route for each node.
	// +kubebuilder:default=true
	// +optional
	ServicePerNode *bool `json:"servicePerNode,omitempty"`
	// NodeIndexes specifies the indexes of the individual nodes in the cluster.
	NodeIndexes []int `json:"nodeIndexes,omitempty"`
	// AllowedSourceRanges specifies the allowed CIDR source ranges for the ingress.
	// +optional
	AllowedSourceRanges []string `json:"allowedSourceRanges,omitempty"`
	// If true enable (proxy) access log for this qdrant cluster.
	// +optional
	EnableAccessLog *bool `json:"enableAccessLog,omitempty"`
	// MultiAZ is true when the Qdrant cluster spans multiple availability
	// zones and traffic should be kept same-zone where possible.
	// +kubebuilder:default=false
	// +optional
	MultiAZ bool `json:"multiAZ,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QdrantClusterRoutingList is the whole list of all QdrantClusterRouting objects.
type QdrantClusterRoutingList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of QdrantClusterRouting objects
	Items []QdrantClusterRouting `json:"items"`
}

// QdrantClusterRoutingStatus defines the observed state of QdrantClusterRouting
// +kubebuilder:pruning:PreserveUnknownFields
type QdrantClusterRoutingStatus struct {
	// Set to true if routing of the Qdrant cluster has been bootstrapped once.
	// This implies that at least one route is bootstrapped, for detailed information see the BootstrapInfos field
	Bootstrapped *bool `json:"bootstrapped,omitempty"`
	// Individual bootstrap status info (e.g. when multiple routes are available for this Qdrant cluster)
	BootstrapInfos *[]BootstrapStatusInfo `json:"bootstrapInfos,omitempty"`
	// MultiAZRoutingPhase tracks the cluster's position in a multi-AZ cutover transition.
	// During transient phases the route-manager serves the cluster on BOTH LBs so DNS can flip safely.
	// +kubebuilder:validation:Enum=SettledSingle;SettledMulti;PromotingToMulti;DemotingToSingle
	// +optional
	MultiAZRoutingPhase string `json:"multiAZRoutingPhase,omitempty"`
	// MultiAZRoutingPhaseEnteredAt records when the current phase was entered.
	// +optional
	MultiAZRoutingPhaseEnteredAt *metav1.Time `json:"multiAZRoutingPhaseEnteredAt,omitempty"`
	// MultiAZDNSGateClearedAt is when the per-cluster DNSEndpoint first reached the
	// desired state for the current transient phase (present for Promoting, absent for
	// Demoting). Grace counts from this timestamp — not from MultiAZRoutingPhaseEnteredAt —
	// so slow DNSEndpoint convergence does not cause premature settling. Reset to nil
	// when the gate condition becomes unsatisfied (DNS flap) or on phase change.
	// +optional
	MultiAZDNSGateClearedAt *metav1.Time `json:"multiAZDNSGateClearedAt,omitempty"`
}

// MultiAZRoutingPhase values written to QdrantClusterRoutingStatus.MultiAZRoutingPhase.
const (
	MultiAZRoutingPhaseSettledSingle    = "SettledSingle"
	MultiAZRoutingPhaseSettledMulti     = "SettledMulti"
	MultiAZRoutingPhasePromotingToMulti = "PromotingToMulti"
	MultiAZRoutingPhaseDemotingToSingle = "DemotingToSingle"
)

// BootstrapStatusInfo is part of QdrantClusterRoutingStatus.
type BootstrapStatusInfo struct {
	// Identifier of the route this bootstrap status info belongs to.
	RouteId string `json:"routeId,omitempty"`
	// Set if the route uses a shared loadbalancer.
	Shared *bool `json:"shared,omitempty"`
	// Set if the route uses a dedicated loadbalancer.
	Dedicated *bool `json:"dedicated,omitempty"`
	// Set to true if routing of the Qdrant cluster has been bootstrapped once for this specific route.
	Bootstrapped *bool `json:"bootstrapped,omitempty"`
}

func init() {
	SchemeBuilder.Register(&QdrantClusterRouting{}, &QdrantClusterRoutingList{})
}
