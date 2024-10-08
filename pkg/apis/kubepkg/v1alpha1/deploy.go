package v1alpha1

import (
	"encoding/json"

	"github.com/octohelm/courier/pkg/openapi/jsonschema/util"
)

type Deployer interface {
	GetKind() string
}

// +gengo:deepcopy=false
type Deploy struct {
	Underlying Deployer `json:"-"`
}

func (in *Deploy) DeepCopy() *Deploy {
	if in == nil {
		return nil
	}
	out := new(Deploy)
	in.Underlying = out.Underlying
	return out
}

func (in *Deploy) DeepCopyInto(out *Deploy) {
	out.Underlying = in.Underlying
}

func (d *Deploy) UnmarshalJSON(data []byte) error {
	dd := Deploy{}
	if err := util.UnmarshalTaggedUnionFromJSON(data, &dd); err != nil {
		return err
	}
	*d = dd
	return nil
}

func (d Deploy) MarshalJSON() ([]byte, error) {
	if d.Underlying == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(d.Underlying)
}

func (d *Deploy) SetUnderlying(u any) {
	d.Underlying = u.(Deployer)
}

func (Deploy) Discriminator() string {
	return "kind"
}

func (Deploy) Mapping() map[string]any {
	return map[string]any{
		"Deployment":  Deployer(&DeployDeployment{}),
		"DaemonSet":   Deployer(&DeployDaemonSet{}),
		"StatefulSet": Deployer(&DeployStatefulSet{}),
		"Job":         Deployer(&DeployJob{}),
		"CronJob":     Deployer(&DeployCronJob{}),

		"Secret":    Deployer(&DeploySecret{}),
		"ConfigMap": Deployer(&DeployConfigMap{}),

		"Endpoints": Deployer(&DeployEndpoints{}),
	}
}
