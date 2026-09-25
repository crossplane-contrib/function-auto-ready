package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckReplicationControllerHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy replicationcontroller - all replicas ready",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ReplicationController",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"replicas":          int64(3),
						"readyReplicas":     int64(3),
						"availableReplicas": int64(3),
					},
				},
			},
			expected: true,
		},
		{
			name: "healthy replicationcontroller - default replicas (1)",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ReplicationController",
					"spec":       map[string]interface{}{},
					"status": map[string]interface{}{
						"replicas":          int64(1),
						"readyReplicas":     int64(1),
						"availableReplicas": int64(1),
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy replicationcontroller - replicas not matching",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ReplicationController",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"replicas":          int64(2),
						"readyReplicas":     int64(2),
						"availableReplicas": int64(2),
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy replicationcontroller - not all available",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ReplicationController",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"replicas":          int64(3),
						"readyReplicas":     int64(3),
						"availableReplicas": int64(2),
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy replicationcontroller - no status",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ReplicationController",
				},
			},
			expected: false,
		},
		{
			name: "unhealthy replicationcontroller - missing readyReplicas",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ReplicationController",
					"status": map[string]interface{}{
						"replicas": int64(1),
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkReplicationControllerHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkReplicationControllerHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
