package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckCronJobHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy cronjob - suspended",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "CronJob",
					"spec": map[string]interface{}{
						"suspend": true,
					},
				},
			},
			expected: true,
		},
		{
			name: "healthy cronjob - last execution succeeded",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "CronJob",
					"spec": map[string]interface{}{
						"suspend": false,
					},
					"status": map[string]interface{}{
						"lastScheduleTime":    "2024-01-01T10:00:00Z",
						"lastSuccessfulTime": "2024-01-01T10:05:00Z",
					},
				},
			},
			expected: true,
		},
		{
			name: "healthy cronjob - job is active",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "CronJob",
					"spec": map[string]interface{}{
						"suspend": false,
					},
					"status": map[string]interface{}{
						"lastScheduleTime": "2024-01-01T10:00:00Z",
						"active": []interface{}{
							map[string]interface{}{
								"name": "job-1",
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy cronjob - last execution failed",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "CronJob",
					"spec": map[string]interface{}{
						"suspend": false,
					},
					"status": map[string]interface{}{
						"lastScheduleTime":    "2024-01-01T10:05:00Z",
						"lastSuccessfulTime": "2024-01-01T10:00:00Z",
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy cronjob - never succeeded",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "CronJob",
					"spec": map[string]interface{}{
						"suspend": false,
					},
					"status": map[string]interface{}{
						"lastScheduleTime": "2024-01-01T10:00:00Z",
					},
				},
			},
			expected: false,
		},
		{
			name: "progressing cronjob - not yet scheduled",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "batch/v1",
					"kind":       "CronJob",
					"spec": map[string]interface{}{
						"suspend": false,
					},
					"status": map[string]interface{}{},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkCronJobHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkCronJobHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
