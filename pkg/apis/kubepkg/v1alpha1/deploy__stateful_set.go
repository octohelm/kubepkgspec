package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
)

type DeployStatefulSet struct {
	Kind        string                 `json:"kind" validate:"@string{StatefulSet}"`
	Annotations map[string]string      `json:"annotations,omitempty"`
	Spec        appsv1.StatefulSetSpec `json:"spec,omitempty"`
}
