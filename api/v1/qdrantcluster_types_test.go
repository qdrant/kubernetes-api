package v1

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		name          string
		spec          QdrantClusterSpec
		expectedError error
	}{
		{
			name: "Storage size is not specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:    "100m",
					Memory: "128Mi",
				},
			},
			expectedError: fmt.Errorf("Spec.Resources.Storage error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'"),
		},
		{
			name: "Invalid storage size",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Memory:  "128Mi",
					Storage: "foo",
				},
			},
			expectedError: fmt.Errorf("Spec.Resources.Storage error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'"),
		},
		{
			name: "CPU amount is not specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					Memory:  "128Mi",
					Storage: "2Gi",
				},
			},
			expectedError: fmt.Errorf("Spec.Resources.CPU error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'"),
		},

		{
			name: "Invalid CPU amount",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "foo",
					Memory:  "128Mi",
					Storage: "2Gi",
				},
			},
			expectedError: fmt.Errorf("Spec.Resources.CPU error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'"),
		},
		{
			name: "Memory amount  is not specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Storage: "2Gi",
				},
			},
			expectedError: fmt.Errorf("Spec.Resources.Memory error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'"),
		},
		{
			name: "Invalid Memory amount",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Memory:  "foo",
					Storage: "2Gi",
				},
			},
			expectedError: fmt.Errorf("Spec.Resources.Memory error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spec.Validate()
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError.Error())
			}
		})
	}
}
