package v1alpha1

import (
	"cmp"
	appsv1 "k8s.io/api/apps/v1"
)

type DeployStatefulSet struct {
	Kind        string            `json:"kind" validate:"@string{StatefulSet}"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        StatefulSetSpec   `json:"spec,omitempty"`
}

func (d DeployStatefulSet) GetKind() string {
	return cmp.Or(d.Kind, "StatefulSet")
}

// +gengo:partialstruct
// +gengo:partialstruct:replace=Template:PodPartialTemplateSpec
// +gengo:partialstruct:omit=Selector
type statefulSetSpec appsv1.StatefulSetSpec
