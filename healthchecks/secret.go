package healthchecks

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerSecretHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	}
	RegisterHealthCheck(gvk, alwaysReady)
}
