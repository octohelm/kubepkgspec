package v1alpha1

import (
	"github.com/octohelm/courier/pkg/validator"
	"github.com/octohelm/courier/pkg/validator/taggedunion"
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
	if err := taggedunion.Unmarshal(data, &dd); err != nil {
		return err
	}
	*d = dd
	return nil
}

func (d Deploy) MarshalJSON() ([]byte, error) {
	if d.Underlying == nil {
		return []byte("{}"), nil
	}
	return validator.Marshal(d.Underlying)
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
