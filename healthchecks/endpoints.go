package healthchecks

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerEndpointsHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Endpoints",
	}
	RegisterHealthCheck(gvk, alwaysReady)
}
