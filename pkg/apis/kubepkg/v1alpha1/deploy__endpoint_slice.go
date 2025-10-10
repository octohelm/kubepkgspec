package v1alpha1

import (
	"cmp"

	discoveryv1 "k8s.io/api/discovery/v1"
)

type DeployEndpointSlice struct {
	Kind string `json:"kind" validate:"@string{EndpointSlice}"`
	DeployInfrastructure

	AddressType discoveryv1.AddressType `json:"addressType"`
	Addresses   []string                `json:"addresses,omitzero"`
	Ports       map[string]int32        `json:"ports"`
}

func (d DeployEndpointSlice) GetKind() string {
	return cmp.Or(d.Kind, "EndpointSlice")
}

func (d *DeployEndpointSlice) SetKind(kind string) {
	d.Kind = d.GetKind()
}
