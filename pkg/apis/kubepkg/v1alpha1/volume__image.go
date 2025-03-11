package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type VolumeImage struct {
	Type string                    `json:"type" validate:"@string{Image}"`
	Opt  *corev1.ImageVolumeSource `json:"opt,omitzero"`

	VolumeMount
}
