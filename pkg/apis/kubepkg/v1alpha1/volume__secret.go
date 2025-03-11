package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type VolumeSecret struct {
	Type string                     `json:"type" validate:"@string{Secret}"`
	Opt  *corev1.SecretVolumeSource `json:"opt,omitzero"`
	Spec *ConfigMapSpec             `json:"spec,omitzero"`

	VolumeMount
}
