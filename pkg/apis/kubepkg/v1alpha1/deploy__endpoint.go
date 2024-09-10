package v1alpha1

import (
	"cmp"
	corev1 "k8s.io/api/core/v1"
)

type DeployEndpoints struct {
	Kind        string                   `json:"kind" validate:"@string{Endpoints}"`
	Annotations map[string]string        `json:"annotations,omitempty"`
	Ports       map[string]int32         `json:"ports"`
	Addresses   []corev1.EndpointAddress `json:"addresses,omitempty"`
}

func (d DeployEndpoints) GetKind() string {
	return cmp.Or(d.Kind, "Endpoints")
}
