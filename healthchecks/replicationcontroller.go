package healthchecks

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerReplicationControllerHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ReplicationController",
	}
	RegisterHealthCheck(gvk, checkReplicationControllerHealth)
}

// checkReplicationControllerHealth implements health check for ReplicationControllers
// A ReplicationController is considered healthy when:
// 1. spec.replicas == status.replicas
// 2. spec.replicas == status.readyReplicas
// 3. spec.replicas == status.availableReplicas
func checkReplicationControllerHealth(obj *unstructured.Unstructured) bool {
	var rc corev1.ReplicationController
	err := convertFromUnstructured(obj, &rc)
	if err != nil {
		return false
	}

	desired := int32(1)
	if rc.Spec.Replicas != nil {
		desired = *rc.Spec.Replicas
	}

	return rc.Status.Replicas == desired &&
		rc.Status.ReadyReplicas == desired &&
		rc.Status.AvailableReplicas == desired
}
