package healthchecks

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerServiceHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Service",
	}
	RegisterHealthCheck(gvk, checkServiceHealth)
}

// checkServiceHealth implements ArgoCD-style health check for Services
// A Service is considered healthy when:
// - If type is LoadBalancer: status.loadBalancer.ingress must be populated
// - For all other types: always healthy (ClusterIP, NodePort, ExternalName)
func checkServiceHealth(obj *unstructured.Unstructured) bool {
	// Get spec.type (defaults to ClusterIP if not specified)
	serviceType, found, err := unstructured.NestedString(obj.Object, "spec", "type")
	if err != nil || !found {
		serviceType = "ClusterIP"
	}

	// Only LoadBalancer services need health checking
	if serviceType != "LoadBalancer" {
		return true
	}

	// Check for status.loadBalancer.ingress
	ingress, found, err := unstructured.NestedSlice(obj.Object, "status", "loadBalancer", "ingress")
	if err != nil || !found {
		return false
	}

	// LoadBalancer is healthy if at least one ingress point is assigned
	return len(ingress) > 0
}
