package v1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//goland:noinspection GoUnusedConst
const (
	KindQdrantCluster     = "QdrantCluster"
	ResourceQdrantCluster = "qdrantclusters"
)

//goland:noinspection GoUnusedConst
const (
	// RestartedAtAnnotationKey is the annotation key to trigger a restart.
	// The annotation should be placed on the QdrantCluster instance.
	// The value should be a [RFC3339 formatted] date.
	// If the value is updated it will retrigger the restart.
	// For historical reasons the key doesn't start with `operator.qdrant.com/`
	RestartedAtAnnotationKey = "restartedAt"
	// RecreateNodeAnnotationKey is the annotation key to recreate a certain node.
	// The annotation should be placed on the pod created by the operator (for the node that need to be recreated).
	// It is allowed to add this annotation to multiple pods, the operator will handle them all.
	// The value should be non-empty, but is free to use, and will be used in Events.
	// This feature requires that the cluster-manager is enabled.
	RecreateNodeAnnotationKey = "operator.qdrant.com/recreate-node"
	// ReinitAnnotationKey is the annotation key to trigger reinitialization of the given cluster.
	// The annotation value is ignored and can be used to document why reinitialization is requested.
	ReinitAnnotationKey = "operator.qdrant.com/reinit"
)

// GPUType specifies the type of GPU to use.
// +kubebuilder:validation:Enum=nvidia;amd
type GPUType string

//goland:noinspection GoUnusedConst
const (
	GPUTypeNvidia GPUType = "nvidia"
	GPUTypeAmd    GPUType = "amd"
)

// RebalanceStrategy specifies the strategy to use for automaticially rebalancing shards the cluster.
// +kubebuilder:validation:Enum=by_count;by_size;by_count_and_size
type RebalanceStrategy string

//goland:noinspection GoUnusedConst
const (
	ByCount        RebalanceStrategy = "by_count"
	BySize         RebalanceStrategy = "by_size"
	ByCountAndSize RebalanceStrategy = "by_count_and_size"
)

