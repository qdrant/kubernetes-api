package v1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestQdrantClusterSpecJSONOmitsUnsetOnDemandReplication(t *testing.T) {
	data, err := json.Marshal(QdrantClusterSpec{})

	assert.NoError(t, err)
	assert.NotContains(t, string(data), "onDemandReplication")
}

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

// TestIngressGetEnableAccessLog pins the three-state override. The false case is
// the one that matters operationally: a single tenant can carry most of a
// region's traffic, so it has to be possible to silence one cluster without
// giving up the log for every other cluster in the region.
func TestIngressGetEnableAccessLog(t *testing.T) {
	testCases := []struct {
		name    string
		ingress *Ingress
		region  bool
		want    bool
	}{
		{"nil ingress follows the region", nil, true, true},
		{"unset follows the region (on)", &Ingress{}, true, true},
		{"unset follows the region (off)", &Ingress{}, false, false},
		{"true forces on against a disabled region", &Ingress{EnableAccessLog: ptr.To(true)}, false, true},
		{"false forces off against an enabled region", &Ingress{EnableAccessLog: ptr.To(false)}, true, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.ingress.GetEnableAccessLog(tc.region))
		})
	}
}

func TestIngressJSONOmitsUnsetEnableAccessLog(t *testing.T) {
	data, err := json.Marshal(Ingress{})

	assert.NoError(t, err)
	assert.NotContains(t, string(data), "enableAccessLog")
}
