package v1alpha1

import "cmp"

type DeployConfigMap struct {
	Kind        string            `json:"kind" validate:"@string{ConfigMap}"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

func (d DeployConfigMap) GetKind() string {
	return cmp.Or(d.Kind, "ConfigMap")
}
