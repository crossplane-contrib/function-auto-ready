# function-auto-ready
[![CI](https://github.com/crossplane-contrib/function-auto-ready/actions/workflows/ci.yml/badge.svg)](https://github.com/crossplane-contrib/function-auto-ready/actions/workflows/ci.yml) ![GitHub release (latest SemVer)](https://img.shields.io/github/release/crossplane-contrib/function-auto-ready)

This [composition function][docs-functions] automatically detects composed
resources that are ready. It considers a composed resource ready if:

* Another function added the composed resource to the desired state.
* The composed resource appears in the observed state (i.e. it exists).
* The composed resource has the status condition `type: Ready`, `status: True`.

Crossplane considers a composite resource (XR) to be ready when all of its
desired composed resources are ready.

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

### Composite Readiness
The default behavior is that Crossplane determines the Composite Readiness
based on whether or not all of the Desired Resources are in the Ready state.
This can cause the Composite to appear `Ready` prematurely in some scenarios.
To avoid this problem the `expectedResourceCount` input can be used to ensure
that the Composite resource waits for _at least_ the indicated number of
resources to be ready before the Composite resource is marked as Ready.

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
    input:
      apiVersion: autoready.fn.crossplane.io/v1beta1
      kind: Input
      expectedResourceCount: 1
```

The function will also accept input from the pipeline context using the key
`autoready.fn.crossplane.io`.  Context input overrides inline input values.

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
          ---
          # Other resources...
          ---
          apiVersion: meta.gotemplating.fn.crossplane.io/v1alpha1
          kind: Context
          data:
            "autoready.fn.crossplane.io":
              expectedResourceCount: {{ $.myExpectedResourceCount }}

  - step: automatically-detect-ready-composed-resources
    functionRef:
      name: function-auto-ready
    input:
      expectedResourceCount: 1  # This is ignored because the context has the key
```

See the [example](example) directory for an example you can run locally using
the Crossplane CLI:

```shell
$ crossplane beta render xr.yaml composition.yaml functions.yaml
```

See the [composition functions documentation][docs-functions] to learn more
about `crossplane beta render`.

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