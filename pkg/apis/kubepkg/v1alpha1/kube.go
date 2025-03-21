package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&KubePkg{}, &KubePkgList{})
}

// KubePkgList
// +gengo:deepcopy:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KubePkgList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KubePkg `json:"items"`
}

// KubePkg
// +gengo:deepcopy:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KubePkg struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitzero"`
	Spec              Spec    `json:"spec"`
	Status            *Status `json:"status,omitzero"`
}

func (k *KubePkg) String() string {
	return fmt.Sprintf("%s.%s", k.Name, k.Kind)
}
