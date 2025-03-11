package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type VolumeHostPath struct {
	Type string                       `json:"type" validate:"@string{HostPath}"`
	Opt  *corev1.HostPathVolumeSource `json:"opt,omitzero"`
	VolumeMount
}
