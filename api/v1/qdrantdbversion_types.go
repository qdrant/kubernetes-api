package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// QdrantVersionSpec defines the desired state of QdrantVersion
type QdrantVersionSpec struct {
	// Version number (should be semver compliant).
	// E.g. "v1.10.1"
	Version string `json:"version,omitempty"`
	// If set, this version is default for new clusters on Cloud.
	// There should be only 1 Qdrant version in the platform set as default.
	// +kubebuilder:default=false
	// +optional
	IsDefault bool `json:"isDefault,omitempty"`
	// Full docker image to use for this version.
	// If empty, a default image will be derived from Version (and qdrant/qdrant is assumed).
	// +optional
	Image string `json:"image,omitempty"`
	// If set, this version cannot be used for new clusters.
	// +kubebuilder:default=false
	// +optional
	Unavailable bool `json:"unavailable,omitempty"`
	// If set, this version is no longer actively supported.
	// +kubebuilder:default=false
	// +optional
	IsEndOfLife bool `json:"isEndOfLife,omitempty"`
	// If set, this version can only be used by accounts with given IDs.
	// +optional
	AccountIDs []string `json:"accountIds,omitempty"`
	// If set, this version can only be used by accounts that have been given the listed privileges.
	// +optional
	AccountPrivileges []string `json:"accountPrivileges,omitempty"`
	// General remarks for human reading
	// +optional
	Remarks string `json:"remarks,omitempty"`
	// Release Notes URL for the specified version
	ReleaseNotesURL string `json:"releaseNotesURL,omitempty"`
}

//+kubebuilder:object:root=true
// +kubebuilder:resource:path=qdrantversions,singular=qdrantversions,shortName=qv;qvs
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="IsDefault",type=boolean,JSONPath=`.spec.isDefault`
// +kubebuilder:printcolumn:name="Unavailable",type=boolean,JSONPath=`.spec.unavailable`
// +kubebuilder:printcolumn:name="IsEndOfLife",type=boolean,JSONPath=`.metadata.isEndOfLife`

// QdrantVersion is the Schema for the qdrantversions API
type QdrantVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec QdrantVersionSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// QdrantVersionList contains a list of QdrantVersion
type QdrantVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QdrantVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QdrantVersion{}, &QdrantVersionList{})
}
