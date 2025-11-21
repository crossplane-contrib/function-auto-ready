package healthchecks

import (
	autoscalingv2 "k8s.io/api/autoscaling/v2"
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

	for _, condition := range hpa.Status.Conditions {
		// Check for degraded conditions
		switch condition.Type {
		case "FailedGetScale", "FailedUpdateScale", "FailedGetResourceMetric", "InvalidSelector":
			if condition.Status == "True" {
				return false
			}
		}

		// Check for healthy conditions
		switch condition.Type {
		case autoscalingv2.ScalingActive:
			if condition.Status == "True" {
				return true
			}
		case autoscalingv2.ScalingLimited:
			if condition.Status == "True" {
				return true
			}
		}
	}

	// Progressing (waiting to autoscale)
	return false
}
