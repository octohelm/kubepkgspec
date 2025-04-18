package v1alpha1

import (
	"cmp"

	appsv1 "k8s.io/api/apps/v1"
)

type DeployDeployment struct {
	Kind string `json:"kind" validate:"@string{Deployment}"`
	DeployInfrastructure
	Spec DeploymentSpec `json:"spec,omitzero"`
}

func (d DeployDeployment) GetKind() string {
	return cmp.Or(d.Kind, "Deployment")
}

func (d *DeployDeployment) SetKind(kind string) {
	d.Kind = d.GetKind()
}

// +gengo:partialstruct
// +gengo:partialstruct:replace=Template:*PodPartialTemplateSpec json:"template,omitzero"
// +gengo:partialstruct:omit=Selector
type deploymentSpec appsv1.DeploymentSpec
