package strfmt

import (
	"bytes"
	"html/template"
	"io"
	"net/url"
	"strings"

	"github.com/octohelm/x/ptr"
	"github.com/pkg/errors"
	networkingv1 "k8s.io/api/networking/v1"
)

func From(templates []GatewayTemplate) *IngressGatewaySet {
	s := &IngressGatewaySet{
		Gateways: map[string]GatewayTemplate{},
	}

	for i := range templates {
		t := templates[i]
		s.Gateways[t.Name] = t
	}

	return s
}

type IngressGatewaySet struct {
	Gateways map[string]GatewayTemplate `json:"gateways"`

	serviceName string
	namespace   string
}

func (s IngressGatewaySet) IsZero() bool {
	return len(s.Gateways) == 0
}

func (s IngressGatewaySet) For(service string, namespace string) *IngressGatewaySet {
	s.serviceName = service
	s.namespace = namespace
	return &s
}

func (s *IngressGatewaySet) Endpoints() map[string]string {
	endpoints := map[string]string{
		"default": "http://" + s.serviceName,
	}

	for name, gt := range s.Gateways {
		host := s.hostFor(&gt)

		if host == "" {
			endpoints[name] = ""
		} else {
			if gt.Https {
				endpoints[name] = "https://" + host
			} else {
				endpoints[name] = "http://" + host
			}
		}
	}

	return endpoints
}

func (s *IngressGatewaySet) hostFor(gt *GatewayTemplate) string {
	b := bytes.NewBuffer(nil)
	_ = gt.Execute(b, map[string]any{
		"Name":      s.serviceName,
		"Namespace": s.namespace,
	})
	return b.String()
}

func (s *IngressGatewaySet) IngressRules(paths map[string]string, gateways ...string) (rules []networkingv1.IngressRule) {
	gatewaySet := map[string]GatewayTemplate{}

	for name, g := range s.Gateways {
		if g.AlwaysEnabled {
			gatewaySet[name] = g
		}
	}

	for i := range gateways {
		gateway := gateways[i]

		if gt, ok := s.Gateways[gateway]; ok {
			gatewaySet[gateway] = gt
		}

		if v := strings.Index(gateway, "+"); v > 0 {
			name := gateway[0:v]
			t, err := ParseGatewayTemplate(gateway[v+1:])
			if err == nil {
				gatewaySet[name] = *t
			}
		}
	}

	for _, gt := range gatewaySet {
		for portName, p := range paths {
			r := networkingv1.IngressRule{}
			r.Host = s.hostFor(&gt)
			r.HTTP = &networkingv1.HTTPIngressRuleValue{
				Paths: []networkingv1.HTTPIngressPath{
					{
						Path:     p,
						PathType: ptr.Ptr(networkingv1.PathTypeImplementationSpecific),
						Backend: networkingv1.IngressBackend{
							Service: &networkingv1.IngressServiceBackend{
								Name: s.serviceName,
								Port: networkingv1.ServiceBackendPort{
									Name: portName,
								},
							},
						},
					},
				},
			}
			rules = append(rules, r)
		}
	}

	return
}

// openapi:strfmt gateway-template
type GatewayTemplate struct {
	Name          string
	HostTemplate  string
	Https         bool
	AlwaysEnabled bool
}

func (t *GatewayTemplate) UnmarshalText(d []byte) error {
	tt, err := ParseGatewayTemplate(string(d))
	if err != nil {
		return err
	}
	*t = *tt
	return nil
}

func (t GatewayTemplate) MarshalText() ([]byte, error) {
	s := bytes.NewBuffer(nil)

	if t.Https {
		s.WriteString("https://")
	} else {
		s.WriteString("http://")
	}

	s.WriteString(t.HostTemplate)

	if t.AlwaysEnabled {
		s.WriteString("?always=true")
	}

	return s.Bytes(), nil
}

func (t *GatewayTemplate) Execute(w io.Writer, m map[string]any) error {
	tt, err := template.New(t.HostTemplate).Parse(t.HostTemplate)
	if err != nil {
		return err
	}
	return tt.Execute(w, m)
}

func ParseGatewayTemplate(t string) (*GatewayTemplate, error) {
	parts := strings.Split(t, "://")
	if len(parts) == 2 {
		gt := &GatewayTemplate{}
		gt.Name = strings.Split(parts[0], "+")[0]
		gt.Https = strings.HasSuffix(parts[0], "https")

		hostAndParams := strings.Split(parts[1], "?")

		if len(hostAndParams) > 1 {
			params, err := url.ParseQuery(hostAndParams[1])
			if err != nil {
				return nil, err
			}

			if params.Get("always") == "true" {
				gt.AlwaysEnabled = true
			}
		}

		gt.HostTemplate = hostAndParams[0]
		return gt, nil
	}
	return nil, errors.Errorf("invalid gateway template %s", t)
}
