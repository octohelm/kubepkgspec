package v1alpha1

import (
	"encoding/json"
	"github.com/octohelm/courier/pkg/openapi/jsonschema/util"
)

type Service struct {
	// Ports [PortName]servicePort
	Ports map[string]int32 `json:"ports,omitempty"`
	// Paths [PortName]BashPath
	Paths map[string]string `json:"paths,omitempty"`
	// ClusterIP
	ClusterIP string `json:"clusterIP,omitempty"`

	Expose *Expose `json:"expose,omitempty"`
}

type Exposer interface {
	GetType() string
}

// +gengo:deepcopy=false
type Expose struct {
	Underlying Exposer `json:"-"`
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
	if err := util.UnmarshalTaggedUnionFromJSON(data, &vv); err != nil {
		return err
	}
	*d = vv
	return nil
}

func (d Expose) MarshalJSON() ([]byte, error) {
	if d.Underlying == nil {
		return nil, nil
	}
	return json.Marshal(d.Underlying)
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
	Gateway []string `json:"gateway,omitempty"`
}

func (ExposeIngress) GetType() string {
	return "Ingress"
}
