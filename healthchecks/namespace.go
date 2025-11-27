package healthchecks

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerNamespaceHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Namespace",
	}
	RegisterHealthCheck(gvk, alwaysReady)
}
