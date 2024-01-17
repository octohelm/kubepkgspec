package v1alpha1

import (
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestEnvVarValueOrFrom(t *testing.T) {
	envVars := []string{
		`@field/metadata.name`,
		`@resource/limit.cpu`,
		`@configMap/xxx-config/X?`,
		`@secret/xxx-config.k3s/X`,
		`just value`,
		`@every 12h`,
	}

	for i := range envVars {
		envVar := envVars[i]

		t.Run(envVar, func(t *testing.T) {
			e2 := EnvVarValueOrFrom{}
			err := e2.UnmarshalText([]byte(envVar))
			testingx.Expect(t, err, testingx.Be[error](nil))

			data, err := e2.MarshalText()
			testingx.Expect(t, err, testingx.Be[error](nil))
			testingx.Expect(t, string(data), testingx.Be(envVar))
		})
	}
}
