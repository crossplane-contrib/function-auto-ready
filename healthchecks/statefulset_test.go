package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckStatefulSetHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy statefulset - all replicas ready and updated",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "StatefulSet",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"readyReplicas":   int64(3),
						"currentReplicas": int64(3),
						"currentRevision": "myapp-5d8f9c7b",
						"updateRevision":  "myapp-5d8f9c7b",
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy statefulset - replicas not ready",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "StatefulSet",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"readyReplicas":   int64(2),
						"currentReplicas": int64(3),
						"currentRevision": "myapp-5d8f9c7b",
						"updateRevision":  "myapp-5d8f9c7b",
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy statefulset - update in progress",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "StatefulSet",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"readyReplicas":   int64(3),
						"currentReplicas": int64(3),
						"currentRevision": "myapp-5d8f9c7b",
						"updateRevision":  "myapp-6e9g0d8c",
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy statefulset - wrong current replicas",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "StatefulSet",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
					"status": map[string]interface{}{
						"readyReplicas":   int64(3),
						"currentReplicas": int64(2),
						"currentRevision": "myapp-5d8f9c7b",
						"updateRevision":  "myapp-5d8f9c7b",
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy statefulset - no status",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "StatefulSet",
					"spec": map[string]interface{}{
						"replicas": int64(3),
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkStatefulSetHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkStatefulSetHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
