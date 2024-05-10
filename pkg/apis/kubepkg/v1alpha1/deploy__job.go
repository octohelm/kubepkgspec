package v1alpha1

import (
	batchv1 "k8s.io/api/batch/v1"
)

type DeployJob struct {
	Kind        string            `json:"kind" validate:"@string{Job}"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        JobSpec           `json:"spec,omitempty"`
}

// +gengo:partialstruct
// +gengo:partialstruct:replace=Template:PodPartialTemplateSpec
type jobSpec batchv1.JobSpec