// QdrantClusterSpec defines the desired state of QdrantCluster
// +kubebuilder:pruning:PreserveUnknownFields
type QdrantClusterSpec struct {
	// Id specifies the unique identifier of the cluster
	Id string `json:"id"`
	// Version specifies the version of Qdrant to deploy
	Version string `json:"version"`
	// Size specifies the desired number of Qdrant nodes in the cluster
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	Size int `json:"size"`
	// ServicePerNode specifies whether the cluster should start a dedicated service for each node.
	// +kubebuilder:default=true
	// +optional
	ServicePerNode *bool `json:"servicePerNode,omitempty"`
	// ClusterManager specifies whether to use the cluster manager for this cluster.
	// The Python-operator will deploy a dedicated cluster manager instance.
	// The Go-operator will use a shared instance.
	// If not set, the default will be taken from the operator config.
	// +optional
	ClusterManager *bool `json:"clusterManager,omitempty"`
	// Suspend specifies whether to suspend the cluster.
	// If enabled, the cluster will be suspended and all related resources will be removed except the PVCs.
	// +kubebuilder:default=false
	// +optional
	Suspend bool `json:"suspend,omitempty"`
	// Pauses specifies a list of pause request by developer for manual maintenance.
	// Operator will skip handling any changes in the CR if any pause request is present.
	// +optional
	Pauses []Pause `json:"pauses,omitempty"`
	// Image specifies the image to use for each Qdrant node.
	// +optional
	Image *QdrantImage `json:"image,omitempty"`
	// Resources specifies the resources to allocate for each Qdrant node.
	Resources Resources `json:"resources,omitempty"`
	// Security specifies the security context for each Qdrant node.
	// +optional
	Security *QdrantSecurityContext `json:"security,omitempty"`
	// Tolerations specifies the tolerations for each Qdrant node.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// NodeSelector specifies the node selector for each Qdrant node.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Config specifies the Qdrant configuration setttings for the clusters.
	// +optional
	Config *QdrantConfiguration `json:"config,omitempty"`
	// Ingress specifies the ingress for the cluster.
	// +optional
	Ingress *Ingress `json:"ingress,omitempty"`
	// Service specifies the configuration of the Qdrant Kubernetes Service.
	// +optional
	Service *KubernetesService `json:"service,omitempty"`
	// GPU specifies GPU configuration for the cluster. If this field is not set, no GPU will be used.
	// +optional
	GPU *GPU `json:"gpu,omitempty"`
	// StatefulSet specifies the configuration of the Qdrant Kubernetes StatefulSet.
	// +optional
	StatefulSet *KubernetesStatefulSet `json:"statefulSet,omitempty"`
	// StorageClassNames specifies the storage class names for db and snapshots.
	// +optional
	StorageClassNames *StorageClassNames `json:"storageClassNames,omitempty"`
	// TopologySpreadConstraints specifies the topology spread constraints for the cluster.
	// +optional
	TopologySpreadConstraints *[]corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	// PodDisruptionBudget specifies the pod disruption budget for the cluster.
	// +optional
	PodDisruptionBudget *policyv1.PodDisruptionBudgetSpec `json:"podDisruptionBudget,omitempty"`
	// RestartAllPodsConcurrently specifies whether to restart all pods concurrently (also called one-shot-restart).
	// If enabled, all the pods in the cluster will be restarted concurrently in situations where multiple pods
	// need to be restarted, like when RestartedAtAnnotationKey is added/updated or the Qdrant version needs to be upgraded.
	// This helps sharded but not replicated clusters to reduce downtime to a possible minimum during restart.
	// If unset, the operator is going to restart nodes concurrently if none of the collections if replicated.
	// +optional
	RestartAllPodsConcurrently *bool `json:"restartAllPodsConcurrently,omitempty"`
	// If StartupDelaySeconds is set (> 0), an additional 'sleep <value>' will be emitted to the pod startup.
	// The sleep will be added when a pod is restarted, it will not force any pod to restart.
	// This feature can be used for debugging the core, e.g. if a pod is in crash loop, it provided a way
	// to inspect the attached storage.
	// +optional
	StartupDelaySeconds *int `json:"startupDelaySeconds,omitempty"`
	// RebalanceStrategy specifies the strategy to use for automaticially rebalancing shards the cluster.
	// Cluster-manager needs to be enabled for this feature to work.
	// +optional
	RebalanceStrategy *RebalanceStrategy `json:"rebalanceStrategy,omitempty"`
}

// Validate if there are incorrect settings in the CRD
func (s QdrantClusterSpec) Validate() error {
	if err := s.Resources.Validate("Spec.Resources"); err != nil {
		return err
	}
	return nil
}

// GetServicePerNode get the service per node, taking the default (true) into concideration
func (s QdrantClusterSpec) GetServicePerNode() bool {
	if s.ServicePerNode == nil {
		return true
	}
	return *s.ServicePerNode
}

