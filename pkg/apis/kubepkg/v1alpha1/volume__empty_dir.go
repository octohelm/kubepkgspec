package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type VolumeEmptyDir struct {
	Type string                       `json:"type" validate:"@string{EmptyDir}"`
	Opt  *corev1.EmptyDirVolumeSource `json:"opt,omitzero"`

	VolumeMount
}
