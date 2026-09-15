package healthchecks

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerPersistentVolumeHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "PersistentVolume",
	}
	RegisterHealthCheck(gvk, checkPersistentVolumeHealth)
}

// checkPersistentVolumeHealth implements health check for PersistentVolumes
// Based on ArgoCD's gitops-engine implementation
func checkPersistentVolumeHealth(obj *unstructured.Unstructured) bool {
	var pv corev1.PersistentVolume
	err := convertFromUnstructured(obj, &pv)
	if err != nil {
		return false
	}

	switch pv.Status.Phase {
	case corev1.VolumeBound:
		return true
	case corev1.VolumeAvailable:
		return true
	case corev1.VolumeReleased:
		return false
	case corev1.VolumeFailed:
		return false
	default:
		return false
	}
}
