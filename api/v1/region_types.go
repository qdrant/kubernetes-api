package v1

import (
	helmapi "github.com/fluxcd/helm-controller/api/v2beta2"
	srcapi "github.com/fluxcd/source-controller/api/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KindQdrantCloudRegion     = "QdrantCloudRegion"
	ResourceQdrantCloudRegion = "qdrantcloudregions"
)

// QdrantCloudRegionSpec defines the desired state of QdrantCloudRegion
type QdrantCloudRegionSpec struct {
	// Id specifies the unique identifier of the region
	Id string `json:"id,omitempty"`
	// HelmRepositories specifies the list of helm repositories to be created to the region
	// +optional
	HelmRepositories []HelmRepository `json:"helmRepositories,omitempty"`
	// HelmReleases specifies the list of helm releases to be created to the region
	// +optional
	HelmReleases []HelmRelease `json:"helmReleases,omitempty"`
}

type HelmRepository struct {
	// MarkedForDeletionAt specifies the time when the helm repository was marked for deletion
	// +optional
	MarkedForDeletionAt *string `json:"markedForDeletionAt,omitempty"`
	// Object specifies the helm repository object
	// +kubebuilder:validation:EmbeddedResource
	Object *srcapi.HelmRepository `json:"object,omitempty"`
}

type HelmRelease struct {
	// MarkedForDeletionAt specifies the time when the helm release was marked for deletion
	// +optional
	MarkedForDeletionAt *string `json:"markedForDeletionAt,omitempty"`
	// Object specifies the helm release object
	// +kubebuilder:validation:EmbeddedResource
	Object *helmapi.HelmRelease `json:"object,omitempty"`
}

type RegionPhase string
type MonitoringSource string

const (
	RegionPhaseReady        RegionPhase      = "Ready"
	RegionPhaseNotReady     RegionPhase      = "NotReady"
	FailedToSync            RegionPhase      = "FailedToSync"
	KubeletMonitoringSource MonitoringSource = "kubelet"
)

type QdrantCloudRegionStatus struct {
	// Phase specifies the current phase of the region
	// +optional
	Phase RegionPhase `json:"phase,omitempty"`
	// KubernetesVersion specifies the version of the Kubernetes cluster
	// +optional
	KubernetesVersion string `json:"k8sVersion,omitempty"`
	// NumberOfNodes specifies the number of nodes in the Kubernetes cluster
	// +optional
	NumberOfNodes int `json:"numberOfNodes,omitempty"`
	// Capabilities specifies the capabilities of the Kubernetes cluster
	// +optional
	Capabilities *RegionCapabilities `json:"capabilities,omitempty"`
	// HelmRepositories specifies the status of the helm repositories
	// +optional
	HelmRepositories []ComponentStatus `json:"helmRepositories,omitempty"`
	// HelmReleases specifies the status of the helm releases
	// +optional
	HelmReleases []ComponentStatus `json:"helmReleases,omitempty"`
	// Message specifies the info explaining the current phase of the region
	// +optional
	Message string `json:"message,omitempty"`
	// KubernetesDistribution specifies the distribution of the Kubernetes cluster
	// +optional
	KubernetesDistribution KubernetesDistribution `json:"k8sDistribution,omitempty"`
	// AlternativeMonitoringSource specifies the alternative monitoring source
	// +optional
	AlternativeMonitoringSource AlternativeMonitoringSource `json:"alternativeMonitoringSource,omitempty"`
}

type AlternativeMonitoringSource struct {
	// CAdvisorMonitoringSource specifies the cAdvisor monitoring source
	// +optional
	CAdvisorMonitoringSource MonitoringSource `json:"cAdvisorMonitoringSource,omitempty"`
	// NodeMonitoringSource specifies the node monitoring source
	// +optional
	NodeMonitoringSource MonitoringSource `json:"nodeMonitoringSource,omitempty"`
}

type RegionCapabilities struct {
	// VolumeSnapshot specifies whether the Kubernetes cluster supports volume snapshot
	// +optional
	VolumeSnapshot *bool `json:"volumeSnapshot,omitempty"`
	// VolumeExpansion specifies whether the Kubernetes cluster supports volume expansion
	// +optional
	VolumeExpansion *bool `json:"volumeExpansion,omitempty"`
}

type KubernetesDistribution string

const (
	K8sDistributionUnknown   KubernetesDistribution = "unknown"
	K8sDistributionAWS       KubernetesDistribution = "aws"
	K8sDistributionGCP       KubernetesDistribution = "gcp"
	K8sDistributionAzure     KubernetesDistribution = "azure"
	K8sDistributionDO        KubernetesDistribution = "do"
	K8sDistributionScaleway  KubernetesDistribution = "scaleway"
	K8sDistributionOpenShift KubernetesDistribution = "openshift"
	K8sDistributionLinode    KubernetesDistribution = "linode"
	K8sDistributionCivo      KubernetesDistribution = "civo"
	K8sDistributionOCI       KubernetesDistribution = "oci"
	K8sDistributionOVHCloud  KubernetesDistribution = "ovhcloud"
	K8sDistributionStackit   KubernetesDistribution = "stackit"
	K8sDistributionVultr     KubernetesDistribution = "vultr"
	K8sDistributionK3s       KubernetesDistribution = "k3s"
)

type ComponentPhase string

const (
	ComponentPhaseReady    ComponentPhase = "Ready"
	ComponentPhaseNotReady ComponentPhase = "NotReady"
	ComponentPhaseUnknown  ComponentPhase = "Unknown"
	ComponentPhaseNotFound ComponentPhase = "NotFound"
)

type ComponentStatus struct {
	// Name specifies the name of the component
	Name string `json:"name"`
	// Namespace specifies the namespace of the component
	Namespace string `json:"namespace,omitempty"`
	// Version specifies the version of the component
	// +optional
	Version string `json:"version,omitempty"`
	// Phase specifies the current phase of the component
	Phase ComponentPhase `json:"phase,omitempty"`
	// Message specifies the info explaining the current phase of the component
	// +optional
	Message string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:pruning:PreserveUnknownFields
// +kubebuilder:resource:path=qdrantcloudregions,scope=Cluster,singular=qdrantcloudregion,shortName=region;regions
// +kubebuilder:printcolumn:name="K8s Version",type=string,JSONPath=`.status.k8sVersion`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// QdrantCloudRegion is the Schema for the qdrantcloudregions API
type QdrantCloudRegion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QdrantCloudRegionSpec   `json:"spec,omitempty"`
	Status QdrantCloudRegionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// QdrantCloudRegionList contains a list of QdrantCloudRegion
type QdrantCloudRegionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QdrantCloudRegion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QdrantCloudRegion{}, &QdrantCloudRegionList{})
}
