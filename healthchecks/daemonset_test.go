package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckDaemonSetHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy daemonset - all pods scheduled and ready",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "DaemonSet",
					"status": map[string]interface{}{
						"desiredNumberScheduled": int64(5),
						"numberReady":            int64(5),
						"updatedNumberScheduled": int64(5),
						"numberAvailable":        int64(5),
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy daemonset - pods not ready",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "DaemonSet",
					"status": map[string]interface{}{
						"desiredNumberScheduled": int64(5),
						"numberReady":            int64(4),
						"updatedNumberScheduled": int64(5),
						"numberAvailable":        int64(5),
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy daemonset - update in progress",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "DaemonSet",
					"status": map[string]interface{}{
						"desiredNumberScheduled": int64(5),
						"numberReady":            int64(5),
						"updatedNumberScheduled": int64(4),
						"numberAvailable":        int64(5),
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy daemonset - pods not available",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "DaemonSet",
					"status": map[string]interface{}{
						"desiredNumberScheduled": int64(5),
						"numberReady":            int64(5),
						"updatedNumberScheduled": int64(5),
						"numberAvailable":        int64(4),
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy daemonset - no status",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "DaemonSet",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkDaemonSetHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkDaemonSetHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
