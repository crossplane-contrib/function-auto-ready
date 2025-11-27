package healthchecks

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerDaemonSetHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "DaemonSet",
	}
	RegisterHealthCheck(gvk, checkDaemonSetHealth)
}

// checkDaemonSetHealth implements ArgoCD-style health check for DaemonSets
// A DaemonSet is considered healthy when:
// 1. status.desiredNumberScheduled == status.numberReady
// 2. status.updatedNumberScheduled == status.desiredNumberScheduled
// 3. status.numberAvailable == status.desiredNumberScheduled
func checkDaemonSetHealth(obj *unstructured.Unstructured) bool {
	// Get status.desiredNumberScheduled
	desiredNumberScheduled, found := getInt64Field(obj.Object, "status", "desiredNumberScheduled")
	if !found {
		return false
	}

	// Get status.numberReady
	numberReady, found := getInt64Field(obj.Object, "status", "numberReady")
	if !found {
		return false
	}

	// Get status.updatedNumberScheduled
	updatedNumberScheduled, found := getInt64Field(obj.Object, "status", "updatedNumberScheduled")
	if !found {
		return false
	}

	// Get status.numberAvailable
	numberAvailable, found := getInt64Field(obj.Object, "status", "numberAvailable")
	if !found {
		return false
	}

	// Check all numbers match
	return desiredNumberScheduled == numberReady &&
		updatedNumberScheduled == desiredNumberScheduled &&
		numberAvailable == desiredNumberScheduled
}
