package healthchecks

import (
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerHorizontalPodAutoscalerHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "autoscaling",
		Version: "v2",
		Kind:    "HorizontalPodAutoscaler",
	}
	RegisterHealthCheck(gvk, checkHorizontalPodAutoscalerHealth)
}

// checkHorizontalPodAutoscalerHealth implements health check for HorizontalPodAutoscalers
// Based on ArgoCD's gitops-engine implementation
func checkHorizontalPodAutoscalerHealth(obj *unstructured.Unstructured) bool {
	var hpa autoscalingv2.HorizontalPodAutoscaler
	err := convertFromUnstructured(obj, &hpa)
	if err != nil {
		return false
	}

	// A failure on any condition makes the HPA unhealthy, wherever it appears
	// in the list. The HPA controller only updates the conditions it
	// evaluated, so a failing ScalingActive can sit next to a ScalingLimited
	// condition that is still True from an earlier, healthy sync.
	for _, condition := range hpa.Status.Conditions {
		if isHorizontalPodAutoscalerDegraded(condition) {
			return false
		}
	}

	for _, condition := range hpa.Status.Conditions {
		switch condition.Type {
		case autoscalingv2.ScalingActive, autoscalingv2.ScalingLimited:
			if condition.Status == corev1.ConditionTrue {
				return true
			}
		}
	}

	// Progressing (waiting to autoscale)
	return false
}

// isHorizontalPodAutoscalerDegraded reports whether the condition says the
// HPA cannot read or update its target's scale, or cannot compute a replica
// count from its metrics. These failures are reported as the reason of an
// AbleToScale or ScalingActive condition whose status is False.
func isHorizontalPodAutoscalerDegraded(condition autoscalingv2.HorizontalPodAutoscalerCondition) bool {
	if condition.Status != corev1.ConditionFalse {
		return false
	}
	switch condition.Type {
	case autoscalingv2.AbleToScale:
		return condition.Reason == "FailedGetScale" || condition.Reason == "FailedUpdateScale"
	case autoscalingv2.ScalingActive:
		// ScalingDisabled means the target was scaled to zero on purpose.
		return condition.Reason != "ScalingDisabled"
	}
	return false
}
