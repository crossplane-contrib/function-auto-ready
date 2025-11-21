package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckJobHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy job - completed",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "Job",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Complete",
								"status": "True",
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy job - failed",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "Job",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Failed",
								"status": "True",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy job - suspended",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "Job",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Suspended",
								"status": "True",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "progressing job - no conditions",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "Job",
					"status":     map[string]interface{}{},
				},
			},
			expected: false,
		},
		{
			name: "progressing job - condition false",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "Job",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Complete",
								"status": "False",
							},
						},
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkJobHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkJobHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
