package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KindEnvoyBootstrapConfig     = "EnvoyBootstrapConfig"
	ResourceEnvoyBootstrapConfig = "envoybootstrapconfigs"
)

//+kubebuilder:object:root=true
// +kubebuilder:resource:path=envoybootstrapconfigs,singular=envoybootstrapconfig,shortName=ebc;ebcs
// +kubebuilder:printcolumn:name="ClusterID",type=string,JSONPath=`.spec.clusterID`
// +kubebuilder:printcolumn:name="Ready",type=boolean,JSONPath=`.status.ready`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:subresource:status

// EnvoyBootstrapConfig is the Schema for the envoybootstrapconfigs API.
// It declares the desired state for an Envoy bootstrap configuration.
// The route-manager reconciler watches these resources and produces a
// ConfigMap (rendered bootstrap JSON) and a Secret (access token) in the
// same namespace.
type EnvoyBootstrapConfig struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EnvoyBootstrapConfigSpec   `json:"spec,omitempty"`
	Status EnvoyBootstrapConfigStatus `json:"status,omitempty"`
}

// EnvoyBootstrapConfigSpec describes the desired Envoy bootstrap configuration.
type EnvoyBootstrapConfigSpec struct {
	// ClusterID identifies the Qdrant cluster this Envoy instance serves.
	// When set the Envoy runs in dedicated mode for this cluster.
	// When nil it runs in shared mode.
	// +optional
	ClusterID *string `json:"clusterID,omitempty"`
	// ProxyProtocolEnabled enables the PROXY protocol on Envoy listeners.
	// +kubebuilder:default=false
	// +optional
	ProxyProtocolEnabled bool `json:"proxyProtocolEnabled,omitempty"`
}

// EnvoyBootstrapConfigStatus defines the observed state of EnvoyBootstrapConfig.
type EnvoyBootstrapConfigStatus struct {
	// Ready is true once the ConfigMap and Secret exist and are up-to-date.
	Ready bool `json:"ready,omitempty"`
	// ConfigMapName is the name of the generated ConfigMap containing the bootstrap JSON.
	ConfigMapName string `json:"configMapName,omitempty"`
	// SecretName is the name of the access-token Secret (created or reused).
	SecretName string `json:"secretName,omitempty"`
	// Conditions represent the latest available observations of the resource's state.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EnvoyBootstrapConfigList contains a list of EnvoyBootstrapConfig objects.
type EnvoyBootstrapConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of EnvoyBootstrapConfig objects
	Items []EnvoyBootstrapConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EnvoyBootstrapConfig{}, &EnvoyBootstrapConfigList{})
}
