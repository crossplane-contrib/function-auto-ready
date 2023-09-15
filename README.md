# function-auto-ready

A Function that automatically detects when composed resources are ready. It
considers a composed resource ready if:

* The desired resource appears in the observed resources (i.e. it exists).
* The observed resource has status condition `type: Ready`, `status: True`.

In future this Function may accept input to configure how it should determine
that a composed resource is ready, but for now it's fixed.

## Developing a Function

This Function doesn't use the typical Crossplane build submodule and Makefile,
since we'd like Functions to have a less heavyweight developer experience.
It mostly relies on regular old Go tools:

```shell
# Run code generation - see input/generate.go
$ go generate ./...

# Run tests
$ go test -cover ./...
?       github.com/crossplane/function-auto-ready/input/v1beta1      [no test files]
ok      github.com/crossplane/function-auto-ready    0.006s  coverage: 25.8% of statements

# Lint the code
$ docker run --rm -v $(pwd):/app -v ~/.cache/golangci-lint/v1.54.2:/root/.cache -w /app golangci/golangci-lint:v1.54.2 golangci-lint run

# Build a Docker image - see Dockerfile
$ docker build .
```

This Function can be pushed to any Docker registry. To push to xpkg.upbound.io
use `docker push` and `docker-credential-up` from
https://github.com/upbound/up/.

[Crossplane]: https://crossplane.io
[function-design]: https://github.com/crossplane/crossplane/blob/3996f20/design/design-doc-composition-functions.md
[function-pr]: https://github.com/crossplane/crossplane/pull/4500
[new-crossplane-issue]: https://github.com/crossplane/crossplane/issues/new?assignees=&labels=enhancement&projects=&template=feature_request.md
[install-master-docs]: https://docs.crossplane.io/v1.13/software/install/#install-pre-release-crossplane-versions
[proto-schema]: https://github.com/crossplane/function-sdk-go/blob/main/proto/v1beta1/run_function.proto
[grpcurl]: https://github.com/fullstorydev/grpcurl