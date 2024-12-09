package v1

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	r := Resources{}
	err := r.Validate("test")
	require.Error(t, err)
	require.ErrorContains(t, err, "test.CPU error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'")

	spec := QdrantClusterSpec{}
	err = spec.Validate()
	require.Error(t, err)
	require.ErrorContains(t, err, "Spec.Resources.CPU error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'")
}

func TestGetQdrantClusterCrdForHash(t *testing.T) {
	qc := QdrantCluster{}
	hash, err := getQdrantClusterCrdHash(qc)
	require.NoError(t, err)
	assert.Equal(t, "523a8cb", hash)

	falseVal := false
	qc.Spec.ServicePerNode = &falseVal
	hash, err = getQdrantClusterCrdHash(qc)
	require.NoError(t, err)
	assert.Equal(t, "523a8cb", hash)

	trueVal := true
	qc.Spec.ServicePerNode = &trueVal
	hash, err = getQdrantClusterCrdHash(qc)
	require.NoError(t, err)
	assert.Equal(t, "523a8cb", hash)
}

// getQdrantClusterCrdHash created a hash for the provided QdrantCluster,
// however a subset only, see GetQdrantClusterCrdForHash for details.
func getQdrantClusterCrdHash(qc QdrantCluster) (string, error) {
	inspect := GetQdrantClusterCrdForHash(qc)
	// Get the hash, so we can diff later
	hash, err := getHash(inspect)
	if err != nil {
		return "", fmt.Errorf("failed to get hash for QdrantCluster: %w", err)
	}
	return hash, err
}

// Get hash of provided value.
// Returns the first 7 characters of the hash (like GitHub).
func getHash(v any) (string, error) {
	json, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshal failed: %w", err)
	}
	// Initialize hash
	hash := sha256.New()
	// add the serialized content
	hash.Write(json)
	// close hash
	sum := hash.Sum(nil)
	// Return first 7 characters
	return fmt.Sprintf("%x", sum)[:7], nil
}
