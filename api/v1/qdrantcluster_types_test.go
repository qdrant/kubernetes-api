package v1

import (
	"testing"

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
