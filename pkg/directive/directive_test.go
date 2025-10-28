package directive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/octohelm/kubepkgspec/pkg/directive"
	"github.com/octohelm/x/testing/bdd"
)

func TestApply(t *testing.T) {
	bdd.FromT(t).Given("some directive", func(b bdd.T) {
		b.When("do apply", func(b bdd.T) {
			ok, v, err := directive.Apply(context.Background(), "/generate password")

			b.Then("password generated",
				bdd.True(ok),
				bdd.Equal("generated", v),
				bdd.NoError(err),
			)
		})

		b.When("do apply with single arg", func(b bdd.T) {
			ok, v, err := directive.Apply(context.Background(), "/generate password {n=16}")

			b.Then("password generated",
				bdd.True(ok),
				bdd.Equal("generated{n=16}", v),
				bdd.NoError(err),
			)
		})

		b.When("do apply with multi args", func(b bdd.T) {
			ok, v, err := directive.Apply(context.Background(), `/generate password {n=16,format="1"}`)

			b.Then("password generated",
				bdd.True(ok),
				bdd.Equal(`generated{format="1",n=16}`, v),
				bdd.NoError(err),
			)
		})
	})

	bdd.FromT(t).Given("simple value or invalid directive", func(b bdd.T) {
		b.When("do apply", func(b bdd.T) {
			ok, v, err := directive.Apply(context.Background(), "v")

			b.Then("nothing patched",
				bdd.False(ok),
				bdd.Equal("v", v),
				bdd.NoError(err),
			)
		})
	})

}

func init() {
	directive.Register(PasswordGenerate{})
}

type PasswordGenerate struct {
}

func (PasswordGenerate) Group() string {
	return "generate"
}

func (PasswordGenerate) Name() string {
	return "password"
}

func (PasswordGenerate) Apply(ctx context.Context, arg directive.Args) (string, error) {
	return fmt.Sprintf("generated%s", arg), nil
}
