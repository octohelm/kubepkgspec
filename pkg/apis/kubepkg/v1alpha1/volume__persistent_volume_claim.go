package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type VolumePersistentVolumeClaim struct {
	Type string                                    `json:"type" validate:"@string{PersistentVolumeClaim}"`
	Opt  *corev1.PersistentVolumeClaimVolumeSource `json:"opt,omitzero"`
	Spec corev1.PersistentVolumeClaimSpec          `json:"spec"`

	VolumeMount
}
