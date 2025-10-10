package v1alpha1

import (
	"cmp"

	corev1 "k8s.io/api/core/v1"
)

type VolumeHostPath struct {
	Type string                       `json:"type" validate:"@string{HostPath}"`
	Opt  *corev1.HostPathVolumeSource `json:"opt,omitzero"`
	VolumeMount
}

func (d VolumeHostPath) GetKind() string {
	return cmp.Or(d.Type, "HostPath")
}

func (d *VolumeHostPath) SetKind(kind string) {
	d.Type = d.GetKind()
}
