package healthchecks

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerJobHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "batch",
		Version: "v1",
		Kind:    "Job",
	}
	RegisterHealthCheck(gvk, checkJobHealth)
}

// checkJobHealth implements health check for Jobs
// Based on ArgoCD's gitops-engine implementation
func checkJobHealth(obj *unstructured.Unstructured) bool {
	var job batchv1.Job
	err := convertFromUnstructured(obj, &job)
	if err != nil {
		return false
	}

	for _, condition := range job.Status.Conditions {
		switch condition.Type {
		case batchv1.JobFailed:
			if condition.Status == "True" {
				return false
			}
		case batchv1.JobSuspended:
			if condition.Status == "True" {
				return false
			}
		case batchv1.JobComplete:
			if condition.Status == "True" {
				return true
			}
		}
	}

	// Job is still progressing
	return false
}
