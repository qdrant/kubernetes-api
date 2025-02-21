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
	// Components specifies the list of components to be installed in the region
	// +optional
	Components []ComponentReference `json:"components,omitempty"`
	// HelmRepositories specifies the list of helm repositories to be created to the region
	// Deprecated: Use "Components" instead
	// +optional
	HelmRepositories []HelmRepository `json:"helmRepositories,omitempty"`
	// HelmReleases specifies the list of helm releases to be created to the region
	// Deprecated: Use "Components" instead
	// +optional
	HelmReleases []HelmRelease `json:"helmReleases,omitempty"`
}

type ComponentReference struct {
	// APIVersion is the group and version of the component being referenced.
	APIVersion string `json:"apiVersion"`
	// Kind is the type of component being referenced
	Kind string `json:"kind"`
	// Name is the name of component being referenced
	Name string `json:"name"`
	// Namespace is the namespace of component being referenced.
	Namespace string `json:"namespace"`
	// MarkedForDeletion specifies whether the component is marked for deletion
	// +optional
	MarkedForDeletion bool `json:"markedForDeletion,omitempty"`
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

const (
	RegionPhaseReady    RegionPhase = "Ready"
	RegionPhaseNotReady RegionPhase = "NotReady"
	FailedToSync        RegionPhase = "FailedToSync"
)

type MetricSource string

const (
	KubeletMetricSource MetricSource = "kubelet"
	ApiMetricSource     MetricSource = "api"
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
	// Monitoring specifies monitoring status
	// +optional
	Monitoring Monitoring `json:"monitoring,omitempty"`
	// StorageClasses contains the availble StorageClasses in the Kubernetes cluster
	// +optional
	StorageClasses []StorageClass `json:"storageClasses,omitempty"`
	// VolumeSnapshotClasses contains the available VolumeSnapshotClasses in the Kubernetes cluster
	// +optional
	VolumeSnapshotClasses []VolumeSnapshotClass `json:"volumeSnapshotClasses,omitempty"`
	// NodeInfos contains the information about the nodes in the Kubernetes cluster
	// +optional
	NodeInfos []NodeInfo `json:"nodeInfos,omitempty"`
}

type StorageClass struct {
	// Name specifies the name of the storage class
	Name string `json:"name"`
	// Default specifies whether the storage class is the default storage class
	Default bool `json:"default"`
	// Provisioner specifies the provisioner of the storage class
	Provisioner string `json:"provisioner"`
	// AllowVolumeExpansion specifies whether the storage class allows volume expansion
	AllowVolumeExpansion bool `json:"allowVolumeExpansion"`
	// ReclaimPolicy specifies the reclaim policy of the storage class
	ReclaimPolicy string `json:"reclaimPolicy"`
	// Parameters specifies the parameters of the storage class
	// +optional
	Parameters map[string]string `json:"parameters"`
}

type VolumeSnapshotClass struct {
	// Name specifies the name of the volume snapshot class
	Name string `json:"name"`
	// Driver specifies the driver of the volume snapshot class
	Driver string `json:"driver"`
}

type NodeInfo struct {
	// Name specifies the name of the node
	Name string `json:"name"`
	// Region specifies the region of the node
	// +optional
	Region string `json:"region,omitempty"`
	// Zone specifies the zone of the node
	// +optional
	Zone string `json:"zone,omitempty"`
	// InstanceType specifies the instance type of the node
	// +optional
	InstanceType string `json:"instanceType,omitempty"`
	// Arch specifies the CPU architecture of the node
	// +optional
	Arch string `json:"arch,omitempty"`
	// Capacity specifies the capacity of the node
	Capacity NodeResourceInfo `json:"capacity"`
	// Allocatable specifies the allocatable resources of the node
	Allocatable NodeResourceInfo `json:"allocatable"`
}

type NodeResourceInfo struct {
	// CPU specifies the CPU resources of the node
	CPU string `json:"cpu"`
	// Memory specifies the memory resources of the node
	Memory string `json:"memory"`
	// Pods specifies the pods resources of the node
	Pods string `json:"pods"`
	// EphemeralStorage specifies the ephemeral storage resources of the node
	EphemeralStorage string `json:"ephemeralStorage"`
}

type Monitoring struct {
	// CAdvisorMetricSource specifies the cAdvisor metric source
	// +optional
	CAdvisorMetricSource MetricSource `json:"cAdvisorMetricSource,omitempty"`
	// NodeMetricSource specifies the node metric source
	// +optional
	NodeMetricSource MetricSource `json:"nodeMetricSource,omitempty"`
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
