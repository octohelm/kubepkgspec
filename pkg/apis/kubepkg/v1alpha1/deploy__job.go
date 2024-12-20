package v1alpha1

import (
	"cmp"
	batchv1 "k8s.io/api/batch/v1"
)

type DeployJob struct {
	Kind string `json:"kind" validate:"@string{Job}"`
	DeployInfrastructure
	Spec JobSpec `json:"spec,omitempty"`
}

func (d DeployJob) GetKind() string {
	return cmp.Or(d.Kind, "Job")
}

// +gengo:partialstruct
// +gengo:partialstruct:replace=Template:*PodPartialTemplateSpec json:"template,omitempty"
// +gengo:partialstruct:omit=Selector
type jobSpec batchv1.JobSpec
