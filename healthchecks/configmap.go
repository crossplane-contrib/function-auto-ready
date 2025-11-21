package healthchecks

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerConfigMapHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ConfigMap",
	}
	RegisterHealthCheck(gvk, alwaysReady)
}
