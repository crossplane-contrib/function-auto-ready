package healthchecks

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerStatefulSetHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "StatefulSet",
	}
	RegisterHealthCheck(gvk, checkStatefulSetHealth)
}

// checkStatefulSetHealth implements ArgoCD-style health check for StatefulSets
// A StatefulSet is considered healthy when:
// 1. status.currentRevision == status.updateRevision (all pods updated)
// 2. spec.replicas == status.readyReplicas (all replicas ready)
// 3. spec.replicas == status.currentReplicas (all replicas at current revision)
func checkStatefulSetHealth(obj *unstructured.Unstructured) bool {
	// Get spec.replicas (may be nil, defaults to 1)
	specReplicas := int64(1)
	if val, found := getInt64Field(obj.Object, "spec", "replicas"); found {
		specReplicas = val
	}

	// Get status.readyReplicas
	readyReplicas, found := getInt64Field(obj.Object, "status", "readyReplicas")
	if !found {
		return false
	}

	// Get status.currentReplicas
	currentReplicas, found := getInt64Field(obj.Object, "status", "currentReplicas")
	if !found {
		return false
	}

	// Check replica counts match
	if specReplicas != readyReplicas || specReplicas != currentReplicas {
		return false
	}

	// Get status.currentRevision
	currentRevision, found, err := unstructured.NestedString(obj.Object, "status", "currentRevision")
	if err != nil || !found {
		return false
	}

	// Get status.updateRevision
	updateRevision, found, err := unstructured.NestedString(obj.Object, "status", "updateRevision")
	if err != nil || !found {
		return false
	}

	// Check that all pods are at the updated revision
	return currentRevision == updateRevision
}
