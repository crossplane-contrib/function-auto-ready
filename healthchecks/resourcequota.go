package healthchecks

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerResourceQuotaHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ResourceQuota",
	}
	RegisterHealthCheck(gvk, alwaysReady)
}
