package strfmt

import (
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestParseIngressGatewaySet(t *testing.T) {
	internal, err := ParseGatewayTemplate("https://{{ .Name }}---{{ .Namespace }}.internal?always=true")
	testingx.Expect(t, err, testingx.Be[error](nil))

	p, err := ParseGatewayTemplate("https://{{ .Name }}.public")
	testingx.Expect(t, err, testingx.Be[error](nil))

	igs := IngressGatewaySet{
		Gateways: map[string]GatewayTemplate{
			"internal": *internal,
			"public":   *p,
		},
	}

	t.Run("should generate rule", func(t *testing.T) {
		rules := igs.
			For("test", "default").
			IngressRules(map[string]string{
				"http": "/",
			}, "public")

		testingx.Expect(t, len(rules), testingx.Be(2))
	})

	t.Run("should generate custom rule", func(t *testing.T) {
		rules := igs.
			For("test", "default").
			IngressRules(map[string]string{
				"http": "/",
			}, "internal+http://test.internal")

		testingx.Expect(t, rules[0].Host, testingx.Be("test.internal"))
	})

}
