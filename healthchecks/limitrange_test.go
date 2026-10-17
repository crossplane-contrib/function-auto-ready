package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckLimitRangeHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "limitrange - with limits",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "LimitRange",
					"spec": map[string]interface{}{
						"limits": []interface{}{
							map[string]interface{}{
								"type": "Container",
								"default": map[string]interface{}{
									"cpu":    "100m",
									"memory": "128Mi",
								},
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "limitrange - without limits",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "LimitRange",
				},
			},
			expected: true,
		},
		{
			name: "limitrange - with status",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "LimitRange",
					"status": map[string]interface{}{
						"hard": map[string]interface{}{
							"cpu": "10",
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkLimitRangeHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkLimitRangeHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
