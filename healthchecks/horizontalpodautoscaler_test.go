package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckHorizontalPodAutoscalerHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy hpa - scaling active",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "autoscaling/v2",
					"kind":       "HorizontalPodAutoscaler",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "ScalingActive",
								"status": "True",
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "healthy hpa - scaling limited",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "autoscaling/v2",
					"kind":       "HorizontalPodAutoscaler",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "ScalingLimited",
								"status": "True",
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy hpa - failed to get scale",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "autoscaling/v2",
					"kind":       "HorizontalPodAutoscaler",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "FailedGetScale",
								"status": "True",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy hpa - failed to update scale",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "autoscaling/v2",
					"kind":       "HorizontalPodAutoscaler",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "FailedUpdateScale",
								"status": "True",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy hpa - invalid selector",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "autoscaling/v2",
					"kind":       "HorizontalPodAutoscaler",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "InvalidSelector",
								"status": "True",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "progressing hpa - no conditions",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "autoscaling/v2",
					"kind":       "HorizontalPodAutoscaler",
					"status":     map[string]interface{}{},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkHorizontalPodAutoscalerHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkHorizontalPodAutoscalerHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
