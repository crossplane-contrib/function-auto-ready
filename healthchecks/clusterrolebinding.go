package healthchecks

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerClusterRoleBindingHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
		Kind:    "ClusterRoleBinding",
	}
	RegisterHealthCheck(gvk, alwaysReady)
}
