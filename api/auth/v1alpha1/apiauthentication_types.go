package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIAuthentication is a configuration for authenticating against Qdrant clusters.
type APIAuthentication struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec APIAuthenticationSpec `json:"spec"`
}

// APIAuthenticationSpec describes the configuration for authenticating against Qdrant clusters.
type APIAuthenticationSpec struct {
	// +kubebuilder:validation:MinLength=128
	// +kubebuilder:validation:MaxLength=128
	// +optional
	// SHA512 hash of an API key.
	SHA512 *string `json:"sha512,omitempty"`

	// +listType=set
	// List of cluster IDs for which the API key is valid
	ClusterIDs []string `json:"clusterIDs"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIAuthenticationList is the whole list of all APIAuthentication objects.
type APIAuthenticationList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of APIAuthentication objects
	Items []APIAuthentication `json:"items"`
}
