package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckDeploymentHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy deployment - all replicas ready",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"updatedReplicas":   int64(3),
						"availableReplicas": int64(3),
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Available",
								"status": "True",
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy deployment - replicas not updated",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"updatedReplicas":   int64(2),
						"availableReplicas": int64(3),
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Available",
								"status": "True",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy deployment - replicas not available",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"updatedReplicas":   int64(3),
						"availableReplicas": int64(2),
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Available",
								"status": "True",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy deployment - Available condition False",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"updatedReplicas":   int64(3),
						"availableReplicas": int64(3),
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Available",
								"status": "False",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy deployment - no status",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
				},
			},
			expected: false,
		},
		{
			name: "healthy deployment - default replicas (1)",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"spec":       map[string]interface{}{},
					"status": map[string]interface{}{
						"updatedReplicas":   int64(1),
						"availableReplicas": int64(1),
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Available",
								"status": "True",
							},
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkDeploymentHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkDeploymentHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
