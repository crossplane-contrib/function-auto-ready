package healthchecks

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestCheckPodHealth(t *testing.T) {
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		expected bool
	}{
		{
			name: "healthy pod - running and ready",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"spec": map[string]interface{}{
						"restartPolicy": "Always",
					},
					"status": map[string]interface{}{
						"phase": "Running",
						"conditions": []interface{}{
							map[string]interface{}{
								"type":   "Ready",
								"status": "True",
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "healthy pod - succeeded",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"spec": map[string]interface{}{
						"restartPolicy": "Never",
					},
					"status": map[string]interface{}{
						"phase": "Succeeded",
					},
				},
			},
			expected: true,
		},
		{
			name: "unhealthy pod - failed",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"spec": map[string]interface{}{
						"restartPolicy": "Always",
					},
					"status": map[string]interface{}{
						"phase": "Failed",
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy pod - pending",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"spec": map[string]interface{}{
						"restartPolicy": "Always",
					},
					"status": map[string]interface{}{
						"phase": "Pending",
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy pod - ImagePullBackOff",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"spec": map[string]interface{}{
						"restartPolicy": "Always",
					},
					"status": map[string]interface{}{
						"phase": "Running",
						"containerStatuses": []interface{}{
							map[string]interface{}{
								"state": map[string]interface{}{
									"waiting": map[string]interface{}{
										"reason": "ImagePullBackOff",
									},
								},
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "unhealthy pod - CrashLoopBackOff",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"spec": map[string]interface{}{
						"restartPolicy": "Always",
					},
					"status": map[string]interface{}{
						"phase": "Running",
						"containerStatuses": []interface{}{
							map[string]interface{}{
								"state": map[string]interface{}{
									"waiting": map[string]interface{}{
										"reason": "CrashLoopBackOff",
									},
								},
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "progressing pod - running with OnFailure restart policy",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"spec": map[string]interface{}{
						"restartPolicy": "OnFailure",
					},
					"status": map[string]interface{}{
						"phase": "Running",
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkPodHealth(tt.obj)
			if result != tt.expected {
				t.Errorf("checkPodHealth() = %v, want %v", result, tt.expected)
			}
		})
	}
}
