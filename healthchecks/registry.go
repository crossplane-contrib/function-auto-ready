package healthchecks

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// HealthCheckFunc is a function that determines if a Kubernetes resource is healthy/ready.
// It returns true if the resource is ready, false otherwise.
type HealthCheckFunc func(obj *unstructured.Unstructured) bool

// registry holds the mapping from GroupVersionKind to health check functions
var registry = make(map[schema.GroupVersionKind]HealthCheckFunc)

// RegisterHealthCheck registers a health check function for a specific GroupVersionKind
func RegisterHealthCheck(gvk schema.GroupVersionKind, fn HealthCheckFunc) {
	registry[gvk] = fn
}

// GetHealthCheck retrieves the health check function for a specific GroupVersionKind
// Returns nil if no health check is registered for the GVK
func GetHealthCheck(gvk schema.GroupVersionKind) HealthCheckFunc {
	return registry[gvk]
}

// alwaysReady is a health check function for resources that are considered ready
// if they exist in the observed state (e.g., ConfigMap, Secret, ServiceAccount, Namespace).
// These resources don't have meaningful status conditions to check.
func alwaysReady(obj *unstructured.Unstructured) bool {
	return true
}

// convertFromUnstructured converts an unstructured object to a typed Kubernetes object
func convertFromUnstructured(obj *unstructured.Unstructured, target interface{}) error {
	return runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, target)
}

// getInt64Field extracts an int64 value from a field, handling multiple numeric types
func getInt64Field(obj map[string]interface{}, path ...string) (int64, bool) {
	val, found, err := unstructured.NestedFieldNoCopy(obj, path...)
	if err != nil || !found || val == nil {
		return 0, false
	}

	switch v := val.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}

func init() {
	// Register all standard Kubernetes resource health checks
	registerConfigMapHealthCheck()
	registerCronJobHealthCheck()
	registerDaemonSetHealthCheck()
	registerDeploymentHealthCheck()
	registerHorizontalPodAutoscalerHealthCheck()
	registerIngressHealthCheck()
	registerJobHealthCheck()
	registerNamespaceHealthCheck()
	registerPersistentVolumeClaimHealthCheck()
	registerPodHealthCheck()
	registerReplicaSetHealthCheck()
	registerRoleHealthCheck()
	registerRoleBindingHealthCheck()
	registerSecretHealthCheck()
	registerServiceHealthCheck()
	registerServiceAccountHealthCheck()
	registerStatefulSetHealthCheck()
}
