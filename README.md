# function-auto-ready
[![CI](https://github.com/crossplane-contrib/function-auto-ready/actions/workflows/ci.yml/badge.svg)](https://github.com/crossplane-contrib/function-auto-ready/actions/workflows/ci.yml) ![GitHub release (latest SemVer)](https://img.shields.io/github/release/crossplane-contrib/function-auto-ready)

This [composition function][docs-functions] automatically detects composed
resources that are ready. It considers a composed resource ready if:

* Another function added the composed resource to the desired state.
* The composed resource appears in the observed state (i.e. it exists).
* **For standard Kubernetes resources** with health check implementations (see list below), the resource passes resource-specific health checks.
* **For all other resources** (Crossplane managed resources, custom resources, etc.), the composed resource has the status condition `type: Ready`, `status: True`.

Crossplane considers a composite resource (XR) to be ready when all of its
desired composed resources are ready.

## Health Checks

This function implements resource-specific health checks for standard Kubernetes resources. The following table shows the current implementation status:

### Core (core/v1)
- [x] Pod - Succeeded, or Running with Ready condition (RestartPolicy: Always)
- [x] Service - ClusterIP/NodePort: immediately ready; LoadBalancer: requires ingress assignment
- [x] Namespace - Always ready if it exists
- [ ] Node
- [x] ConfigMap - Always ready if it exists
- [x] Secret - Always ready if it exists
- [x] ServiceAccount - Always ready if it exists
- [ ] Endpoints
- [ ] PersistentVolume
- [x] PersistentVolumeClaim - Phase is Bound
- [ ] ReplicationController
- [ ] ResourceQuota
- [ ] LimitRange
- [ ] Event

### Apps (apps/v1)
- [x] Deployment - `spec.replicas == status.availableReplicas`, all replicas updated, `Available` condition is `True`
- [x] StatefulSet - `spec.replicas == status.readyReplicas`, all replicas at current revision
- [x] DaemonSet - All desired pods are scheduled, ready, updated, and available
- [x] ReplicaSet - Observed generation matches, available replicas match desired, no replica failures

### Batch (batch/v1)
- [x] Job - Complete condition is True (not Failed or Suspended)
- [x] CronJob - Suspended, has active jobs, or last execution succeeded

### Autoscaling (autoscaling/v2)
- [x] HorizontalPodAutoscaler - ScalingActive or ScalingLimited, no failed conditions

### Networking (networking.k8s.io/v1)
- [x] Ingress - Load balancer ingress is assigned
- [ ] IngressClass
- [ ] NetworkPolicy

### RBAC (rbac.authorization.k8s.io/v1)
- [ ] Role
- [ ] ClusterRole
- [ ] RoleBinding
- [ ] ClusterRoleBinding

### Storage (storage.k8s.io/v1)
- [ ] StorageClass
- [ ] VolumeAttachment
- [ ] CSIDriver
- [ ] CSINode

### Policy (policy/v1)
- [ ] PodDisruptionBudget

For all other resource types (Crossplane managed resources, custom resources, etc.), the function falls back to checking the standard Ready status condition.

## CEL-based health checks (alpha)

Some resource types — notably Crossplane `Configuration` and `Provider`
packages, Cluster API `Cluster` and many more — do not surface a `Ready` status condition, so the default
fallback never considers them ready. To handle these cases the function
supports user-defined readiness expressions written in
[CEL][cel], gated behind the `CELHealthcheckCustomizations` alpha feature
gate.

Enable the feature gate when running the function:

```shell
function-auto-ready --feature-gates=CELHealthcheckCustomizations=true
```

Customizations are supplied as a map keyed by `<group>_<version>_<kind>`
(the group's dots replaced with underscores; the core group is the empty
string). Each value is a CEL expression evaluated against the observed
composed resource, bound to the variable `object`. The expression must
return a boolean — any other result, or an evaluation error, is treated
as not ready. When a customization exists for a given GVK it takes
precedence over the built-in health check.

The map is read from the function's request context. Point the function
at it with the `celHealthCheckCustomizationFrom` input field, which
accepts a [field path][fieldpath]:

```yaml
- step: automatically-detect-ready-composed-resources
  functionRef:
    name: function-auto-ready
  input:
    apiVersion: autoready.fn.crossplane.io/v1alpha1
    kind: Input
    celHealthCheckCustomizationFrom: "[apiextensions.crossplane.io/environment].celHealthCheckCustomizations"
```

Any source that populates the function context can supply the map —
typically `function-environment-configs` reading an `EnvironmentConfig`,
but an earlier pipeline step works just as well. See
[`example/cel-healthcheck`](example/cel-healthcheck/) for a runnable
example that checks the `Installed` and `Healthy` conditions on a
Crossplane `Configuration`.

[cel]: https://github.com/google/cel-spec
[fieldpath]: https://pkg.go.dev/github.com/crossplane/function-sdk-go/request

In this example, the [Go Templating][fn-go-templating] function is used to add
a desired composed resource - an Amazon Web Services S3 Bucket. Once Crossplane
has created the Bucket, the Auto Ready function will let Crossplane know when it
is ready. Because this XR only has one composed resource, the XR will become
ready when the Bucket becomes ready.

```yaml
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: example
spec:
  compositeTypeRef:
    apiVersion: example.crossplane.io/v1beta1
    kind: XR
  mode: Pipeline
  pipeline:
  - step: create-a-bucket
    functionRef:
      name: function-go-templating
    input:
      apiVersion: gotemplating.fn.crossplane.io/v1beta1
      kind: GoTemplate
      source: Inline
      inline:
        template: |
          apiVersion: s3.aws.upbound.io/v1beta1
          kind: Bucket
          metadata:
            annotations:
              gotemplating.fn.crossplane.io/composition-resource-name: bucket
          spec:
            forProvider:
              region: {{ .observed.composite.resource.spec.region }}
  - step: automatically-detect-ready-composed-resources
    functionRef:
      name: function-auto-ready
```

See the [example](example/basic/) directory for an example you can run locally using
the Crossplane CLI:

```shell
$ crossplane render xr.yaml composition.yaml functions.yaml
```

To test with observed resources (simulating resources that already exist):

```shell
$ crossplane render xr.yaml composition-k8s.yaml functions.yaml -o observed-k8s.yaml
```

See the [composition functions documentation][docs-functions] to learn more
about `crossplane render`.

## Developing this function

This function uses [Go][go], [Docker][docker], and the [Crossplane CLI][cli] to
build functions.

```shell
# Run code generation - see input/generate.go
$ go generate ./...

# Run tests - see fn_test.go
$ go test ./...

# Build the function's runtime image - see Dockerfile
$ docker build . --tag=runtime

# Build a function package - see package/crossplane.yaml
$ crossplane xpkg build -f package --embed-runtime-image=runtime
```

[docs-functions]: https://docs.crossplane.io/v1.14/concepts/composition-functions/
[fn-go-templating]: https://github.com/crossplane-contrib/function-go-templating/tree/main
[go]: https://go.dev
[docker]: https://www.docker.com
[cli]: https://docs.crossplane.io/latest/cli