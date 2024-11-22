package v1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KindQdrantCluster     = "QdrantCluster"
	ResourceQdrantCluster = "qdrantclusters"
)

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
)

// GetQdrantClusterCrdForHash creates a QdrantCluster for the use of creating a hash for the provided QdrantCluster,
// It uses a subset only, so it ignores e.g. majority of the the status and some fields from the spec.
func GetQdrantClusterCrdForHash(qc QdrantCluster) QdrantCluster {
	// Only include the items to be inspected (basically the spec and a single annot)
	result := QdrantCluster{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
		},
	}
	// Add the restartAt annot if needed
	if annots := qc.GetAnnotations(); annots != nil {
		if val, found := annots[RestartedAtAnnotationKey]; found {
			result.ObjectMeta.Annotations[RestartedAtAnnotationKey] = val
		}
	}
	cloned := qc.Spec.DeepCopy()
	// Version is a special case, we can upgrade via an upgade-path,
	// so we should take care of the version in status instead
	if v := qc.Status.Version; v != "" {
		cloned.Version = v
	}
	// Remove all fields (aka set a fixed value) which shouldn't restart the pod
	// The list is sorted alphabetically, for easier maintainability
	cloned.ClusterManager = nil
	cloned.Distributed = false
	cloned.Ingress = nil
	cloned.Pauses = nil
	cloned.Resources.Storage = ""
	cloned.RestartAllPodsConcurrently = false
	cloned.Service = nil
	cloned.ServicePerNode = false
	cloned.Size = 1
	if v := cloned.StatefulSet; v != nil {
		v.Annotations = nil
	}
	cloned.StorageClassNames = nil
	cloned.Suspend = false
	// Set Spec for result
	result.Spec = *cloned
	// Return result
	return result
}

// QdrantClusterSpec defines the desired state of QdrantCluster
// +kubebuilder:pruning:PreserveUnknownFields
type QdrantClusterSpec struct {
	// Id specifies the unique identifier of the cluster
	Id string `json:"id"`
	// Version specifies the version of Qdrant to deploy
	Version string `json:"version"`
	// Size specifies the desired number of Qdrant nodes in the cluster
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=30
	Size int `json:"size"`
	// ServicePerNode specifies whether the cluster should start a dedicated service for each node.
	// +kubebuilder:default=true
	// +optional
	ServicePerNode bool `json:"servicePerNode"`
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
	// Deprecated
	// +optional
	Distributed bool `json:"distributed,omitempty"`
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
	// need to be restarted like when RestartedAtAnnotationKey is added/updated or the Qdrant version need to be upgraded.
	// This helps sharded but not replicated clusters to reduce downtime to possible minimum during restart.
	// +optional
	RestartAllPodsConcurrently bool `json:"restartAllPodsConcurrently,omitempty"`
}

// Validates if there are incorrect settings in the CRD
func (s QdrantClusterSpec) Validate() error {
	if err := s.Resources.Validate("Spec.Resources"); err != nil {
		return err
	}
	return nil
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

// Validates if there are incorrect settings in the CRD
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

// Validates if there are incorrect settings in the CRD
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
}

type StorageConfig struct {
	// Performance configuration
	// +optional
	Performance *StoragePerformanceConfig `json:"performance,omitempty"`
}

type StoragePerformanceConfig struct {
	// OptimizerCPUBudget defines the number of CPU allocation.
	// If 0 - auto selection, keep 1 or more CPUs unallocated depending on CPU size
	// If negative - subtract this number of CPUs from the available CPUs.
	// If positive - use this exact number of CPUs.
	// +optional
	OptimizerCPUBudget *int64 `json:"optimizerCPUBudget,omitempty"`
	// AsyncScorer enables io_uring when rescoring
	// +optional
	AsyncScorer *bool `json:"asyncScorer,omitempty"`
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
}

