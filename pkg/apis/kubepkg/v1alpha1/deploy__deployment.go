package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
)

type DeployDeployment struct {
	Kind        string                `json:"kind" validate:"@string{Deployment}"`
	Annotations map[string]string     `json:"annotations,omitempty"`
	Spec        appsv1.DeploymentSpec `json:"spec,omitempty"`
}
