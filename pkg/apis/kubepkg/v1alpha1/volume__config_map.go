package v1alpha1

import (
	"cmp"

	corev1 "k8s.io/api/core/v1"
)

type VolumeConfigMap struct {
	Type string                        `json:"type" validate:"@string{ConfigMap}"`
	Opt  *corev1.ConfigMapVolumeSource `json:"opt,omitzero"`
	Spec *ConfigMapSpec                `json:"spec,omitzero"`

	VolumeMount
}

func (d VolumeConfigMap) GetKind() string {
	return cmp.Or(d.Type, "ConfigMap")
}

func (d *VolumeConfigMap) SetKind(kind string) {
	d.Type = d.GetKind()
}

type ConfigMapSpec struct {
	Data map[string]string `json:"data"`
}
