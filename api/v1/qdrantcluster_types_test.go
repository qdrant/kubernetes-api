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
			name:          "Neither .spec.storage.size nor .spec.resource.storage is specified",
			spec:          QdrantClusterSpec{},
			expectedError: fmt.Errorf("must specify either .spec.storage.size or .spec.resources.storage"),
		},
		{
			name: "Only .spec.storage.size specified but not .spec.resources.storage",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:    "100m",
					Memory: "128Mi",
				},
				Storage: Storage{
					Size: "2Gi",
				},
			},
			expectedError: nil,
		},
		{
			name: "Only .spec.resources.storage specified but not .spec.storage.size",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Memory:  "128Mi",
					Storage: "2Gi",
				},
			},
			expectedError: nil,
		},
		{
			name: "Both .spec.storage.size and not .spec.resources.storage, specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:     "100m",
					Memory:  "128Mi",
					Storage: "2Gi",
				},
				Storage: Storage{
					Size: "2Gi",
				},
			},
			expectedError: nil,
		},
		{
			name: "Invalid storage size",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:    "100m",
					Memory: "128Mi",
				},
				Storage: Storage{
					Size: "foo",
				},
			},
			expectedError: fmt.Errorf("invalid storage size: foo error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'"),
		},
		{
			name: "CPU amount not specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					Memory: "128Mi",
				},
				Storage: Storage{
					Size: "2Gi",
				},
			},
			expectedError: fmt.Errorf("Spec.Resources.CPU error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'"),
		},

		{
			name: "Invalid CPU amount",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:    "foo",
					Memory: "128Mi",
				},
				Storage: Storage{
					Size: "2Gi",
				},
			},
			expectedError: fmt.Errorf("Spec.Resources.CPU error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'"),
		},
		{
			name: "Memory amount not specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU: "100m",
				},
				Storage: Storage{
					Size: "2Gi",
				},
			},
			expectedError: fmt.Errorf("Spec.Resources.Memory error: quantities must match the regular expression '^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$'"),
		},
		{
			name: "Invalid Memory amount",
			spec: QdrantClusterSpec{
				Resources: Resources{
					CPU:    "100m",
					Memory: "foo",
				},
				Storage: Storage{
					Size: "2Gi",
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

func TestGetStorageSize(t *testing.T) {
	testCases := []struct {
		name            string
		spec            QdrantClusterSpec
		expectedStorage string
	}{
		{
			name:            "Neither .spec.storage.size nor .spec.resources.storage specified",
			spec:            QdrantClusterSpec{},
			expectedStorage: "",
		},
		{
			name: "Only .spec.storage.size specified",
			spec: QdrantClusterSpec{
				Storage: Storage{
					Size: "2Gi",
				},
			},
			expectedStorage: "2Gi",
		},
		{
			name: "Only .spec.resources.storage specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					Storage: "2Gi",
				},
			},
			expectedStorage: "2Gi",
		},
		{
			name: "Both .spec.storage.size and .spec.resources.storage specified",
			spec: QdrantClusterSpec{
				Resources: Resources{
					Storage: "1Gi",
				},
				Storage: Storage{
					Size: "2Gi",
				},
			},
			expectedStorage: "2Gi",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.spec.GetStorageSize()
			assert.Equal(t, tt.expectedStorage, actual)
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
