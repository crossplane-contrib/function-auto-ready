package healthchecks

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	RegisterHealthCheck(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ServiceAccount",
	}, checkServiceAccountHealth)
}

func checkServiceAccountHealth(obj *unstructured.Unstructured) bool {
	return true
}
