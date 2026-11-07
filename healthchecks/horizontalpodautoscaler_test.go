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
								"type":   "AbleToScale",
								"status": "False",
								"reason": "FailedGetScale",
							},
							map[string]interface{}{
								"type":   "ScalingActive",
								"status": "True",
								"reason": "ValidMetricFound",
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
								"type":   "AbleToScale",
								"status": "False",
								"reason": "FailedUpdateScale",
							},
							map[string]interface{}{
								"type":   "ScalingActive",
								"status": "True",
								"reason": "ValidMetricFound",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy hpa - failed to get resource metric",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "autoscaling/v2",
					"kind":       "HorizontalPodAutoscaler",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "AbleToScale",
								"status": "True",
								"reason": "SucceededGetScale",
							},
							map[string]interface{}{
								"type":   "ScalingActive",
								"status": "False",
								"reason": "FailedGetResourceMetric",
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
								"type":   "AbleToScale",
								"status": "True",
								"reason": "SucceededGetScale",
							},
							map[string]interface{}{
								"type":   "ScalingActive",
								"status": "False",
								"reason": "InvalidSelector",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy hpa - failing metric next to a stale scaling limited condition",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "autoscaling/v2",
					"kind":       "HorizontalPodAutoscaler",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "AbleToScale",
								"status": "True",
								"reason": "SucceededGetScale",
							},
							map[string]interface{}{
								"type":   "ScalingActive",
								"status": "False",
								"reason": "FailedGetExternalMetric",
							},
							map[string]interface{}{
								"type":   "ScalingLimited",
								"status": "True",
								"reason": "TooFewReplicas",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy hpa - stale scaling limited condition listed before the failing metric",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "autoscaling/v2",
					"kind":       "HorizontalPodAutoscaler",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "ScalingLimited",
								"status": "True",
								"reason": "TooFewReplicas",
							},
							map[string]interface{}{
								"type":   "ScalingActive",
								"status": "False",
								"reason": "FailedGetResourceMetric",
							},
							map[string]interface{}{
								"type":   "AbleToScale",
								"status": "True",
								"reason": "SucceededGetScale",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "healthy hpa - backing off is not a failure",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "autoscaling/v2",
					"kind":       "HorizontalPodAutoscaler",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "AbleToScale",
								"status": "False",
								"reason": "BackoffBoth",
							},
							map[string]interface{}{
								"type":   "ScalingActive",
								"status": "True",
								"reason": "ValidMetricFound",
							},
							map[string]interface{}{
								"type":   "ScalingLimited",
								"status": "False",
								"reason": "DesiredWithinRange",
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "progressing hpa - scaling disabled",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "autoscaling/v2",
					"kind":       "HorizontalPodAutoscaler",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "AbleToScale",
								"status": "True",
								"reason": "SucceededGetScale",
							},
							map[string]interface{}{
								"type":   "ScalingActive",
								"status": "False",
								"reason": "ScalingDisabled",
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
