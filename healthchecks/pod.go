package healthchecks

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerPodHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	}
	RegisterHealthCheck(gvk, checkPodHealth)
}

// checkPodHealth implements health check for Pods
// Based on ArgoCD's gitops-engine implementation
func checkPodHealth(obj *unstructured.Unstructured) bool {
	var pod corev1.Pod
	err := convertFromUnstructured(obj, &pod)
	if err != nil {
		return false
	}

	switch pod.Status.Phase {
	case corev1.PodSucceeded:
		return true
	case corev1.PodRunning:
		// For pods with RestartPolicy Always, check if ready
		if pod.Spec.RestartPolicy == corev1.RestartPolicyAlways {
			// Check if pod is ready
			for _, condition := range pod.Status.Conditions {
				if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
					return true
				}
			}
			// Check for container failures
			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.State.Waiting != nil {
					waiting := containerStatus.State.Waiting
					// Common failure reasons
					if waiting.Reason == "ImagePullBackOff" ||
						waiting.Reason == "ErrImagePull" ||
						waiting.Reason == "CrashLoopBackOff" ||
						waiting.Reason == "CreateContainerConfigError" {
						return false
					}
				}
				if containerStatus.State.Terminated != nil {
					return false
				}
			}
			// Pod is running but not ready yet
			return false
		}
		// For OnFailure/Never restart policies, running means progressing
		return false
	case corev1.PodFailed:
		return false
	case corev1.PodPending:
		return false
	default:
		return false
	}
}
