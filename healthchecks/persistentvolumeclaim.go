package healthchecks

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerPersistentVolumeClaimHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "PersistentVolumeClaim",
	}
	RegisterHealthCheck(gvk, checkPersistentVolumeClaimHealth)
}

// checkPersistentVolumeClaimHealth implements health check for PersistentVolumeClaims
// Based on ArgoCD's gitops-engine implementation
func checkPersistentVolumeClaimHealth(obj *unstructured.Unstructured) bool {
	var pvc corev1.PersistentVolumeClaim
	err := convertFromUnstructured(obj, &pvc)
	if err != nil {
		return false
	}

	switch pvc.Status.Phase {
	case corev1.ClaimBound:
		return true
	case corev1.ClaimLost:
		return false
	case corev1.ClaimPending:
		return false
	default:
		return false
	}
}