type GPU struct {
	// GPUType specifies the type of the GPU to use. If set, GPU indexing is enabled.
	// +kubebuilder:validation:Enum=nvidia;amd
	GPUType GPUType `json:"gpuType"`
	// ForceHalfPrecision for `f32` values while indexing.
	// `f16` conversion will take place
	// only inside GPU memory and won't affect storage type.
	// +kubebuilder:default=false
	ForceHalfPrecision bool `json:"forceHalfPrecision"`
	// DeviceFilter for GPU devices by hardware name. Case-insensitive.
	// List of substrings to match against the gpu device name.
	// Example: [- "nvidia"]
	// If not specified, all devices are accepted.
	// +kubebuilder:validation:MinItems:=1
	// +optional
	DeviceFilter []string `json:"deviceFilter,omitempty"`
	// Devices is a List of explicit GPU devices to use.
	// If host has multiple GPUs, this option allows to select specific devices
	// by their index in the list of found devices.
	// If `deviceFilter` is set, indexes are applied after filtering.
	// If not specified, all devices are accepted.
	// +kubebuilder:validation:MinItems:=1
	// +optional
	Devices []string `json:"devices,omitempty"`
	// ParallelIndexes is the number of parallel indexes to run on the GPU.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum:=1
	ParallelIndexes int `json:"parallelIndexes"`
	// GroupsCount is the amount of used vulkan "groups" of GPU.
	// In other words, how many parallel points can be indexed by GPU.
	// Optimal value might depend on the GPU model.
	// Proportional, but doesn't necessary equal to the physical number of warps.
	// Do not change this value unless you know what you are doing.
	// +optional
	// +kubebuilder:validation:Minimum:=1
	GroupsCount int `json:"groupsCount,omitempty"`
	// AllowIntegrated specifies whether to allow integrated GPUs to be used.
	// +kubebuilder:default=false
	AllowIntegrated bool `json:"allowIntegrated"`
}

func (g *GPU) GetGPUType() GPUType {
	if g == nil {
		return ""
	}
	return g.GPUType
}

