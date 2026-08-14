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

// GetSpec returns the routing spec, or the zero spec if the object is nil.
// Returning a value (rather than a pointer) keeps the spec getters below usable
// on the result, so callers can chain safely: qcrt.GetSpec().GetEnabled().
func (r *QdrantClusterRouting) GetSpec() QdrantClusterRoutingSpec {
	if r == nil {
		return QdrantClusterRoutingSpec{}
	}
	return r.Spec
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

// GetClusterId returns the unique identifier of the cluster.
func (s QdrantClusterRoutingSpec) GetClusterId() string {
	return s.ClusterId
}

// GetFQDN returns the fully qualified domain name of the cluster.
func (s QdrantClusterRoutingSpec) GetFQDN() string {
	return s.FQDN
}

// GetEnabled returns whether ingress is enabled for the cluster.
// Unset means enabled, matching the CRD default.
func (s QdrantClusterRoutingSpec) GetEnabled() bool {
	if s.Enabled == nil {
		return true
	}
	return *s.Enabled
}

// GetShared returns whether the cluster uses at least one shared loadbalancer.
func (s QdrantClusterRoutingSpec) GetShared() bool {
	if s.Shared == nil {
		return false
	}
	return *s.Shared
}

// GetDedicated returns whether the cluster uses at least one dedicated loadbalancer.
func (s QdrantClusterRoutingSpec) GetDedicated() bool {
	if s.Dedicated == nil {
		return false
	}
	return *s.Dedicated
}

// GetTLS returns whether tls is enabled at qdrant level.
func (s QdrantClusterRoutingSpec) GetTLS() bool {
	if s.TLS == nil {
		return false
	}
	return *s.TLS
}

// GetServicePerNode returns whether each node has a dedicated route.
// Unset means enabled, matching the CRD default.
func (s QdrantClusterRoutingSpec) GetServicePerNode() bool {
	if s.ServicePerNode == nil {
		return true
	}
	return *s.ServicePerNode
}

// GetNodeIndexes returns the indexes of the individual nodes in the cluster.
func (s QdrantClusterRoutingSpec) GetNodeIndexes() []int {
	return s.NodeIndexes
}

// GetAllowedSourceRanges returns the allowed CIDR source ranges for the ingress.
// An empty result means no restriction, so callers that enforce it must treat
// "no ranges" as "allow all" deliberately rather than by omission.
func (s QdrantClusterRoutingSpec) GetAllowedSourceRanges() []string {
	return s.AllowedSourceRanges
}

// GetEnableAccessLog returns whether the (proxy) access log is enabled.
func (s QdrantClusterRoutingSpec) GetEnableAccessLog() bool {
	if s.EnableAccessLog == nil {
		return false
	}
	return *s.EnableAccessLog
}

// GetMultiAZ returns whether the cluster spans multiple availability zones.
func (s QdrantClusterRoutingSpec) GetMultiAZ() bool {
	return s.MultiAZ
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
