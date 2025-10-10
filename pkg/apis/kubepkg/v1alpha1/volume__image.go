package v1alpha1

import (
	"cmp"

	corev1 "k8s.io/api/core/v1"
)

type VolumeImage struct {
	Type string                    `json:"type" validate:"@string{Image}"`
	Opt  *corev1.ImageVolumeSource `json:"opt,omitzero"`

	VolumeMount
}

func (d VolumeImage) GetKind() string {
	return cmp.Or(d.Type, "Image")
}

func (d *VolumeImage) SetKind(kind string) {
	d.Type = d.GetKind()
}
