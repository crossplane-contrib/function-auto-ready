package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckPersistentVolumeClaimHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy pvc - bound",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "PersistentVolumeClaim",
					"status": map[string]interface{}{
						"phase": "Bound",
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy pvc - lost",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "PersistentVolumeClaim",
					"status": map[string]interface{}{
						"phase": "Lost",
					},
				},
			},
			expected: false,
		},
		{
			name: "progressing pvc - pending",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "PersistentVolumeClaim",
					"status": map[string]interface{}{
						"phase": "Pending",
					},
				},
			},
			expected: false,
		},
		{
			name: "unknown pvc - no status",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "PersistentVolumeClaim",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkPersistentVolumeClaimHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkPersistentVolumeClaimHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
