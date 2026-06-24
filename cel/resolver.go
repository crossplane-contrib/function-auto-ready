package cel

import (
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/google/cel-go/cel"
	celtypes "github.com/google/cel-go/common/types"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Resolver struct {
	HealthCheckRegistry map[string]string
}

const (
	errCelQueryFailedToCompile           = "failed to compile query"
	errCelQueryReturnTypeNotBool         = "celQuery does not return a bool type"
	errCelQueryFailedToCreateProgram     = "failed to create program from the cel query"
	errCelQueryFailedToEvalProgram       = "failed to eval the program"
	errCelQueryFailedToCreateEnvironment = "cel query failed to create environment"
)

func (r Resolver) GetHealthCheck(gvk schema.GroupVersionKind) (celQuery string, found bool) {
	gvkKey := gvk.Group + "_" + gvk.Version + "_" + gvk.Kind

	celQuery, found = r.HealthCheckRegistry[gvkKey]
	return
}

func (r Resolver) HealthDeriveFromCelQuery(celQuery string, obj map[string]any) (ready resource.Ready, err error) {
	ready = resource.ReadyUnspecified

	env, err := cel.NewEnv(
		cel.Variable("object", cel.AnyType),
	)
	if err != nil {
		err = errors.Wrap(err, errCelQueryFailedToCreateEnvironment)
		return ready, err
	}

	ast, iss := env.Compile(celQuery)
	if iss.Err() != nil {
		err = errors.Wrap(iss.Err(), errCelQueryFailedToCompile)
		return ready, err
	}

	if !ast.OutputType().IsExactType(cel.BoolType) {
		err = errors.New(errCelQueryReturnTypeNotBool)
		return ready, err
	}

	program, err := env.Program(ast)
	if err != nil {
		err = errors.Wrap(err, errCelQueryFailedToCreateProgram)
		return ready, err
	}

	val, _, err := program.Eval(map[string]any{
		"object": obj,
	})
	if err != nil {
		err = errors.Wrap(err, errCelQueryFailedToEvalProgram)
		return ready, err
	}

	if val == celtypes.True {
		ready = resource.ReadyTrue
	} else {
		ready = resource.ReadyFalse
	}
	return ready, err
}
