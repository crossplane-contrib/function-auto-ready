package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckIngressHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy ingress - loadbalancer assigned",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "networking.k8s.io/v1",
					"kind":       "Ingress",
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
			name: "healthy ingress - hostname assigned",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "networking.k8s.io/v1",
					"kind":       "Ingress",
					"status": map[string]interface{}{
						"loadBalancer": map[string]interface{}{
							"ingress": []interface{}{
								map[string]interface{}{
									"hostname": "example.com",
								},
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy ingress - no loadbalancer",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "networking.k8s.io/v1",
					"kind":       "Ingress",
					"status": map[string]interface{}{
						"loadBalancer": map[string]interface{}{},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy ingress - no status",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "networking.k8s.io/v1",
					"kind":       "Ingress",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkIngressHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkIngressHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
