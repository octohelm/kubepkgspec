package v1alpha1

import (
	"os"
	"strings"
	"testing"

	"github.com/octohelm/courier/pkg/validator"
	"github.com/octohelm/x/testing/bdd"
)

func TestKubePkg(t *testing.T) {
	b := bdd.FromT(t)

	b.Given("a simple kube pkg", func(b bdd.T) {
		raw := bdd.Must(os.ReadFile("./testdata/simple.kube-pkg.json"))

		k := &KubePkg{}
		b.Then("unmarshal success",
			bdd.NoError(validator.Unmarshal(raw, k)),
		)
	})

	b.Given("a simple invalid kube pkg", func(b bdd.T) {
		raw := bdd.Must(os.ReadFile("./testdata/simple.kube-pkg.invalid.json"))

		k := &KubePkg{}
		b.Then("unmarshal with error",
			bdd.True(
				strings.Contains(
					validator.Unmarshal(raw, k).Error(),
					"unmarshal JSON number into Go []string at /spec/services/#/paths/http",
				),
			),
		)
	})
}
