package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckServiceHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy loadbalancer - ingress assigned",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Service",
					"spec": map[string]interface{}{
						"type": "LoadBalancer",
					},
					"status": map[string]interface{}{
						"loadBalancer": map[string]interface{}{
							"ingress": []interface{}{
								map[string]interface{}{
									"ip": "1.2.3.4",
								},
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy loadbalancer - no ingress",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Service",
					"spec": map[string]interface{}{
						"type": "LoadBalancer",
					},
					"status": map[string]interface{}{
						"loadBalancer": map[string]interface{}{},
					},
				},
			},
			expected: false,
		},
		{
			name: "healthy clusterip - always ready",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Service",
					"spec": map[string]interface{}{
						"type": "ClusterIP",
					},
				},
			},
			expected: true,
		},
		{
			name: "healthy nodeport - always ready",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Service",
					"spec": map[string]interface{}{
						"type": "NodePort",
					},
				},
			},
			expected: true,
		},
		{
			name: "healthy service - default type (ClusterIP)",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Service",
					"spec":       map[string]interface{}{},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkServiceHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkServiceHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
