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

func TestGetStorageClassNames(t *testing.T) {
	testCases := []struct {
		name     string
		spec     QdrantClusterSpec
		expected *StorageClassNames
	}{
		{
			name:     "Neither .spec.storageClassNames nor .spec.storage.storageClassNames specified",
			spec:     QdrantClusterSpec{},
			expected: nil,
		},
		{
			name: "Only .spec.storageClassNames specified",
			spec: QdrantClusterSpec{
				StorageClassNames: &StorageClassNames{
					DB:        ptr.To("foo"),
					Snapshots: ptr.To("bar"),
				},
			},
			expected: &StorageClassNames{
				DB:        ptr.To("foo"),
				Snapshots: ptr.To("bar"),
			},
		},
		{
			name: "Only .spec.storage.storageClassNames specified",
			spec: QdrantClusterSpec{
				Storage: Storage{
					StorageClassNames: &StorageClassNames{
						DB:        ptr.To("foo"),
						Snapshots: ptr.To("bar"),
					},
				},
			},
			expected: &StorageClassNames{
				DB:        ptr.To("foo"),
				Snapshots: ptr.To("bar"),
			},
		},
		{
			name: "Both .spec.storageClassNames and .spec.storage.storageClassNames specified",
			spec: QdrantClusterSpec{
				StorageClassNames: &StorageClassNames{
					DB:        ptr.To("foo-old"),
					Snapshots: ptr.To("bar-old"),
				},
				Storage: Storage{
					StorageClassNames: &StorageClassNames{
						DB:        ptr.To("foo"),
						Snapshots: ptr.To("bar"),
					},
				},
			},
			expected: &StorageClassNames{
				DB:        ptr.To("foo"),
				Snapshots: ptr.To("bar"),
			},
		},
		{
			name: "DB and Snapshot storageclass specified in different places",
			spec: QdrantClusterSpec{
				StorageClassNames: &StorageClassNames{
					DB: ptr.To("foo"),
				},
				Storage: Storage{
					StorageClassNames: &StorageClassNames{
						Snapshots: ptr.To("bar"),
					},
				},
			},
			expected: &StorageClassNames{
				DB:        ptr.To("foo"),
				Snapshots: ptr.To("bar"),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.spec.GetStorageClassNames()
			if tt.expected == nil || actual == nil {
				assert.Equal(t, tt.expected, actual)
			} else {
				assert.EqualValues(t, tt.expected, actual)
			}
		})
	}
}
