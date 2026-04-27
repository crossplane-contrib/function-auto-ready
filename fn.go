package main

import (
	"context"
	"fmt"
	"maps"
	"regexp"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/crossplane/function-sdk-go/errors"
	"github.com/crossplane/function-sdk-go/logging"
	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"

	"github.com/crossplane/function-auto-ready/cel"
	"github.com/crossplane/function-auto-ready/features"
	"github.com/crossplane/function-auto-ready/healthchecks"
	input "github.com/crossplane/function-auto-ready/input/v1beta1"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
	TTL time.Duration
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1.RunFunctionRequest) (*fnv1.RunFunctionResponse, error) {
	f.log.Debug("Running Function", "tag", req.GetMeta().GetTag())

	rsp := response.To(req, f.TTL)

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

	// Read function input
	in := &input.Input{}
	if err := request.GetInput(req, in); err != nil {
		response.Fatal(rsp, errors.Wrap(err, "invalid input"))
		return rsp, nil
	}

	// First mark resources based on CEL customizations if CELHealthcheckCustomizations alpha feature is enabled
	if features.FeatureGate.Enabled(features.CELHealthcheckCustomizations) {
		// Evaluate the CEL health checks customizations
		// both CELHealthCheckCustomizationFrom and CELHealthCheckCustomization are merged into celHealthChecks
		// with inline CELHealthCheckCustomization taking precedence over customization passed via context
		celHealthchecks := make(map[string]string)

		if in.CELHealthCheckCustomizationFrom != nil {
			// Initialize celHealthchecks with context entries
			celHealthchecks = GetNestedMap(req.GetContext().AsMap(), *in.CELHealthCheckCustomizationFrom)
		}

		if in.CELHealthCheckCustomization != nil {
			// Merge inline cel health checks with existing health checks, overwrite existing entries if they exist
			maps.Copy(celHealthchecks, *in.CELHealthCheckCustomization)
		}

		celResolver := cel.Resolver{
			HealthCheckRegistry: celHealthchecks,
		}

		for name, dr := range desired {
			log := log.WithValues("composed-resource-name", name)

			// Skip if resource doesn't exist yet
			or, ok := observed[name]
			if !ok {
				continue
			}

			// Skip if readiness already explicitly set
			if dr.Ready != resource.ReadyUnspecified {
				continue
			}

			// Check if this resource type has a registered health check customization
			gvk := or.Resource.GroupVersionKind()

			if celQuery, found := celResolver.GetHealthCheck(gvk); found {
				log.Debug("Using resource-specific health check customization", "gvk", gvk.String())
				ready, err := celResolver.HealthDeriveFromCelQuery(celQuery, or.Resource.Object)
				if err != nil {
					response.Warning(rsp, err)
					log.Debug(fmt.Sprintf("Encountered error during resource-specific health check customization evaluation: %s", err.Error()), "gvk", gvk.String())
					continue
				}

				dr.Ready = ready
			}
		}
	}

	// Second, mark standard Kubernetes resources as ready using resource-specific health checks
	for name, dr := range desired {
		log := log.WithValues("composed-resource-name", name)

		// Skip if resource doesn't exist yet
		or, ok := observed[name]
		if !ok {
			continue
		}

		// Skip if readiness already explicitly set
		if dr.Ready != resource.ReadyUnspecified {
			continue
		}

		// Check if this resource type has a registered health check
		// Get GVK from the unstructured object (apiVersion and kind fields)
		// composed.Unstructured embeds unstructured.Unstructured, so we can use it directly
		gvk := or.Resource.GroupVersionKind()
		if healthCheck := healthchecks.GetHealthCheck(gvk); healthCheck != nil {
			log.Debug("Using resource-specific health check", "gvk", gvk.String())
			if healthCheck(&or.Resource.Unstructured) {
				log.Debug("Marked resource as ready via resource-specific health check", "gvk", gvk.String())
				dr.Ready = resource.ReadyTrue
			}
		}
	}

	// Third, check remaining resources using the Ready status condition
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

func GetNestedMap(context map[string]any, key string) map[string]string {
	parts, err := ParseNestedKey(key)
	if err != nil {
		return nil
	}

	currentValue := any(context)
	for _, k := range parts {
		// Check if the current value is a map
		if nestedMap, ok := currentValue.(map[string]any); ok {
			// Get the next value in the nested map
			if nextValue, exists := nestedMap[k]; exists {
				currentValue = nextValue
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	// Convert the final value to a map[string]string
	if resultAny, ok := currentValue.(map[string]any); ok {
		result := make(map[string]string)
		for k, vAny := range resultAny {
			v, ok := vAny.(string)
			if ok {
				result[k] = v
			}
		}
		return result
	}
	return nil
}

// ParseNestedKey enables the bracket and dot notation to key reference
func ParseNestedKey(key string) ([]string, error) {
	var parts []string
	// Regular expression to extract keys, supporting both dot and bracket notation
	regex := regexp.MustCompile(`\[([^\[\]]+)\]|([^.\[\]]+)`)
	matches := regex.FindAllStringSubmatch(key, -1)
	for _, match := range matches {
		if match[1] != "" {
			parts = append(parts, match[1]) // Bracket notation
		} else if match[2] != "" {
			parts = append(parts, match[2]) // Dot notation
		}
	}

	if len(parts) == 0 {
		return nil, errors.New("invalid key")
	}
	return parts, nil
}
