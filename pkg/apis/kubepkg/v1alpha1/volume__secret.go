package v1alpha1

import (
	"cmp"
	corev1 "k8s.io/api/core/v1"
)

type VolumeSecret struct {
	Type string                     `json:"type" validate:"@string{Secret}"`
	Opt  *corev1.SecretVolumeSource `json:"opt,omitzero"`
	Spec *ConfigMapSpec             `json:"spec,omitzero"`

	VolumeMount
}

func (d VolumeSecret) GetKind() string {
	return cmp.Or(d.Type, "Secret")
}

func (d *VolumeSecret) SetKind(kind string) {
	d.Type = d.GetKind()
}
