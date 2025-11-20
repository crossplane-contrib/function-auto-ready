package healthchecks

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerIngressHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "networking.k8s.io",
		Version: "v1",
		Kind:    "Ingress",
	}
	RegisterHealthCheck(gvk, checkIngressHealth)
}

// checkIngressHealth implements ArgoCD-style health check for Ingresses
// An Ingress is considered healthy when:
// - status.loadBalancer.ingress is populated (LoadBalancer has been assigned)
func checkIngressHealth(obj *unstructured.Unstructured) bool {
	// Check for status.loadBalancer.ingress
	ingress, found, err := unstructured.NestedSlice(obj.Object, "status", "loadBalancer", "ingress")
	if err != nil || !found {
		return false
	}

	// Ingress is healthy if at least one ingress point is assigned
	return len(ingress) > 0
}
