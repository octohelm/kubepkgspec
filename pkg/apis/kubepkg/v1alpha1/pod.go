package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type PodPartialTemplateSpec struct {
	Spec PodPartialSpec `json:"spec,omitempty"`
}

func (in *PodPartialTemplateSpec) DeepCopyAs() *corev1.PodTemplateSpec {
	if in == nil {
		return nil
	}

	out := new(corev1.PodTemplateSpec)
	in.DeepCopyIntoAs(out)

	return out
}

func (in *PodPartialTemplateSpec) DeepCopyIntoAs(templateSpec *corev1.PodTemplateSpec) {
	if in == nil {
		return
	}

	in.Spec.DeepCopyIntoAs(&templateSpec.Spec)
}

// +gengo:partialstruct
// +gengo:partialstruct:omit=Volumes
// +gengo:partialstruct:omit=InitContainers
// +gengo:partialstruct:omit=Containers
// +gengo:partialstruct:omit=EphemeralContainers
// +gengo:partialstruct:omit=ServiceAccountName
// +gengo:partialstruct:omit=DeprecatedServiceAccount
type podPartialSpec corev1.PodSpec
