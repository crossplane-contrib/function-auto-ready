package healthchecks

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerConfigMapHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ConfigMap",
	}
	RegisterHealthCheck(gvk, checkConfigMapHealth)
}

// checkConfigMapHealth implements health check for ConfigMaps
// ConfigMaps don't have meaningful status conditions, so if the ConfigMap exists
// in the observed state, it's considered ready.
func checkConfigMapHealth(obj *unstructured.Unstructured) bool {
	// ConfigMaps are always ready if they exist
	return true
}
