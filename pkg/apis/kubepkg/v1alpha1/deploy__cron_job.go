package v1alpha1

import (
	"cmp"

	batchv1 "k8s.io/api/batch/v1"
)

type DeployCronJob struct {
	Kind string `json:"kind" validate:"@string{CronJob}"`
	DeployInfrastructure
	Spec CronJobSpec `json:"spec,omitzero"`
}

func (d DeployCronJob) GetKind() string {
	return cmp.Or(d.Kind, "CronJob")
}

func (d *DeployCronJob) SetKind(kind string) {
	d.Kind = d.GetKind()
}

// +gengo:partialstruct
// +gengo:partialstruct:replace=JobTemplate:*JobTemplateSpec json:"template,omitzero"
// +gengo:partialstruct:omit=Selector
type cronJobSpec batchv1.CronJobSpec

type JobTemplateSpec struct {
	Spec JobSpec `json:"spec,omitzero"`
}

func (in *JobTemplateSpec) DeepCopyAs() *batchv1.JobTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(batchv1.JobTemplateSpec)
	in.DeepCopyIntoAs(out)
	return out
}

func (t *JobTemplateSpec) DeepCopyIntoAs(templateSpec *batchv1.JobTemplateSpec) {
	t.Spec.DeepCopyIntoAs(&templateSpec.Spec)
}