type KubernetesService struct {
	// Type specifies the type of the Service: "ClusterIP", "NodePort", "LoadBalancer".
	// +kubebuilder:default="ClusterIP"
	// +optional
	Type corev1.ServiceType `json:"type,omitempty"`
	// Annotations specifies the annotations for the Service.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

func (s *KubernetesService) GetType() corev1.ServiceType {
	if s == nil {
		return corev1.ServiceTypeClusterIP
	}
	return s.Type
}

func (s *KubernetesService) GetAnnotations() map[string]string {
	if s == nil {
		return nil
	}
	return s.Annotations
}

type KubernetesStatefulSet struct {
	// Annotations specifies the annotations for the StatefulSet.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Pods  specifies the configuration of the Pods of the Qdrant StatefulSet.
	// +optional
	Pods *KubernetesPod `json:"pods,omitempty"`
}

func (kss *KubernetesStatefulSet) GetPods() *KubernetesPod {
	if kss == nil {
		return nil
	}
	return kss.Pods
}

type KubernetesPod struct {
	// Annotations specifies the annotations for the Pods.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels specifies the labels for the Pods.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// ExtraEnv specifies the extra environment variables for the Pods.
	// +optional
	ExtraEnv []corev1.EnvVar `json:"extraEnv,omitempty"`
}

func (kp *KubernetesPod) GetAnnotations() map[string]string {
	if kp == nil {
		return nil
	}
	return kp.Annotations
}

func (kp *KubernetesPod) GetLabels() map[string]string {
	if kp == nil {
		return nil
	}
	return kp.Labels
}

func (kp *KubernetesPod) GetExtraEnv() []corev1.EnvVar {
	if kp == nil {
		return nil
	}
	return kp.ExtraEnv
}

type Pause struct {
	// Owner specifies the owner of the pause request.
	Owner string `json:"owner,omitempty"`
	// Reason specifies the reason for the pause request.
	Reason string `json:"reason,omitempty"`
	// CreationTimestamp specifies the time when the pause request was created.
	CreationTimestamp string `json:"creationTimestamp,omitempty"`
}

type QdrantImage struct {
	// Repository specifies the repository of the Qdrant image.
	// If not specified defaults the config of the operator (or qdrant/qdrant if not specified in operator).
	// +optional
	Repository *string `json:"repository,omitempty"`
	// PullPolicy specifies the image pull policy for the Qdrant image.
	// If not specified defaults the config of the operator (or IfNotPresent if not specified in operator).
	// +optional
	PullPolicy *corev1.PullPolicy `json:"pullPolicy,omitempty"`
	// PullSecretName specifies the pull secret for the Qdrant image.
	// +optional
	PullSecretName *string `json:"pullSecretName,omitempty"`
}

func (qi *QdrantImage) GetRepository() *string {
	if qi == nil {
		return nil
	}
	return qi.Repository
}

func (qi *QdrantImage) GetPullPolicy() *corev1.PullPolicy {
	if qi == nil {
		return nil
	}
	return qi.PullPolicy
}

func (qi *QdrantImage) GetImagePullSecrets() *string {
	if qi == nil {
		return nil
	}
	return qi.PullSecretName
}

type Resources struct {
	// CPU specifies the CPU limit for each Qdrant node.
	CPU string `json:"cpu,omitempty"`
	// Memory specifies the memory limit for each Qdrant node.
	Memory string `json:"memory,omitempty"`
	// Storage specifies the storage amount for each Qdrant node.
	Storage string `json:"storage,omitempty"`
	// Requests specifies the resource requests for each Qdrant node.
	// +optional
	Requests ResourceRequests `json:"requests,omitempty"`
}

// Validate if there are incorrect settings in the CRD
func (s Resources) Validate(base string) error {
	if _, err := resource.ParseQuantity(s.CPU); err != nil {
		return fmt.Errorf("%s.CPU error: %w", base, err)
	}
	if _, err := resource.ParseQuantity(s.Memory); err != nil {
		return fmt.Errorf("%s.Memory error: %w", base, err)
	}
	if _, err := resource.ParseQuantity(s.Storage); err != nil {
		return fmt.Errorf("%s.Storage error: %w", base, err)
	}
	if err := s.Requests.Validate(base + ".Request"); err != nil {
		return err
	}
	return nil
}

func (r Resources) GetRequestCPU() string {
	if r.Requests.CPU != "" {
		return r.Requests.CPU
	}
	return r.CPU
}

func (r Resources) GetRequestMemory() string {
	if r.Requests.Memory != "" {
		return r.Requests.Memory
	}
	return r.Memory
}

type ResourceRequests struct {
	// CPU specifies the CPU request for each Qdrant node.
	// +optional
	CPU string `json:"cpu,omitempty"`
	// Memory specifies the memory request for each Qdrant node.
	// +optional
	Memory string `json:"memory,omitempty"`
}

// Validate if there are incorrect settings in the CRD
func (s ResourceRequests) Validate(base string) error {
	if s.CPU != "" {
		if _, err := resource.ParseQuantity(s.CPU); err != nil {
			return fmt.Errorf("%s.CPU error: %w", base, err)
		}
	}
	if s.Memory != "" {
		if _, err := resource.ParseQuantity(s.Memory); err != nil {
			return fmt.Errorf("%s.Memory error: %w", base, err)
		}
	}
	return nil
}

type QdrantSecurityContext struct {
	// User specifies the user to run the Qdrant process as.
	User int64 `json:"user,omitempty"`
	// Group specifies the group to run the Qdrant process as.
	Group int64 `json:"group,omitempty"`
	// FsGroup specifies file system group to run the Qdrant process as.
	// +optional
	FsGroup *int64 `json:"fsGroup,omitempty"`
}

func (c *QdrantSecurityContext) GetUser() *int64 {
	if c == nil {
		return nil
	}
	return &c.User
}

func (c *QdrantSecurityContext) GetGroup() *int64 {
	if c == nil {
		return nil
	}
	return &c.Group
}

func (c *QdrantSecurityContext) GetFsGroup() *int64 {
	if c == nil {
		return nil
	}
	return c.FsGroup
}

type QdrantConfiguration struct {
	// Collection specifies the default collection configuration for Qdrant.
	// +optional
	Collection *QdrantConfigurationCollection `json:"collection,omitempty"`
	// LogLevel specifies the log level for Qdrant.
	// +optional
	LogLevel *string `json:"log_level,omitempty"`
	// Service specifies the service level configuration for Qdrant.
	// +optional
	Service *QdrantConfigurationService `json:"service,omitempty"`
	// TLS specifies the TLS configuration for Qdrant.
	// +optional
	TLS *QdrantConfigurationTLS `json:"tls,omitempty"`
	// Storage specifies the storage configuration for Qdrant.
	// +optional
	Storage *StorageConfig `json:"storage,omitempty"`
	// Inference configuration. This is used in Qdrant Managed Cloud only. If not set Inference is not available to this cluster.
	// +optional
	Inference *InferenceConfig `json:"inference,omitempty"`
}

type InferenceConfig struct {
	// Enabled specifies whether to enable inference for the cluster or not.
	// +kubebuilder:default=false
	// +optional
	Enabled bool `json:"enabled"`
}

type StorageConfig struct {
	// Performance configuration
	// +optional
	Performance *StoragePerformanceConfig `json:"performance,omitempty"`
	// MaxCollections represents the maximal number of collections allowed to be created.
	// It can be set for Qdrant version >= 1.14.1
	// Default to 1000 if omitted and Qdrant version >= 1.15.0
	// +optional
	// +kubebuilder:validation:Minimum:=1
	MaxCollections *uint `json:"maxCollections,omitempty"`
}

type StoragePerformanceConfig struct {
	// OptimizerCPUBudget defines the number of CPU allocation.
	// If 0 - auto selection, keep 1 or more CPUs unallocated depending on CPU size
	// If negative - subtract this number of CPUs from the available CPUs.
	// If positive - use this exact number of CPUs.
	// +optional
	OptimizerCPUBudget *int64 `json:"optimizer_cpu_budget,omitempty"`
	// AsyncScorer enables io_uring when rescoring
	// +optional
	AsyncScorer *bool `json:"async_scorer,omitempty"`
}

func (c *QdrantConfiguration) GetService() *QdrantConfigurationService {
	if c == nil {
		return nil
	}
	return c.Service
}

func (c *QdrantConfiguration) GetTLS() *QdrantConfigurationTLS {
	if c == nil {
		return nil
	}
	return c.TLS
}

type QdrantConfigurationCollection struct {
	// ReplicationFactor specifies the default number of replicas of each shard
	// +optional
	ReplicationFactor *int64 `json:"replication_factor,omitempty"`
	// WriteConsistencyFactor specifies how many replicas should apply the operation to consider it successful
	// +optional
	WriteConsistencyFactor *int64 `json:"write_consistency_factor,omitempty"`
	// Vectors specifies the default parameters for vectors
	// +optional
	Vectors *QdrantConfigurationCollectionVectors `json:"vectors,omitempty"`
	// StrictMode specifies the strict mode configuration for the collection
	// +optional
	StrictMode *QdrantConfigurationCollectionStrictMode `json:"strict_mode,omitempty"`
}

type QdrantConfigurationCollectionStrictMode struct {
	// MaxPayloadIndexCount represents the maximal number of payload indexes allowed to be created.
	// It can be set for Qdrant version >= 1.16.0
	// Default to 100 if omitted and Qdrant version >= 1.16.0
	// +optional
	// +kubebuilder:validation:Minimum:=1
	MaxPayloadIndexCount *uint `json:"max_payload_index_count,omitempty"`
}

type QdrantConfigurationCollectionVectors struct {
	// OnDisk specifies whether vectors should be stored in memory or on disk.
	// +optional
	OnDisk *bool `json:"on_disk,omitempty"`
}

type QdrantConfigurationService struct {
	// ApiKey for the qdrant instance
	// +optional
	ApiKey *QdrantSecretKeyRef `json:"api_key,omitempty"`
	// ReadOnlyApiKey for the qdrant instance
	// +optional
	ReadOnlyApiKey *QdrantSecretKeyRef `json:"read_only_api_key,omitempty"`
	// JwtRbac specifies whether to enable jwt rbac for the qdrant instance
	// Default is false
	// +optional
	JwtRbac *bool `json:"jwt_rbac,omitempty"`
	// HideJwtDashboard specifies whether to hide the JWT dashboard of the embedded UI
	// Default is false
	// +optional
	HideJwtDashboard *bool `json:"hide_jwt_dashboard,omitempty"`
	// EnableTLS specifies whether to enable tls for the qdrant instance
	// Default is false
	// +optional
	EnableTLS *bool `json:"enable_tls,omitempty"`
	// MaxRequestSizeMb specifies them maximum size of POST data in a single request in megabytes
	// Default, if not set is 32 (MB)
	// +optional
	MaxRequestSizeMb *int64 `json:"max_request_size_mb,omitempty"`
}

func (c *QdrantConfigurationService) GetApiKey() *QdrantSecretKeyRef {
	if c == nil {
		return nil
	}
	return c.ApiKey
}

func (c *QdrantConfigurationService) GetReadOnlyApiKey() *QdrantSecretKeyRef {
	if c == nil {
		return nil
	}
	return c.ReadOnlyApiKey
}

func (c *QdrantConfigurationService) GetJwtRbac() bool {
	if c == nil || c.JwtRbac == nil {
		return false
	}
	return *c.JwtRbac
}

func (c *QdrantConfigurationService) GetHideJwtDashboard() bool {
	if c == nil || c.HideJwtDashboard == nil {
		return false
	}
	return *c.HideJwtDashboard
}

func (c *QdrantConfigurationService) GetEnableTLS() bool {
	if c == nil || c.EnableTLS == nil {
		return false
	}
	return *c.EnableTLS
}

func (c *QdrantConfigurationService) GetMaxRequestSizeMb() int64 {
	if c == nil || c.MaxRequestSizeMb == nil {
		return 32
	}
	return *c.MaxRequestSizeMb
}

type QdrantSecretKeyRef struct {
	// SecretKeyRef to the secret containing data to configure the qdrant instance
	// +optional
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

func (r *QdrantSecretKeyRef) GetQdrantSecretKeyRef() *corev1.SecretKeySelector {
	if r == nil {
		return nil
	}
	return r.SecretKeyRef
}

type QdrantConfigurationTLS struct {
	// Reference to the secret containing the server certificate chain file
	// +optional
	Cert *QdrantSecretKeyRef `json:"cert,omitempty"`
	// Reference to the secret containing the server private key file
	// +optional
	Key *QdrantSecretKeyRef `json:"key,omitempty"`
	// Reference to the secret containing the CA certificate file
	// +optional
	CaCert *QdrantSecretKeyRef `json:"caCert,omitempty"`
}

func (c *QdrantConfigurationTLS) GetCert() *QdrantSecretKeyRef {
	if c == nil {
		return nil
	}
	return c.Cert
}

func (c *QdrantConfigurationTLS) GetCaCert() *QdrantSecretKeyRef {
	if c == nil {
		return nil
	}
	return c.CaCert
}

func (c *QdrantConfigurationTLS) GetKey() *QdrantSecretKeyRef {
	if c == nil {
		return nil
	}
	return c.Key
}

type Ingress struct {
	// Enabled specifies whether to enable ingress for the cluster or not.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// Annotations specifies annotations for the ingress.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// IngressClassName specifies the name of the ingress class
	// +optional
	IngressClassName *string `json:"ingressClassName,omitempty"`
	// Host specifies the host for the ingress.
	// +optional
	Host string `json:"host,omitempty"`
	// TLS specifies whether to enable tls for the ingress.
	// The default depends on the ingress provider:
	// - KubernetesIngress: False
	// - NginxIngress: False
	// - QdrantCloudTraefik: Depending on the config.tls setting of the operator.
	// +optional
	TLS *bool `json:"tls,omitempty"`
	// TLSSecretName specifies the name of the secret containing the tls certificate.
	// +optional
	TLSSecretName string `json:"tlsSecretName,omitempty"`
	// NGINX specifies the nginx ingress specific configurations.
	// +optional
	NGINX *NGINXConfig `json:"nginx,omitempty"`
	// Traefik specifies the traefik ingress specific configurations.
	// +optional
	Traefik *TraefikConfig `json:"traefik,omitempty"`
}

func (i *Ingress) GetAnnotations() map[string]string {
	if i == nil {
		return nil
	}
	return i.Annotations
}

func (i *Ingress) GetIngressClassName() *string {
	if i == nil {
		return nil
	}
	return i.IngressClassName
}

func (i *Ingress) GetTls(def bool) bool {
	if i == nil || i.TLS == nil {
		return def
	}
	return *i.TLS
}

func (i *Ingress) GetNGINX() *NGINXConfig {
	if i == nil {
		return nil
	}
	return i.NGINX
}

func (i *Ingress) GetTraefik() *TraefikConfig {
	if i == nil {
		return nil
	}
	return i.Traefik
}

type NGINXConfig struct {
	// AllowedSourceRanges specifies the allowed CIDR source ranges for the ingress.
	// +optional
	AllowedSourceRanges []string `json:"allowedSourceRanges,omitempty"`
	// GRPCHost specifies the host name for the GRPC ingress.
	// +optional
	GRPCHost *string `json:"grpcHost,omitempty"`
}

func (c *NGINXConfig) GetAllowedSourceRanges() []string {
	if c == nil {
		return nil
	}
	return c.AllowedSourceRanges
}

func (c *NGINXConfig) GetGrpcHost() *string {
	if c == nil {
		return nil
	}
	return c.GRPCHost
}

type TraefikConfig struct {
	// AllowedSourceRanges specifies the allowed CIDR source ranges for the ingress.
	// +optional
	AllowedSourceRanges []string `json:"allowedSourceRanges,omitempty"`
	// EntryPoints is the list of traefik entry points to use for the ingress route.
	// If nothing is set, it will take the entryPoints configured in the operator config.
	EntryPoints []string `json:"entryPoints,omitempty"`
}

func (c *TraefikConfig) GetAllowedSourceRanges() []string {
	if c == nil {
		return nil
	}
	return c.AllowedSourceRanges
}

type StorageClassNames struct {
	// DB specifies the storage class name for db volume.
	// +optional
	DB *string `json:"db,omitempty"`
	// Snapshots specifies the storage class name for snapshots volume.
	// +optional
	Snapshots *string `json:"snapshots,omitempty"`
}

func (n *StorageClassNames) GetDB() *string {
	if n == nil {
		return nil
	}
	return n.DB
}

func (n *StorageClassNames) GetSnapshots() *string {
	if n == nil {
		return nil
	}
	return n.Snapshots
}

type ClusterPhase string

//goland:noinspection GoUnusedConst
const (
	ClusterActiveStateSuffix = "ing"
	ClusterFailedStatePrefix = "FailedTo"

	ClusterCreating       ClusterPhase = "Creating"
	ClusterFailedToCreate ClusterPhase = "FailedToCreate"

	ClusterUpdating       ClusterPhase = "Updating"
	ClusterFailedToUpdate ClusterPhase = "FailedToUpdate"

	ClusterScaling   ClusterPhase = "Scaling"
	ClusterUpgrading ClusterPhase = "Upgrading"

	ClusterSuspending      ClusterPhase = "Suspending"
	ClusterSuspended       ClusterPhase = "Suspended"
	ClusterFailedToSuspend ClusterPhase = "FailedToSuspend"

	ClusterResuming       ClusterPhase = "Resuming"
	ClusterFailedToResume ClusterPhase = "FailedToResume"

	ClusterHealthy           ClusterPhase = "Healthy"
	ClusterNotReady          ClusterPhase = "NotReady"
	ClusterRecoveryMode      ClusterPhase = "RecoveryMode"
	ClusterManualMaintenance ClusterPhase = "ManualMaintenance"
)

type ClusterCondition string

//goland:noinspection GoUnusedConst
const (
	ClusterConditionAcceptingConnection ClusterCondition = "AcceptingConnection"
	ClusterConditionRecoveryMode        ClusterCondition = "RecoveryMode"
)

// QdrantClusterStatus defines the observed state of QdrantCluster
// +kubebuilder:pruning:PreserveUnknownFields
type QdrantClusterStatus struct {
	// Phase specifies the phase of the cluster
	// +optional
	Phase ClusterPhase `json:"phase,omitempty"`
	// Reason specifies the reason for the phase of the cluster
	// +optional
	Reason string `json:"reason,omitempty"`
	// AvailableNodes specifies the number of available nodes in the cluster
	// +optional
	AvailableNodes int `json:"availableNodes,omitempty"`
	// AvailableNodeIndexes specifies the indexes of the individual nodes in the cluster
	// The number of indexes should be equal with the AvailableNodes field.
	// +optional
	AvailableNodeIndexes []int `json:"availableNodeIndexes,omitempty"`
	// DeleteInProgessNodeIndexes specifies the indexes of the nodes in the cluster which are in progress of deleting and required to be deleted.
	// The indexes in this field are part of the AvailableNodeIndexes as well and cannot be re-used anymore before they are fully dropped.
	// Meaning that if the cluster-manager has (async) started the process of deleting nodes, due to a scale-down, there is no way to revert this operaton.
	// If the cluster want to scale-up concurrently (aka the delete is in progress), new nodes are required to accomplish.
	// +optional
	DeleteInProgessNodeIndexes []int `json:"deleteInProgressNodeIndexes,omitempty"`
	// BootstrapNode specifies the node in the cluster which will be used for bootstrapping a new node.
	// Should be a value from AvailableNodeIndexes.
	// As default the value from AvailableNodeIndexes[0] will be used.
	// +optional
	BootstrapNode int `json:"bootstrapNode,omitempty"`
	// If set the operator will scale down 1 peer and reset the bool.
	// If you want to remove more then 1 peer, you need to repeat setting this bool
	// +optional
	ScaleDownAllowed bool `json:"scaleDownAllowed,omitempty"`
	// The node index used in a scale down (see ScaleDownAllowed)
	// If this field is not set the last index in AvailableNodeIndexes will be used.
	ScaleDownNodeIndex *int `json:"ScaleDownNodeIndex,omitempty"`
	// Conditions specifies the conditions of different checks on the cluster
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Nodes specifies the status of the nodes in the cluster
	// +optional
	Nodes map[string]NodeStatus `json:"nodes,omitempty"`
	// The next time to invoke the cluster-manager in UTC
	// +optional
	NextClusterManagerInvocation metav1.Time `json:"nextClusterManagerInvocation,omitempty"`
	// The version (to be) used in the cluster.
	// This version can differ from the spec, because version updates need to be done in order (see `update-path` annotation)
	// +optional
	Version string `json:"version,omitempty"`
}

type NodeStatus struct {
	// Name specifies the name of the node
	// +optional
	Name string `json:"name,omitempty"`
	// StartedAt specifies the time when the node started (in RFC3339 format)
	// +optional
	StartedAt string `json:"started_at,omitempty"`
	// States specifies the condition states of the node
	// +optional
	State map[corev1.PodConditionType]corev1.ConditionStatus `json:"state,omitempty"`
	// Version specifies the version of Qdrant running on the node
	// +optional
	Version string `json:"version,omitempty"`
}

//+kubebuilder:object:root=true
// +kubebuilder:resource:path=qdrantclusters,singular=qdrantcluster,shortName=qc;qcs
// +kubebuilder:printcolumn:name="Nodes",type=integer,JSONPath=`.spec.size`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// QdrantCluster is the Schema for the qdrantclusters API
type QdrantCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QdrantClusterSpec   `json:"spec,omitempty"`
	Status QdrantClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// QdrantClusterList contains a list of QdrantCluster
type QdrantClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QdrantCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QdrantCluster{}, &QdrantClusterList{})
}
