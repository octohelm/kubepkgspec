package v1alpha1

type DeployConfigMap struct {
	Kind        string            `json:"kind" validate:"@string{ConfigMap}"`
	Annotations map[string]string `json:"annotations,omitempty"`
}
