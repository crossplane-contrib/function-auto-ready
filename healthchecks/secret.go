package healthchecks

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerSecretHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	}
	RegisterHealthCheck(gvk, checkSecretHealth)
}

// checkSecretHealth implements health check for Secrets
// Secrets don't have meaningful status conditions, so if the Secret exists
// in the observed state, it's considered ready.
func checkSecretHealth(obj *unstructured.Unstructured) bool {
	// Secrets are always ready if they exist
	return true
}
