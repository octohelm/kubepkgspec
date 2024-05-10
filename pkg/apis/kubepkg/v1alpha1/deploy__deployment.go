package v1alpha1

import appsv1 "k8s.io/api/apps/v1"

type DeployDeployment struct {
	Kind        string            `json:"kind" validate:"@string{Deployment}"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        DeploymentSpec    `json:"spec,omitempty"`
}

// +gengo:partialstruct
// +gengo:partialstruct:replace=Template:PodPartialTemplateSpec
// +gengo:partialstruct:omit=Selector
type deploymentSpec appsv1.DeploymentSpec
