package v1alpha1

import (
	"fmt"

	"github.com/octohelm/courier/pkg/validator"
	"github.com/octohelm/courier/pkg/validator/taggedunion"
)

type Service struct {
	// Ports [PortName]servicePort
	Ports map[string]int32 `json:"ports,omitzero"`
	// Paths [PortName]PathRuleOrMatch
	Paths map[string]StringOrSlice `json:"paths,omitzero"`
	// ClusterIP
	ClusterIP string `json:"clusterIP,omitzero"`

	Expose *Expose `json:"expose,omitzero"`
}

type Exposer interface {
	GetType() string
}

// +gengo:deepcopy=false
type Expose struct {
	Underlying Exposer `json:"-"`
}

func (e Expose) IsZero() bool {
	return e.Underlying == nil
}

func (in Expose) GetType() string {
	if in.Underlying != nil {
		return in.Underlying.GetType()
	}
	return ""
}

func (in *Expose) SetUnderlying(u any) {
	in.Underlying = u.(Exposer)
}

func (in *Expose) DeepCopy() *Expose {
	if in == nil {
		return nil
	}
	out := new(Expose)
	in.DeepCopyInto(out)
	return out
}

func (in *Expose) DeepCopyInto(out *Expose) {
	out.Underlying = in.Underlying
}

func (d *Expose) UnmarshalJSON(data []byte) error {
	vv := Expose{}
	if err := taggedunion.Unmarshal(data, &vv); err != nil {
		return fmt.Errorf("unmarshal %s failed to Expose: %w", data, err)
	}
	*d = vv
	return nil
}

func (d Expose) MarshalJSON() ([]byte, error) {
	if d.Underlying == nil {
		return nil, nil
	}
	return validator.Marshal(d.Underlying)
}

func (Expose) Discriminator() string {
	return "type"
}

func (Expose) Mapping() map[string]any {
	return map[string]any{
		"NodePort": Exposer(&ExposeNodePort{}),
		"Ingress":  Exposer(&ExposeIngress{}),
	}
}

type ExposeNodePort struct {
	Type string `json:"type" validate:"@string{NodePort}"`
}

func (ExposeNodePort) GetType() string {
	return "NodePort"
}

type ExposeIngress struct {
	Type string `json:"type" validate:"@string{Ingress}"`
	// Gateway
	Gateway StringSlice `json:"gateway,omitzero"`
	// Options
	Options map[string]string `json:"options,omitzero"`
}

func (ExposeIngress) GetType() string {
	return "Ingress"
}
