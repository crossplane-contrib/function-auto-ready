package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckEventHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "event - normal event",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Event",
					"metadata": map[string]interface{}{
						"name": "event.name",
					},
					"reason": "Killing",
					"message": "Stopped container mycontainer",
				},
			},
			expected: true,
		},
		{
			name: "event - old event",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Event",
					"metadata": map[string]interface{}{
						"name": "event.name",
					},
					"count": int64(5),
					"type":  "Normal",
				},
			},
			expected: true,
		},
		{
			name: "event - no metadata",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Event",
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkEventHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkEventHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
