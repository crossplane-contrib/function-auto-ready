package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckPersistentVolumeHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy pv - bound",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "PersistentVolume",
					"status": map[string]interface{}{
						"phase": "Bound",
					},
				},
			},
			expected: true,
		},
		{
			name: "healthy pv - available",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "PersistentVolume",
					"status": map[string]interface{}{
						"phase": "Available",
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy pv - released",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "PersistentVolume",
					"status": map[string]interface{}{
						"phase": "Released",
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy pv - failed",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "PersistentVolume",
					"status": map[string]interface{}{
						"phase": "Failed",
					},
				},
			},
			expected: false,
		},
		{
			name: "unknown pv - no status",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "PersistentVolume",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkPersistentVolumeHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkPersistentVolumeHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
