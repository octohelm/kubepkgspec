package v1alpha1

import (
	"maps"
	"regexp"

	"github.com/octohelm/courier/pkg/validator"
	"github.com/octohelm/courier/pkg/validator/taggedunion"
	"github.com/octohelm/courier/pkg/validator/validators"
)

type DeployInfrastructure struct {
	Labels      map[string]string `json:"labels,omitzero" validate:"@map<@qualified-name,@string[0,63]>"`
	Annotations map[string]string `json:"annotations,omitzero" validate:"@map<@qualified-name,@string[0,4096]>"`
}

func (d DeployInfrastructure) GetLabels() map[string]string {
	return d.Labels
}

func (d *DeployInfrastructure) SetLabels(labels map[string]string) {
	if d.Labels == nil {
		d.Labels = map[string]string{}
	}
	maps.Copy(d.Labels, labels)
}

func init() {
	validator.Register(validator.NewFormatValidatorProvider("qualified-name", func(format string) validator.Validator {
		return &validators.StringValidator{
			Format:        format,
			MinLength:     1,
			Pattern:       regexp.MustCompile(`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]{0,61})?[A-Za-z0-9]$`),
			PatternErrMsg: "必须由字母、数字、'-'、'_' 和 '.' 组成，且以字母或数字开头和结尾，长度不得超过 63，可以带 DNS 子域名前缀和 '/', 如 'example.com/MyName', 'MyName', 'my.name', '123-abc'",
		}
	}))
}

type Deployer interface {
	GetKind() string
	GetLabels() map[string]string
	SetLabels(labels map[string]string)
}

// +gengo:deepcopy=false
type Deploy struct {
	Underlying Deployer `json:"-"`
}

func (in Deploy) IsZero() bool {
	return in.Underlying == nil
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

		"EndpointSlice": Deployer(&DeployEndpointSlice{}),
	}
}
