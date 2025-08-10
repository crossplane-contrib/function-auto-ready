// Package v1beta1 contains the input type for this Function
// +kubebuilder:object:generate=true
// +groupName=autoready.fn.crossplane.io
// +versionName=v1beta1
package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// This isn't a custom resource, in the sense that we never install its CRD.
// It is a KRM-like object, so we generate a CRD to describe its schema.

// Input can be used to provide input to this Function.
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=crossplane
type Input struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ExpectedResourceCount is used to determine composite readiness
	ExpectedResourceCount *int `json:"expectedResourceCount,omitempty"`
}
