package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckResourceQuotaHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "resourcequota - with quota used",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ResourceQuota",
					"status": map[string]interface{}{
						"hard": map[string]interface{}{
							"cpu":    "10",
							"memory": "16Gi",
						},
						"used": map[string]interface{}{
							"cpu":    "5",
							"memory": "8Gi",
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "resourcequota - full quota",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ResourceQuota",
					"status": map[string]interface{}{
						"hard": map[string]interface{}{
							"cpu":    "10",
							"memory": "16Gi",
						},
						"used": map[string]interface{}{
							"cpu":    "10",
							"memory": "16Gi",
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "resourcequota - without status",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ResourceQuota",
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkResourceQuotaHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkResourceQuotaHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
