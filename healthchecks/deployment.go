package healthchecks

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registerDeploymentHealthCheck() {
	gvk := schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
	}
	RegisterHealthCheck(gvk, checkDeploymentHealth)
}

// checkDeploymentHealth implements ArgoCD-style health check for Deployments
// A Deployment is considered healthy when:
// 1. spec.replicas == status.updatedReplicas
// 2. spec.replicas == status.availableReplicas
// 3. status.conditions contains "Available" with status "True"
func checkDeploymentHealth(obj *unstructured.Unstructured) bool {
	// Get spec.replicas (may be nil, defaults to 1)
	specReplicas := int64(1)
	if val, found := getInt64Field(obj.Object, "spec", "replicas"); found {
		specReplicas = val
	}

	// Get status.updatedReplicas
	updatedReplicas, found := getInt64Field(obj.Object, "status", "updatedReplicas")
	if !found {
		return false
	}

	// Get status.availableReplicas
	availableReplicas, found := getInt64Field(obj.Object, "status", "availableReplicas")
	if !found {
		return false
	}

	// Check replica counts match
	if specReplicas != updatedReplicas || specReplicas != availableReplicas {
		return false
	}

	// Check for Available condition
	conditions, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil || !found {
		return false
	}

	for _, cond := range conditions {
		condMap, ok := cond.(map[string]interface{})
		if !ok {
			continue
		}

		condType, found, err := unstructured.NestedString(condMap, "type")
		if err != nil || !found || condType != "Available" {
			continue
		}

		condStatus, found, err := unstructured.NestedString(condMap, "status")
		if err != nil || !found {
			continue
		}

		if condStatus == "True" {
			return true
		}
	}

	return false
}
