// Package v1beta1 contains the input type for this Function
// +kubebuilder:object:generate=true
// +groupName=autoready.fn.crossplane.io
// +versionName=v1beta1
package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This isn't a custom resource, in the sense that we never install its CRD.
// It is a KRM-like object, so we generate a CRD to describe its schema.

// Input is used to provide inputs to this Function.
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=crossplane
type Input struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	// TTL for which a response can be cached in time.Duration format
	// +kubebuilder:default="1m0s"
	// +optional
	TTL string `json:"ttl"`

	// CELHealthCheckCustomizationFrom is a reference to fetch CEL health check customizations from context
	// +kubebuilder:validation:Optional
	CELHealthCheckCustomizationFrom *string `json:"celHealthCheckCustomizationFrom,omitempty"`
}
