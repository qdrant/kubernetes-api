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
	// For shared routing this Id will be used for SNI resolving.
	ClusterId string `json:"clusterId"`
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
	// If true enabled envoy access log
	// +optional
	EnableEnvoyAccessLog *bool `json:"enableEnvoyAccessLog,omitempty"`
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
}

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
