package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
)

type VolumeConfigMap struct {
	Type string                    `json:"type" validate:"@string{ConfigMap}"`
	Opt  *v1.ConfigMapVolumeSource `json:"opt,omitempty"`
	Spec *ConfigMapSpec            `json:"spec,omitempty"`

	VolumeMount
}

type ConfigMapSpec struct {
	Data map[string]string `json:"data"`
}
