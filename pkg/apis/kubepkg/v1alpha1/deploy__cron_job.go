package v1alpha1

import (
	batchv1 "k8s.io/api/batch/v1"
)

type DeployCronJob struct {
	Kind        string              `json:"kind" validate:"@string{CronJob}"`
	Annotations map[string]string   `json:"annotations,omitempty"`
	Spec        batchv1.CronJobSpec `json:"spec,omitempty"`
}
