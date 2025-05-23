package main

import (
	"context"

	corev1 "k8s.io/api/core/v1"

	"github.com/crossplane/function-sdk-go/errors"
	"github.com/crossplane/function-sdk-go/logging"
	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1.RunFunctionRequest) (*fnv1.RunFunctionResponse, error) {
	f.log.Debug("Running Function", "tag", req.GetMeta().GetTag())

	rsp := response.To(req, response.DefaultTTL)

	oxr, err := request.GetObservedCompositeResource(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get observed composite resource from %T", req))
		return rsp, nil
	}
	log := f.log.WithValues(
		"xr-apiversion", oxr.Resource.GetAPIVersion(),
		"xr-kind", oxr.Resource.GetKind(),
		"xr-name", oxr.Resource.GetName(),
	)

	observed, err := request.GetObservedComposedResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get observed composed resources from %T", req))
		return rsp, nil
	}

	desired, err := request.GetDesiredComposedResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get desired composed resources from %T", req))
		return rsp, nil
	}

	f.log.Debug("Found desired resources", "count", len(desired))

	// Our goal here is to automatically determine (from the Ready status
	// condition) whether existing composed resources are ready.
	for name, dr := range desired {
		log := log.WithValues("composed-resource-name", name)

		// If this desired resource doesn't exist in the observed resources, it
		// can't be ready because it doesn't yet exist.
		or, ok := observed[name]
		if !ok {
			log.Debug("Ignoring desired resource that does not appear in observed resources")
			continue
		}

		// A previous Function in the pipeline either said this resource was
		// explicitly ready, or explicitly not ready. We only want to
		// automatically determine readiness for desired resources where no
		// other Function has an opinion about their readiness.
		if dr.Ready != resource.ReadyUnspecified {
			log.Debug("Ignoring desired resource that already has explicit readiness", "ready", dr.Ready)
			continue
		}

		// Now we know this resource exists, and no Function that ran before us
		// has an opinion about whether it's ready.

		log.Debug("Found desired resource with unknown readiness")
		// If this observed resource has a status condition with type: Ready,
		// status: True, we set its readiness to true.
		c := or.Resource.GetCondition(xpv1.TypeReady)
		if c.Status == corev1.ConditionTrue {
			log.Debug("Automatically determined that composed resource is ready")
			dr.Ready = resource.ReadyTrue
		}
	}

	if err := response.SetDesiredComposedResources(rsp, desired); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot set desired composed resources from %T", req))
		return rsp, nil
	}

	return rsp, nil
}
