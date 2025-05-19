package main

import (
	"context"
	"github.com/crossplane/function-auto-ready/input/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"
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
	var exp1 = 1
	var exp2 = 2
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
		"ExpectedCountTrue": {
			reason: "Composite should be marked ready when the expected resource count matches the ready resource count",
			args: args{
				req: &fnv1.RunFunctionRequest{
					Meta:  &fnv1.RequestMeta{Tag: "hello"},
					Input: resource.MustStructObject(&v1beta1.Input{ExpectedResourceCount: &exp1}),
					Observed: &fnv1.State{
						Composite: &fnv1.Resource{
							Ready: fnv1.Ready_READY_UNSPECIFIED,
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
						Composite: &fnv1.Resource{
							Ready:    fnv1.Ready_READY_TRUE,
							Resource: resource.MustStructJSON(`{}`),
						},
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
		"ExpectedCountFalse": {
			reason: "Composite should not be marked ready when the expected resource count does not match the ready resource count",
			args: args{
				req: &fnv1.RunFunctionRequest{
					Meta:  &fnv1.RequestMeta{Tag: "hello"},
					Input: resource.MustStructObject(&v1beta1.Input{ExpectedResourceCount: &exp2}),
					Observed: &fnv1.State{
						Composite: &fnv1.Resource{
							Ready: fnv1.Ready_READY_UNSPECIFIED,
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
						Composite: &fnv1.Resource{
							Ready:    fnv1.Ready_READY_FALSE,
							Resource: resource.MustStructJSON(`{}`),
						},
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
		"InputFromContext": {
			reason: "Function should prioritize input from the context when present",
			args: args{
				req: &fnv1.RunFunctionRequest{
					Meta:    &fnv1.RequestMeta{Tag: "hello"},
					Context: &structpb.Struct{Fields: map[string]*structpb.Value{KeyContext: structpb.NewStructValue(resource.MustStructObject(&v1beta1.Input{ExpectedResourceCount: &exp1}))}},
					Input:   resource.MustStructObject(&v1beta1.Input{ExpectedResourceCount: &exp2}),
					Observed: &fnv1.State{
						Composite: &fnv1.Resource{
							Ready: fnv1.Ready_READY_UNSPECIFIED,
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
					Meta:    &fnv1.ResponseMeta{Tag: "hello", Ttl: durationpb.New(response.DefaultTTL)},
					Context: resource.MustStructJSON(`{"autoready.fn.crossplane.io": {"expectedResourceCount": 1, "metadata": {"generation": 0}}}`),
					Desired: &fnv1.State{
						Composite: &fnv1.Resource{
							Ready:    fnv1.Ready_READY_TRUE,
							Resource: resource.MustStructJSON(`{}`),
						},
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
