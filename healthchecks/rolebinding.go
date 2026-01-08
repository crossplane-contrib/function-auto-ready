package healthchecks

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerRoleBindingHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
		Kind:    "RoleBinding",
	}
	RegisterHealthCheck(gvk, alwaysReady)
}
