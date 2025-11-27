package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckReplicaSetHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy replicaset - all replicas available",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "ReplicaSet",
					"metadata": map[string]interface{}{
						"generation": int64(2),
					},
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"observedGeneration": int64(2),
						"availableReplicas":  int64(3),
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy replicaset - observed generation mismatch",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "ReplicaSet",
					"metadata": map[string]interface{}{
						"generation": int64(3),
					},
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"observedGeneration": int64(2),
						"availableReplicas":  int64(3),
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy replicaset - replica failure",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "ReplicaSet",
					"metadata": map[string]interface{}{
						"generation": int64(1),
					},
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"observedGeneration": int64(1),
						"availableReplicas":  int64(3),
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "ReplicaFailure",
								"status": "True",
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy replicaset - not enough available replicas",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "ReplicaSet",
					"metadata": map[string]interface{}{
						"generation": int64(1),
					},
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"observedGeneration": int64(1),
						"availableReplicas":  int64(2),
					},
				},
			},
			expected: false,
		},
		{
			name: "healthy replicaset - default replicas (1)",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "ReplicaSet",
					"metadata": map[string]interface{}{
						"generation": int64(1),
					},
					"spec": map[string]interface{}{},
					"status": map[string]interface{}{
						"observedGeneration": int64(1),
						"availableReplicas":  int64(1),
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkReplicaSetHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkReplicaSetHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
