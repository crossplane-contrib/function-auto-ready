# Example: Testing function-auto-ready with Kubernetes Resources

This example demonstrates how the function-auto-ready composition function automatically detects readiness for standard Kubernetes resources using resource-specific health checks.

## Files

- **composition-k8s.yaml**: A composition that creates three Kubernetes resources (Service, Deployment, ConfigMap) and uses function-auto-ready to detect their readiness
- **xr.yaml**: A simple composite resource (XR) instance to render the composition
- **functions.yaml**: Function definitions required by the Crossplane CLI
- **observed-k8s.yaml**: Simulated observed resources with status fields, used to test readiness detection with resources that already exist

## What This Example Demonstrates

The composition creates three Kubernetes resources:
- **Service** (ClusterIP type) - Ready immediately since ClusterIP services don't require external provisioning
- **Deployment** (3 replicas, nginx) - Ready when all replicas are available and the Available condition is True
- **ConfigMap** - Ready immediately if it exists (no status conditions to check)

The function-auto-ready function automatically applies resource-specific health checks to determine when each resource is ready.

## Testing Locally

### Basic Rendering (Without Observed Resources)

To see what resources the composition creates:

```shell
crossplane render xr.yaml composition-k8s.yaml functions.yaml
```

This renders the desired state. Resources will not be marked as ready because they don't exist in the observed state yet.

### Rendering with Observed Resources

To simulate resources that already exist in a cluster with their status fields populated:

```shell
crossplane render xr.yaml composition-k8s.yaml functions.yaml -o observed-k8s.yaml
```

This simulates the function behavior after Kubernetes has created the resources. The observed-k8s.yaml file contains:
- Service with empty status (ready immediately)
- Deployment with full status showing 3/3 replicas available and Available condition True (ready)
- ConfigMap with no status (ready immediately)

### Expected Output

When rendered with observed resources, you should see:
- All three composed resources in the desired state
- Each resource marked with `Ready: True` in their annotations
- The composite resource (XR) marked as ready because all composed resources are ready

## See Also

For a simpler example using Crossplane managed resources (AWS S3 Bucket), see the main README.md file in the repository root.
