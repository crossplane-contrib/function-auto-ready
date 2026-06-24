package features

import (
	"k8s.io/component-base/featuregate"
)

const (
	// CELHealthcheckCustomizations enables the capability to define healthchecks as CEL queries.
	// CEL-based healthchecks overwrite the built-in Go-based healthchecks if defined for one of the supported types.
	CELHealthcheckCustomizations featuregate.Feature = "CELHealthcheckCustomizations"
)

// defaultAutoReadyFeatureGates consists of all known function-auto-ready specific feature keys.
// To add a new feature, define a Feature constant above and add it here with
// its default state and maturity stage (Alpha, Beta, or GA).
var defaultAutoReadyFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	CELHealthcheckCustomizations: {Default: false, PreRelease: featuregate.Alpha},
}

// FeatureGate is the shared global MutableFeatureGate for function-auto-ready.
// It is populated at init time and can be configured via the --feature-gates
// command-line flag in the controller binary.
var FeatureGate featuregate.MutableFeatureGate = featuregate.NewFeatureGate()

func init() {
	if err := FeatureGate.Add(defaultAutoReadyFeatureGates); err != nil {
		panic(err)
	}
}
