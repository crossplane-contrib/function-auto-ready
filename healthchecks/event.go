package healthchecks

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerEventHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Event",
	}
	RegisterHealthCheck(gvk, alwaysReady)
}
