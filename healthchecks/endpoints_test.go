package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckEndpointsHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "endpoints - empty",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Endpoints",
					"metadata": map[string]interface{}{
						"name": "my-endpoints",
					},
				},
			},
			expected: true,
		},
		{
			name: "endpoints - with subsets",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Endpoints",
					"status": map[string]interface{}{
						"subsets": []interface{}{
							map[string]interface{}{
								"addresses": []interface{}{
									map[string]interface{}{
										"ip": "10.0.0.1",
									},
								},
								"ports": []interface{}{
									map[string]interface{}{
										"port": 80,
									},
								},
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "endpoints - no status",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Endpoints",
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkEndpointsHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkEndpointsHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
