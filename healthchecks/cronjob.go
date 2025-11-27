package healthchecks

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerCronJobHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "batch",
		Version: "v1",
		Kind:    "CronJob",
	}
	RegisterHealthCheck(gvk, checkCronJobHealth)
}

// checkCronJobHealth implements health check for CronJobs
// Based on ArgoCD's gitops-engine implementation
func checkCronJobHealth(obj *unstructured.Unstructured) bool {
	var cronJob batchv1.CronJob
	err := convertFromUnstructured(obj, &cronJob)
	if err != nil {
		return false
	}

	// If suspended, consider it healthy (suspended is a valid state)
	if cronJob.Spec.Suspend != nil && *cronJob.Spec.Suspend {
		return true
	}

	// If there's no status yet, it's progressing
	if cronJob.Status.LastScheduleTime == nil {
		return false
	}

	// If there are active jobs, it's healthy (job is running)
	if len(cronJob.Status.Active) > 0 {
		return true
	}

	// Check if last execution was successful
	// If lastSuccessfulTime is after lastScheduleTime, the last execution succeeded
	if cronJob.Status.LastSuccessfulTime != nil && cronJob.Status.LastScheduleTime != nil {
		if cronJob.Status.LastSuccessfulTime.Time.After(cronJob.Status.LastScheduleTime.Time) ||
			cronJob.Status.LastSuccessfulTime.Time.Equal(cronJob.Status.LastScheduleTime.Time) {
			return true
		}
		// Last execution failed
		return false
	}

	// Never completed successfully
	if cronJob.Status.LastSuccessfulTime == nil {
		return false
	}

	return true
}
