# Example: CEL-based readiness checks for custom resources

This example shows how to configure `function-auto-ready` to evaluate readiness for
resource types it does not have a built-in health check for, by supplying a
[CEL][cel] expression per GVK.

By default, `function-auto-ready` falls back to checking the standard
`Ready` status condition for any resource type it does not recognize. Some
resources — like Crossplane `Configuration` packages — never surface a
`Ready` condition, and instead report their state via `Installed` and
`Healthy` conditions. For these resources you can provide a CEL expression
that evaluates the observed object and returns a boolean.

## How it works

The composition has three pipeline steps:

1. **`create-k8s-resources`** (`function-go-templating`) — renders a
   `pkg.crossplane.io/v1` `Configuration` as a composed resource.
2. **`fetch-cel-healthcheck-customizations`** (`function-environment-configs`) —
   loads the `healthcheck-customizations` `EnvironmentConfig` and merges its
   `data` into the composition environment under the
   `apiextensions.crossplane.io/environment` context key.
3. **`automatically-detect-ready-composed-resources`** (`function-auto-ready`) —
   reads CEL customizations from the environment via
   `celHealthCheckCustomizationFrom` and uses them when evaluating readiness.

### The CEL customization map

`extra-resources.yaml` defines the customizations keyed by
`<group>_<version>_<kind>` (the group's dots replaced with underscores; the
core group is the empty string):

```yaml
data:
  celHealthCheckCustomizations:
    pkg.crossplane.io_v1_Configuration: >
      object.status.conditions.exists(c, c.type == "Installed" && c.status == "True")
      && object.status.conditions.exists(c, c.type == "Healthy" && c.status == "True")
```

The CEL expression is evaluated against the observed composed resource, which
is bound to the variable `object`. The expression must return a boolean; any
other result, or an evaluation error, is treated as not ready.

### Wiring it into the function

`composition.yaml` points the function at the environment key that holds the
map:

```yaml
- step: automatically-detect-ready-composed-resources
  functionRef:
    name: function-auto-ready
  input:
    apiVersion: autoready.fn.crossplane.io/v1alpha1
    kind: Input
    celHealthCheckCustomizationFrom: "[apiextensions.crossplane.io/environment].celHealthCheckCustomizations"
```

The value of `celHealthCheckCustomizationFrom` is a
[field path][fieldpath] into the function's request context. Any source that
populates that context (environment configs, an earlier function, etc.) can
supply the map.

## Running the example

In a separate shell run the local process of function-auto-ready from root of the repository:
```shell
go run . --insecure --feature-gates=CELHealthcheckCustomizations=true
```

Run render command from example's directory:
```shell
crossplane render \
  --extra-resources extra-resources.yaml \
  --observed-resources observed.yaml \
  --include-context \
  xr.yaml composition.yaml functions.yaml
```

`observed.yaml` simulates a `Configuration` that already exists in the cluster
with both `Installed=True` and `Healthy=True`. Because the supplied CEL
expression evaluates to `true` against that observed state, `function-auto-ready`
marks the composed resource — and therefore the XR — ready.

To see the fallback behavior, edit `observed.yaml` to flip one of the
conditions to `False` and re-run; the XR should no longer be reported as
ready.

[cel]: https://github.com/google/cel-spec
[fieldpath]: https://pkg.go.dev/github.com/crossplane/function-sdk-go/request
