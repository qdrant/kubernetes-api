package v1

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
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
		{
			name: "No storage configuration",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Memory:  "1Gi",
					Storage: "2Gi",
				},
			},
			expectedError: nil,
		},
		{
			name: "Empty storage configuration",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Memory:  "1Gi",
					Storage: "2Gi",
				},
				Storage: &Storage{},
			},
			expectedError: nil,
		},
		{
			name: "Only VolumeAttributeClassName specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Memory:  "1Gi",
					Storage: "2Gi",
				},
				Storage: &Storage{
					VolumeAttributesClassName: ptr.To("foo"),
				},
			},
			expectedError: nil,
		},

		{
			name: "Both VolumeAttributeClassName and IOPS/Throughput specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Memory:  "1Gi",
					Storage: "2Gi",
				},
				Storage: &Storage{
					VolumeAttributesClassName: ptr.To("foo"),
					IOPS:                      ptr.To(10000),
					Throughput:                ptr.To(500),
				},
			},
			expectedError: fmt.Errorf(".spec.storage: can not specify both VolumeAttributesClassName and IOPS/Throughput"),
		},
		{
			name: "Only IOPS specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Memory:  "1Gi",
					Storage: "2Gi",
				},
				Storage: &Storage{
					IOPS: ptr.To(10000),
				},
			},
			expectedError: fmt.Errorf(".spec.storage: must specify both IOPS and Throughput"),
		},
		{
			name: "Only Throughput specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Memory:  "1Gi",
					Storage: "2Gi",
				},
				Storage: &Storage{
					Throughput: ptr.To(500),
				},
			},
			expectedError: fmt.Errorf(".spec.storage: must specify both IOPS and Throughput"),
		},
		{
			name: "Both IOPS/Throughput specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Memory:  "1Gi",
					Storage: "2Gi",
				},
				Storage: &Storage{
					IOPS:       ptr.To(10000),
					Throughput: ptr.To(500),
				},
			},
			expectedError: nil,
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
