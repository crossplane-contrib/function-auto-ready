package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckNodeHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy node - Ready condition True",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Node",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":    "Ready",
								"status":  "True",
								"reason":  "KubeletReady",
								"message": "kubelet is posting ready status",
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy node - Ready condition False",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Node",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":    "Ready",
								"status":  "False",
								"reason":  "KubeletNotReady",
								"message": "container runtime is down",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy node - No Ready condition",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Node",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "MemoryPressure",
								"status": "False",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "progressing node - no status",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Node",
				},
			},
			expected: false,
		},
		{
			name: "unhealthy node - Ready condition False among others",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Node",
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "MemoryPressure",
								"status": "False",
							},
							map[string]interface{}{
								"type":   "DiskPressure",
								"status": "False",
							},
							map[string]interface{}{
								"type":    "Ready",
								"status":  "False",
								"reason":  "KubeletNotReady",
								"message": "container runtime is down",
							},
							map[string]interface{}{
								"type":   "PIDPressure",
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
			result := checkNodeHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkNodeHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
