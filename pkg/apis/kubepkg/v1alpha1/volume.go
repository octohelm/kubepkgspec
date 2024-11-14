package v1alpha1

import (
	"github.com/octohelm/courier/pkg/validator"
	"github.com/octohelm/courier/pkg/validator/taggedunion"
)

// +gengo:deepcopy=false
type Volume struct {
	Underlying VolumeMounter `json:"-"`
}

func (in *Volume) SetUnderlying(u any) {
	in.Underlying = u.(VolumeMounter)
}

func (in *Volume) DeepCopy() *Volume {
	if in == nil {
		return nil
	}
	out := new(Volume)
	in.DeepCopyInto(out)
	return out
}

func (in *Volume) DeepCopyInto(out *Volume) {
	out.Underlying = in.Underlying
}

func (d *Volume) UnmarshalJSON(data []byte) error {
	vv := Volume{}
	if err := taggedunion.Unmarshal(data, &vv); err != nil {
		return err
	}
	*d = vv
	return nil
}

func (d Volume) MarshalJSON() ([]byte, error) {
	if d.Underlying == nil {
		return nil, nil
	}
	return validator.Marshal(d.Underlying)
}

func (Volume) Discriminator() string {
	return "type"
}

func (Volume) Mapping() map[string]any {
	return map[string]any{
		"EmptyDir":              VolumeMounter(&VolumeEmptyDir{}),
		"HostPath":              VolumeMounter(&VolumeHostPath{}),
		"Image":                 VolumeMounter(&VolumeImage{}),
		"PersistentVolumeClaim": VolumeMounter(&VolumePersistentVolumeClaim{}),

		"Secret":    VolumeMounter(&VolumeSecret{}),
		"ConfigMap": VolumeMounter(&VolumeConfigMap{}),
	}
}

type VolumeMounter interface {
	GetVolumeMount() *VolumeMount
}

type VolumeMount struct {
	MountPath string `json:"mountPath"`

	MountPropagation string `json:"mountPropagation,omitempty" validate:"@string{Bidirectional,HostToContainer}"`

	// Prefix mountPath == export, use as envFrom
	Prefix   string `json:"prefix,omitempty"`
	Optional *bool  `json:"optional,omitempty"`

	// else volumeMounts
	ReadOnly bool   `json:"readOnly,omitempty"`
	SubPath  string `json:"subPath,omitempty"`
}

func (vm *VolumeMount) GetVolumeMount() *VolumeMount {
	return vm
}
