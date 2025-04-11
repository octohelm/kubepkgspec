package v1alpha1

import (
	"cmp"

	appsv1 "k8s.io/api/apps/v1"
)

type DeployDaemonSet struct {
	Kind string `json:"kind" validate:"@string{DaemonSet}"`
	DeployInfrastructure
	Spec DaemonSetSpec `json:"spec,omitzero"`
}

func (d DeployDaemonSet) GetKind() string {
	return cmp.Or(d.Kind, "DaemonSet")
}

func (d *DeployDaemonSet) SetKind(kind string) {
	d.Kind = d.GetKind()
}

// +gengo:partialstruct
// +gengo:partialstruct:replace=Template:*PodPartialTemplateSpec json:"template,omitzero"
// +gengo:partialstruct:omit=Selector
type daemonSetSpec appsv1.DaemonSetSpec
