package v1alpha1

type DeploySecret struct {
	Kind        string            `json:"kind" validate:"@string{Secret}"`
	Annotations map[string]string `json:"annotations,omitempty"`
}
