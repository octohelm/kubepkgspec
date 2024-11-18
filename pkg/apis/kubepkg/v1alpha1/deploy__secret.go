package v1alpha1

import "cmp"

type DeploySecret struct {
	Kind string `json:"kind" validate:"@string{Secret}"`
	DeployInfrastructure
}

func (d DeploySecret) GetKind() string {
	return cmp.Or(d.Kind, "Secret")
}
