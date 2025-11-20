package main

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/crossplane/function-sdk-go/logging"
	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"
)

func TestRunFunction(t *testing.T) {

	type args struct {
		ctx context.Context
		req *fnv1.RunFunctionRequest
	}
	type want struct {
		rsp *fnv1.RunFunctionResponse
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"AutoDetectReadiness": {
			reason: "An existing composed resource with unspecified readiness and a Ready: True status condition should be detected as ready",
			args: args{
				req: &fnv1.RunFunctionRequest{
					Meta: &fnv1.RequestMeta{Tag: "hello"},
					Observed: &fnv1.State{
						Composite: &fnv1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "test.crossplane.io/v1",
								"kind": "TestXR",
								"metadata": {
									"name": "my-test-xr"
								}
							}`),
						},
						Resources: map[string]*fnv1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestComposed",
									"metadata": {
										"name": "my-test-composed"
									},
									"spec": {},
									"status": {
										"conditions": [
											{
												"type": "Ready",
												"status": "True"
											}
										]
									}
								}`),
							},
						},
					},
					Desired: &fnv1.State{
						Resources: map[string]*fnv1.Resource{
							// This function doesn't care about the desired
							// resource schema. In practice it would match
							// observed (without status), but for this test it
							// doesn't matter.
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{}`),
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1.RunFunctionResponse{
					Meta: &fnv1.ResponseMeta{Tag: "hello", Ttl: durationpb.New(response.DefaultTTL)},
					Desired: &fnv1.State{
						Resources: map[string]*fnv1.Resource{
							// This function doesn't care about the desired
							// resource schema. In practice it would match
							// observed (without status), but for this test it
							// doesn't matter.
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{}`),
								Ready:    fnv1.Ready_READY_TRUE,
							},
						},
					},
				},
			},
		},
		"DeploymentHealthCheck": {
			reason: "A Deployment with all replicas ready and Available condition should be detected as ready via health check",
			args: args{
				req: &fnv1.RunFunctionRequest{
					Meta: &fnv1.RequestMeta{Tag: "hello"},
					Observed: &fnv1.State{
						Composite: &fnv1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "test.crossplane.io/v1",
								"kind": "TestXR",
								"metadata": {
									"name": "my-test-xr"
								}
							}`),
						},
						Resources: map[string]*fnv1.Resource{
							"my-deployment": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "apps/v1",
									"kind": "Deployment",
									"metadata": {
										"name": "my-deployment"
									},
									"spec": {
										"replicas": 3
									},
									"status": {
										"updatedReplicas": 3,
										"availableReplicas": 3,
										"conditions": [
											{
												"type": "Available",
												"status": "True"
											}
										]
									}
								}`),
							},
						},
					},
					Desired: &fnv1.State{
						Resources: map[string]*fnv1.Resource{
							// This function doesn't care about the desired
							// resource schema. In practice it would match
							// observed (without status), but for this test it
							// doesn't matter.
							"my-deployment": {
								Resource: resource.MustStructJSON(`{}`),
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1.RunFunctionResponse{
					Meta: &fnv1.ResponseMeta{Tag: "hello", Ttl: durationpb.New(response.DefaultTTL)},
					Desired: &fnv1.State{
						Resources: map[string]*fnv1.Resource{
							// This function doesn't care about the desired
							// resource schema. In practice it would match
							// observed (without status), but for this test it
							// doesn't matter.
							"my-deployment": {
								Resource: resource.MustStructJSON(`{}`),
								Ready:    fnv1.Ready_READY_TRUE,
							},
						},
					},
				},
			},
		},
		"ServiceHealthCheck": {
			reason: "A ClusterIP Service should be immediately ready via health check",
			args: args{
				req: &fnv1.RunFunctionRequest{
					Meta: &fnv1.RequestMeta{Tag: "hello"},
					Observed: &fnv1.State{
						Composite: &fnv1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "test.crossplane.io/v1",
								"kind": "TestXR",
								"metadata": {
									"name": "my-test-xr"
								}
							}`),
						},
						Resources: map[string]*fnv1.Resource{
							"my-service": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "v1",
									"kind": "Service",
									"metadata": {
										"name": "my-service"
									},
									"spec": {
										"type": "ClusterIP"
									}
								}`),
							},
						},
					},
					Desired: &fnv1.State{
						Resources: map[string]*fnv1.Resource{
							// This function doesn't care about the desired
							// resource schema. In practice it would match
							// observed (without status), but for this test it
							// doesn't matter.
							"my-service": {
								Resource: resource.MustStructJSON(`{}`),
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1.RunFunctionResponse{
					Meta: &fnv1.ResponseMeta{Tag: "hello", Ttl: durationpb.New(response.DefaultTTL)},
					Desired: &fnv1.State{
						Resources: map[string]*fnv1.Resource{
							// This function doesn't care about the desired
							// resource schema. In practice it would match
							// observed (without status), but for this test it
							// doesn't matter.
							"my-service": {
								Resource: resource.MustStructJSON(`{}`),
								Ready:    fnv1.Ready_READY_TRUE,
							},
						},
					},
				},
			},
		},
		"FallbackToReadyCondition": {
			reason: "Resources without registered health checks should fall back to Ready condition check",
			args: args{
				req: &fnv1.RunFunctionRequest{
					Meta: &fnv1.RequestMeta{Tag: "hello"},
					Observed: &fnv1.State{
						Composite: &fnv1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "test.crossplane.io/v1",
								"kind": "TestXR",
								"metadata": {
									"name": "my-test-xr"
								}
							}`),
						},
						Resources: map[string]*fnv1.Resource{
							"managed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "rds.aws.crossplane.io/v1alpha1",
									"kind": "DBInstance",
									"metadata": {
										"name": "my-db"
									},
									"spec": {},
									"status": {
										"conditions": [
											{
												"type": "Ready",
												"status": "True"
											}
										]
									}
								}`),
							},
						},
					},
					Desired: &fnv1.State{
						Resources: map[string]*fnv1.Resource{
							// This function doesn't care about the desired
							// resource schema. In practice it would match
							// observed (without status), but for this test it
							// doesn't matter.
							"managed-resource": {
								Resource: resource.MustStructJSON(`{}`),
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1.RunFunctionResponse{
					Meta: &fnv1.ResponseMeta{Tag: "hello", Ttl: durationpb.New(response.DefaultTTL)},
					Desired: &fnv1.State{
						Resources: map[string]*fnv1.Resource{
							// This function doesn't care about the desired
							// resource schema. In practice it would match
							// observed (without status), but for this test it
							// doesn't matter.
							"managed-resource": {
								Resource: resource.MustStructJSON(`{}`),
								Ready:    fnv1.Ready_READY_TRUE,
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := &Function{log: logging.NewNopLogger()}
			rsp, err := f.RunFunction(tc.args.ctx, tc.args.req)

			if diff := cmp.Diff(tc.want.rsp, rsp, protocmp.Transform()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want rsp, +got rsp:\n%s", tc.reason, diff)
			}

			if diff := cmp.Diff(tc.want.err, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want err, +got err:\n%s", tc.reason, diff)
			}
		})
	}
}
