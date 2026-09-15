package healthchecks

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerNodeHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Node",
	}
	RegisterHealthCheck(gvk, checkNodeHealth)
}

// checkNodeHealth implements health check for Nodes
// Based on ArgoCD's gitops-engine implementation
func checkNodeHealth(obj *unstructured.Unstructured) bool {
	var node corev1.Node
	err := convertFromUnstructured(obj, &node)
	if err != nil {
		return false
	}

	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
			return true
		}
	}

	return false
}
