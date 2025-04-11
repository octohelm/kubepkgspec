package v1alpha1

import (
	"cmp"

	batchv1 "k8s.io/api/batch/v1"
)

type DeployJob struct {
	Kind string `json:"kind" validate:"@string{Job}"`
	DeployInfrastructure
	Spec JobSpec `json:"spec,omitzero"`
}

func (d DeployJob) GetKind() string {
	return cmp.Or(d.Kind, "Job")
}

func (d *DeployJob) SetKind(kind string) {
	d.Kind = d.GetKind()
}

// +gengo:partialstruct
// +gengo:partialstruct:replace=Template:*PodPartialTemplateSpec json:"template,omitzero"
// +gengo:partialstruct:omit=Selector
type jobSpec batchv1.JobSpec
