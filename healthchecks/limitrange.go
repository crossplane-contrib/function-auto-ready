package healthchecks

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerLimitRangeHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "LimitRange",
	}
	RegisterHealthCheck(gvk, alwaysReady)
}
