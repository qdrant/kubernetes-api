package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// QdrantReleaseSpec defines the desired state of QdrantRelease
type QdrantReleaseSpec struct {
	// Version number (should be semver compliant).
	// E.g. "v1.10.1"
	Version string `json:"version,omitempty"`
	// If set, this version is default for new clusters on Cloud.
	// There should be only 1 Qdrant version in the platform set as default.
	// +kubebuilder:default=false
	// +optional
	Default bool `json:"default,omitempty"`
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
	EndOfLife bool `json:"endOfLife,omitempty"`
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

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=qdrantreleases,singular=qdrantrelease,shortName=qr;qrs
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Default",type=boolean,JSONPath=`.spec.default`
// +kubebuilder:printcolumn:name="Unavailable",type=boolean,JSONPath=`.spec.unavailable`
// +kubebuilder:printcolumn:name="EndOfLife",type=boolean,JSONPath=`.spec.endOfLife`

// QdrantRelease describes an available Qdrant release
type QdrantRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec QdrantReleaseSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// QdrantReleaseList contains a list of QdrantRelease
type QdrantReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QdrantRelease `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QdrantRelease{}, &QdrantReleaseList{})
}
