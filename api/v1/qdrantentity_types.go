package v1

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

//goland:noinspection GoUnusedConst
const (
	KindQdrantEntity     = "QdrantEntity"
	ResourceQdrantEntity = "qdrantentities"
)

// QdrantEntitySpec defines the desired state of QdrantEntity
type QdrantEntitySpec struct {
	// The unique identifier of the entity (in UUID format).
	Id string `json:"id,omitempty"`
	// The type of the entity.
	EntityType string `json:"entity_type,omitempty"`
	// Timestamp when the entity was created.
	CreatedAt metav1.Time `json:"created_at,omitempty"`
	// Timestamp when the entity was last updated.
	LastUpdatedAt metav1.Time `json:"last_updated_at,omitempty"`
	// Timestamp when the entity was deleted (or is started to be deleting).
	// If not set the entity is not deleted
	DeletedAt metav1.Time `json:"deleted_at,omitempty"`
	// Generic payload for this entity
	Payload apiextensions.JSON `json:"payload,omitempty"`
	// The hash (sha-256) for the current spec for this entity
	// (including Id, EntityTpe, CreatedAt, LastUpdatedDate, DeletedAt and the Payload)
	// This is for easier comparisson.
	SpecHash string `json:"spechash,omitempty"`
}

// GetPayloadForGRPC gets the current payload
func (r QdrantEntitySpec) GetPayloadForGRPC() (*structpb.Struct, error) {
	return apiextensionsJSONToStructpb(r.Payload)
}

// SetPayloadFromGRPC sets the current payload
func (r *QdrantEntitySpec) SetPayloadFromGRPC(payload *structpb.Struct) error {
	if r == nil {
		return nil
	}
	jsonPayload, err := structpbToApiextensionsJSON(payload)
	if err != nil {
		return err
	}
	r.Payload = jsonPayload
	return nil
}

type EntityPhase string

//goland:noinspection GoUnusedConst
const (
	EntityPhaseCreating EntityPhase = "Creating"
	EntityPhaseReady    EntityPhase = "Ready"
	EntityPhaseFailing  EntityPhase = "Failing"
	EntityPhaseDeleting EntityPhase = "Deleting"
	EntityPhaseDeleted  EntityPhase = "Deleted"
)

// QdrantEntitySpecStatus defines the observed state of QdrantEntitySpec
// +kubebuilder:pruning:PreserveUnknownFields

type QdrantEntityStatus struct {
	// Phase is the current phase of the entity
	// +kubebuilder:validation:Enum=Creating;Ready;Failing;Deleting;Deleted
	Phase EntityPhase `json:"phase,omitempty"`
	// Result is the last result from the invocation to a manager
	Result QdrantEntityStatusResult `json:"result,omitempty"`
	// The hash (sha-256) for the current status
	// (including Phase and Result, not including the hashes).
	// This is for easier comparisson.
	StatusHash string `json:"status_hash,omitempty"`
	// Last used hash from the spec to invoke the manager
	// If the hash of the spec and this hash differs the manager needs to be invoked.
	LastUsedSpecHash string `json:"last_used_spec_hash,omitempty"`
	// Last used hash from the status update towards the backend.
	// If the hash of the status and this hash differs the backend needs to be invoked.
	LastUsedStatusHash string `json:"last_used_status_hash,omitempty"`
}

// EntityResult is the last result from the invocation to a manager
type EntityResult string

//goland:noinspection GoUnusedConst
const (
	EntityResultOk       EntityResult = "Ok"
	EntityRersultPending EntityResult = "Pending"
	EntityResultError    EntityResult = "Error"
)

// QdrantEntityStatusResult is the last result from the invocation to a manager
type QdrantEntityStatusResult struct {
	// The result of last reconcile of the entity
	// +kubebuilder:validation:Enum=Ok;Pending;Error
	Result EntityResult `json:"result,omitempty"`
	// The reason of the result (e.g. in case of an error)
	Reason string `json:"reason,omitempty"`
	// The optional payload of the status.
	Payload apiextensions.JSON `json:"payload,omitempty"`
}

// GetPayloadForGRPC gets the current payload
func (r QdrantEntityStatusResult) GetPayloadForGRPC() (*structpb.Struct, error) {
	return apiextensionsJSONToStructpb(r.Payload)
}

// SetPayloadFromGRPC sets the current payload
func (r *QdrantEntityStatusResult) SetPayloadFromGRPC(payload *structpb.Struct) error {
	if r == nil {
		return nil
	}
	jsonPayload, err := structpbToApiextensionsJSON(payload)
	if err != nil {
		return err
	}
	r.Payload = jsonPayload
	return nil
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=qdrantentities,singular=qdrantentity,shortName=qe
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// QdrantEntity is the Schema for the qdrantentities API
type QdrantEntity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QdrantEntitySpec   `json:"spec,omitempty"`
	Status QdrantEntityStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// QdrantEntityList contains a list of QdrantEntity objects
type QdrantEntityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QdrantEntity `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QdrantEntity{}, &QdrantEntityList{})
}

// apiextensionsJSONToStructpb converts apiextensions.JSON to *structpb.Struct.
func apiextensionsJSONToStructpb(in apiextensions.JSON) (*structpb.Struct, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(in.Raw, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal apiextensions.JSON: %w", err)
	}

	result, err := structpb.NewStruct(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create *structpb.Struct: %w", err)
	}
	return result, nil
}

// structpbToApiextensionsJSON converts *structpb.Struct to apiextensions.JSON.
func structpbToApiextensionsJSON(in *structpb.Struct) (apiextensions.JSON, error) {
	if in == nil {
		return apiextensions.JSON{}, nil
	}
	data := in.AsMap()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return apiextensions.JSON{}, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return apiextensions.JSON{Raw: jsonData}, nil
}
