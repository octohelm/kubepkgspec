package v1alpha1

import (
	"cmp"
	appsv1 "k8s.io/api/apps/v1"
)

type DeployDaemonSet struct {
	Kind string `json:"kind" validate:"@string{DaemonSet}"`
	DeployInfrastructure
	Spec DaemonSetSpec `json:"spec,omitempty"`
}

func (d DeployDaemonSet) GetKind() string {
	return cmp.Or(d.Kind, "DaemonSet")
}

// +gengo:partialstruct
// +gengo:partialstruct:replace=Template:*PodPartialTemplateSpec json:"template,omitempty"
// +gengo:partialstruct:omit=Selector
type daemonSetSpec appsv1.DaemonSetSpec
