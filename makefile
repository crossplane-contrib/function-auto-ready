.PHONY: build

build:
	podman build . --tag=runtime --no-cache
	rm -f image.tar
	podman save -o image.tar runtime
	rm -f package/*.xpkg
	crossplane xpkg build -f package --embed-runtime-image-tarball=image.tar
	up xpkg push xpkg.upbound.io/luktom/function-auto-ready:${version} -f package/*.xpkg
	rm -f package/*.xpkg
	rm -f image.tar
