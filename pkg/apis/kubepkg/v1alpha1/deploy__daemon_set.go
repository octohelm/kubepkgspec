package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
)

type DeployDaemonSet struct {
	Kind        string               `json:"kind" validate:"@string{DaemonSet}"`
	Annotations map[string]string    `json:"annotations,omitempty"`
	Spec        appsv1.DaemonSetSpec `json:"spec,omitempty"`
}
