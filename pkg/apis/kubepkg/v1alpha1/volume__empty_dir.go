package v1alpha1

import (
	"cmp"

	corev1 "k8s.io/api/core/v1"
)

type VolumeEmptyDir struct {
	Type string                       `json:"type" validate:"@string{EmptyDir}"`
	Opt  *corev1.EmptyDirVolumeSource `json:"opt,omitzero"`

	VolumeMount
}

func (d VolumeEmptyDir) GetKind() string {
	return cmp.Or(d.Type, "EmptyDir")
}

func (d *VolumeEmptyDir) SetKind(kind string) {
	d.Type = d.GetKind()
}
