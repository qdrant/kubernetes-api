package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

// TestQdrantClusterRoutingGetSpecIsNilSafe covers the entry point of the getter
// chain: a nil object must read as the zero spec so callers can chain without a
// nil check of their own.
func TestQdrantClusterRoutingGetSpecIsNilSafe(t *testing.T) {
	var qcrt *QdrantClusterRouting

	assert.Equal(t, QdrantClusterRoutingSpec{}, qcrt.GetSpec())
	// The whole point of the value return: this chain must not panic.
	assert.Empty(t, qcrt.GetSpec().GetAllowedSourceRanges())
	assert.Empty(t, qcrt.GetSpec().GetClusterId())
	assert.Empty(t, qcrt.GetSpec().GetFQDN())
	assert.False(t, qcrt.GetSpec().GetShared())
	assert.False(t, qcrt.GetSpec().GetDedicated())
	assert.False(t, qcrt.GetSpec().GetTLS())
	assert.False(t, qcrt.GetSpec().GetEnableAccessLog())
	assert.False(t, qcrt.GetSpec().GetMultiAZ())
	assert.Empty(t, qcrt.GetSpec().GetNodeIndexes())
}

// TestQdrantClusterRoutingGetSpecReturnsSpec asserts GetSpec hands back the
// actual spec, not a zero value, for a populated object.
func TestQdrantClusterRoutingGetSpecReturnsSpec(t *testing.T) {
	qcrt := &QdrantClusterRouting{
		Spec: QdrantClusterRoutingSpec{
			ClusterId:           "cluster-id",
			FQDN:                "cluster-id.example.com",
			AllowedSourceRanges: []string{"1.2.3.4/32"},
		},
	}

	assert.Equal(t, "cluster-id", qcrt.GetSpec().GetClusterId())
	assert.Equal(t, "cluster-id.example.com", qcrt.GetSpec().GetFQDN())
	assert.Equal(t, []string{"1.2.3.4/32"}, qcrt.GetSpec().GetAllowedSourceRanges())
}

// TestQdrantClusterRoutingSpecDefaults pins the unset behaviour of every
// optional field. Enabled and ServicePerNode default to true to match their
// +kubebuilder:default markers, so an unset pointer must NOT read as false the
// way a plain dereference helper would report it.
func TestQdrantClusterRoutingSpecDefaults(t *testing.T) {
	var spec QdrantClusterRoutingSpec

	assert.True(t, spec.GetEnabled(), "unset Enabled defaults to true (CRD default)")
	assert.True(t, spec.GetServicePerNode(), "unset ServicePerNode defaults to true (CRD default)")

	assert.False(t, spec.GetShared())
	assert.False(t, spec.GetDedicated())
	assert.False(t, spec.GetTLS())
	assert.False(t, spec.GetEnableAccessLog())
	assert.False(t, spec.GetMultiAZ())

	assert.Empty(t, spec.GetClusterId())
	assert.Empty(t, spec.GetFQDN())
	assert.Empty(t, spec.GetNodeIndexes())
	assert.Empty(t, spec.GetAllowedSourceRanges())
}

// TestQdrantClusterRoutingSpecGetters asserts every getter reports the value it
// was given, including the false-but-set case that the default-aware getters
// must not override.
func TestQdrantClusterRoutingSpecGetters(t *testing.T) {
	spec := QdrantClusterRoutingSpec{
		ClusterId:           "cluster-id",
		FQDN:                "cluster-id.example.com",
		Enabled:             ptr.To(false),
		Shared:              ptr.To(true),
		Dedicated:           ptr.To(true),
		TLS:                 ptr.To(true),
		ServicePerNode:      ptr.To(false),
		NodeIndexes:         []int{0, 1, 2},
		AllowedSourceRanges: []string{"1.2.3.4/32", "10.0.0.0/8"},
		EnableAccessLog:     ptr.To(true),
		MultiAZ:             true,
	}

	assert.Equal(t, "cluster-id", spec.GetClusterId())
	assert.Equal(t, "cluster-id.example.com", spec.GetFQDN())
	assert.False(t, spec.GetEnabled(), "an explicit false must win over the default")
	assert.True(t, spec.GetShared())
	assert.True(t, spec.GetDedicated())
	assert.True(t, spec.GetTLS())
	assert.False(t, spec.GetServicePerNode(), "an explicit false must win over the default")
	assert.Equal(t, []int{0, 1, 2}, spec.GetNodeIndexes())
	assert.Equal(t, []string{"1.2.3.4/32", "10.0.0.0/8"}, spec.GetAllowedSourceRanges())
	assert.True(t, spec.GetEnableAccessLog())
	assert.True(t, spec.GetMultiAZ())
}
