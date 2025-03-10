package v1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Test apiextensionsJSONToStructpb
func TestApiextensionsJSONToStructpb(t *testing.T) {
	testCases := []struct {
		name          string
		input         apiextensions.JSON
		expected      *structpb.Struct
		expectedError string
	}{
		{
			name: "valid JSON",
			input: apiextensions.JSON{
				Raw: []byte(`{"name": "John Doe", "age": 30}`),
			},
			expected: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"name": {Kind: &structpb.Value_StringValue{StringValue: "John Doe"}},
					"age":  {Kind: &structpb.Value_NumberValue{NumberValue: 30}},
				},
			},
		},
		{
			name: "valid JSON with nested struct",
			input: apiextensions.JSON{
				Raw: []byte(`{"name": "John Doe", "age": 30, "address": {"street": "mainstreet", "zip": "1234"}}`),
			},
			expected: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"name": {Kind: &structpb.Value_StringValue{StringValue: "John Doe"}},
					"age":  {Kind: &structpb.Value_NumberValue{NumberValue: 30}},
					"address": {Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"street": {Kind: &structpb.Value_StringValue{StringValue: "mainstreet"}},
								"zip":    {Kind: &structpb.Value_StringValue{StringValue: "1234"}},
							},
						},
					}},
				},
			},
		},
		{
			name: "valid JSON with array",
			input: apiextensions.JSON{
				Raw: []byte(`{"name": "John Doe", "age": 30, "tags": ["tag1", "tag2"]}`),
			},
			expected: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"name": {Kind: &structpb.Value_StringValue{StringValue: "John Doe"}},
					"age":  {Kind: &structpb.Value_NumberValue{NumberValue: 30}},
					"tags": {Kind: &structpb.Value_ListValue{
						ListValue: &structpb.ListValue{
							Values: []*structpb.Value{
								{Kind: &structpb.Value_StringValue{StringValue: "tag1"}},
								{Kind: &structpb.Value_StringValue{StringValue: "tag2"}},
							},
						},
					}},
				},
			},
		},
		{
			name:          "invalid JSON",
			input:         apiextensions.JSON{Raw: []byte(`{"name": "John Doe", "age": }`)},
			expected:      nil,
			expectedError: "failed to unmarshal apiextensions.JSON",
		},
		{
			name:          "empty JSON",
			input:         apiextensions.JSON{Raw: []byte(``)},
			expected:      &structpb.Struct{Fields: map[string]*structpb.Value{}},
			expectedError: "",
		},
		{
			name:          "empty JSON-object",
			input:         apiextensions.JSON{Raw: []byte(`{}`)},
			expected:      &structpb.Struct{Fields: map[string]*structpb.Value{}},
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		// Run the individual test case (added the name to the runner)
		t.Run(tc.name, func(t *testing.T) {
			result, err := apiextensionsJSONToStructpb(tc.input)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				return
			}
			require.NoError(t, err)

			// We use this special method for comparing protobufs
			if diff := cmp.Diff(tc.expected, result, cmpopts.IgnoreUnexported(structpb.Struct{}, structpb.Value{}, structpb.ListValue{})); diff != "" {
				t.Errorf("apiextensionsJSONToStructpb() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// Test structpbToApiextensionsJSON
func TestStructpbToApiextensionsJSON(t *testing.T) {
	testCases := []struct {
		name          string
		input         *structpb.Struct
		expected      apiextensions.JSON
		expectedError string
	}{
		{
			name: "valid struct",
			input: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"name": {Kind: &structpb.Value_StringValue{StringValue: "John Doe"}},
					"age":  {Kind: &structpb.Value_NumberValue{NumberValue: 30}},
				},
			},
			expected: apiextensions.JSON{Raw: []byte(`{"age":30,"name":"John Doe"}`)},
		},
		{
			name: "valid struct with nested struct",
			input: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"name": {Kind: &structpb.Value_StringValue{StringValue: "John Doe"}},
					"age":  {Kind: &structpb.Value_NumberValue{NumberValue: 30}},
					"address": {Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"street": {Kind: &structpb.Value_StringValue{StringValue: "mainstreet"}},
								"zip":    {Kind: &structpb.Value_StringValue{StringValue: "1234"}},
							},
						},
					}},
				},
			},
			expected: apiextensions.JSON{Raw: []byte(`{"address":{"street":"mainstreet","zip":"1234"},"age":30,"name":"John Doe"}`)},
		},
		{
			name: "valid struct with array",
			input: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"name": {Kind: &structpb.Value_StringValue{StringValue: "John Doe"}},
					"age":  {Kind: &structpb.Value_NumberValue{NumberValue: 30}},
					"tags": {Kind: &structpb.Value_ListValue{
						ListValue: &structpb.ListValue{
							Values: []*structpb.Value{
								{Kind: &structpb.Value_StringValue{StringValue: "tag1"}},
								{Kind: &structpb.Value_StringValue{StringValue: "tag2"}},
							},
						},
					}},
				},
			},
			expected: apiextensions.JSON{Raw: []byte(`{"age":30,"name":"John Doe","tags":["tag1","tag2"]}`)},
		},
		{
			name:          "nil struct",
			input:         nil,
			expected:      apiextensions.JSON{Raw: []byte{}},
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		// Run the individual test case (added the name to the runner)
		t.Run(tc.name, func(t *testing.T) {
			result, err := structpbToApiextensionsJSON(tc.input)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
