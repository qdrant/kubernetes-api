// Package v1 contains API Schema definitions for the qdrant.io v1 API group
// +kubebuilder:object:generate=true
// +groupName=qdrant.io
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "qdrant.io", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = runtime.NewSchemeBuilder(addToGroupVersion)

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// registerTypes adds the given objects to the SchemeBuilder,
// so they are registered for GroupVersion by AddToScheme.
func registerTypes(objects ...runtime.Object) {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(GroupVersion, objects...)
		return nil
	})
}

func addToGroupVersion(s *runtime.Scheme) error {
	metav1.AddToGroupVersion(s, GroupVersion)
	return nil
}