func (c *QdrantConfigurationTLS) GetCert() *QdrantSecretKeyRef {
	if c == nil {
		return nil
	}
	return c.Cert
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

const (
	ClusterActiveStateSuffix = "ing"
	ClusterFailedStatePrefix = "FailedTo"

	ClusterCreating       ClusterPhase = "Creating"
	ClusterFailedToCreate ClusterPhase = "FailedToCreate"

	ClusterUpdating       ClusterPhase = "Updating"
	ClusterFailedToUpdate ClusterPhase = "FailedToUpdate"

	ClusterScaling       ClusterPhase = "Scaling"
	ClusterFailedToScale ClusterPhase = "FailedToScale"

	ClusterRestarting      ClusterPhase = "Restarting"
	ClusterFailedToRestart ClusterPhase = "FailedToRestart"

	ClusterResyncing      ClusterPhase = "Resyncing"
	ClusterFailedToResync ClusterPhase = "FailedToResync"

	ClusterUpgrading       ClusterPhase = "Upgrading"
	ClusterFailedToUpgrade ClusterPhase = "FailedToUpgrade"

	ClusterBackupRunning  ClusterPhase = "BackupRunning"
	ClusterFailedToBackup ClusterPhase = "FailedToBackup"

	ClusterRestoring       ClusterPhase = "Restoring"
	ClusterFailedToRestore ClusterPhase = "FailedToRestore"

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
	// Operations tracks list of recent operation on the cluster. Operator uses this field to track the progress of an operation.
	// In operator V2 this field is deprecated and Kubernetes events are used instead.
	// +optional
	Operations []Operation `json:"operations,omitempty"`
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

type Operation struct {
	// Type specifies the type of the operation
	Type OperationType `json:"type"`
	// Phase specifies the phase of the operation
	// +optional
	Phase OperationPhase `json:"phase,omitempty"`
	// Id specifies the id of the operation
	// +optional
	Id int `json:"id"`
	// StartTime specifies the time when the operation started
	// +optional
	StartTime string `json:"startTime,omitempty"`
	// CompletionTime specifies the time when the operation completed
	// +optional
	CompletionTime string `json:"completionTime,omitempty"`
	// Message specifies the message of the operation
	// +optional
	Message string `json:"message,omitempty"`
	// SubOperation specifies whether the operation is a sub-operation of another operation
	// +optional
	SubOperation bool `json:"subOperation,omitempty"`
	// Steps specifies the steps the operation has performed
	// +optional
	Steps []OperationStep `json:"steps,omitempty"`
}

type OperationType string

const (
	BackupOperation                   OperationType = "Backup"
	ClusterCreationOperation          OperationType = "ClusterCreation"
	HorizontalScalingOperation        OperationType = "HorizontalScaling"
	VerticalScalingOperation          OperationType = "VerticalScaling"
	VersionUpdateOperation            OperationType = "VersionUpdate"
	SuspendOperation                  OperationType = "Suspend"
	ResumeOperation                   OperationType = "Resume"
	RestartOperation                  OperationType = "Restart"
	ResyncOperation                   OperationType = "Resync"
	RecoveryOperation                 OperationType = "Recovery"
	CrossNamespacedMigrationOperation OperationType = "CrossNamespacedMigration"
)

type OperationPhase string

const (
	OperationPending    OperationPhase = "Pending"
	OperationInProgress OperationPhase = "InProgress"
	OperationCompleted  OperationPhase = "Completed"
	OperationFailed     OperationPhase = "Failed"
)

type OperationStep struct {
	// Name specifies the name of the step
	Name string `json:"name"`
	// Id specifies the id of the step
	// +optional
	Id int `json:"id,omitempty"`
	// Phase specifies the phase of the step
	// +optional
	Phase StepPhase `json:"phase,omitempty"`
	// Message specifies the reason in case of failure
	// +optional
	Message string `json:"message,omitempty"`
}

type StepPhase string

const (
	StepInProgress StepPhase = "InProgress"
	StepCompleted  StepPhase = "Completed"
	StepFailed     StepPhase = "Failed"
)

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
