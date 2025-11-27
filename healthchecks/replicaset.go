package healthchecks

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerReplicaSetHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "ReplicaSet",
	}
	RegisterHealthCheck(gvk, checkReplicaSetHealth)
}

// checkReplicaSetHealth implements health check for ReplicaSets
// Based on ArgoCD's gitops-engine implementation
func checkReplicaSetHealth(obj *unstructured.Unstructured) bool {
	var rs appsv1.ReplicaSet
	err := convertFromUnstructured(obj, &rs)
	if err != nil {
		return false
	}

	// Check if observed generation matches
	if rs.Status.ObservedGeneration < rs.Generation {
		return false
	}

	// Check for replica failure condition
	for _, condition := range rs.Status.Conditions {
		if condition.Type == appsv1.ReplicaSetReplicaFailure && condition.Status == "True" {
			return false
		}
	}

	// Check if available replicas match desired replicas
	desiredReplicas := int32(1)
	if rs.Spec.Replicas != nil {
		desiredReplicas = *rs.Spec.Replicas
	}

	if rs.Status.AvailableReplicas < desiredReplicas {
		return false
	}

	return true
}
