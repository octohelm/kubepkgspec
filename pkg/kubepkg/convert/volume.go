package convert

import (
	"sort"

	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/object"
	"github.com/octohelm/kubepkgspec/pkg/reloader"
	"github.com/octohelm/x/ptr"
	corev1 "k8s.io/api/core/v1"
)

type VolumeConverter interface {
	ToResource(kpkg *v1alpha1.KubePkg) (object.Object, error)
	MountTo(c *corev1.Container) (*corev1.Volume, error)
}

func VolumeConvertorsFrom(kpkg *v1alpha1.KubePkg) map[string]VolumeConverter {
	volumes := map[string]VolumeConverter{}

	for n, v := range kpkg.Spec.Volumes {
		switch x := v.Underlying.(type) {
		case *v1alpha1.VolumeConfigMap:
			volumes[n] = &volumeConfigMapConverter{
				ResourceName:    SubResourceName(kpkg, n),
				VolumeConfigMap: x,
				External:        IsGlobalRef(n),
			}
		case *v1alpha1.VolumeSecret:
			volumes[n] = &volumeSecretConverter{
				ResourceName: SubResourceName(kpkg, n),
				VolumeSecret: x,
				External:     IsGlobalRef(n),
			}
		case *v1alpha1.VolumeEmptyDir:
			volumes[n] = &volumeEmptyDirConverter{
				ResourceName:   SubResourceName(kpkg, n),
				VolumeEmptyDir: x,
			}

		case *v1alpha1.VolumeImage:
			volumes[n] = &volumeImageConverter{
				ResourceName: SubResourceName(kpkg, n),
				VolumeImage:  x,
			}

		case *v1alpha1.VolumePersistentVolumeClaim:
			volumes[n] = &volumePersistentVolumeClaimConverter{
				ResourceName:                SubResourceName(kpkg, n),
				VolumePersistentVolumeClaim: x,
			}
		case *v1alpha1.VolumeHostPath:
			volumes[n] = &volumeHostPathConverter{
				ResourceName:   SubResourceName(kpkg, n),
				VolumeHostPath: x,
			}
		}
	}

	switch kpkg.Spec.Deploy.Underlying.(type) {
	case *v1alpha1.DeploySecret, *v1alpha1.DeployConfigMap:
		// skip
	default:
		// make config as config
		data := map[string]string{}
		for k, c := range kpkg.Spec.Config {
			if c.ValueFrom == nil {
				data[k] = c.Value
			}
		}
		volumes["#"] = &volumeConfigMapConverter{
			ResourceName: SubResourceName(kpkg, "#"),
			VolumeConfigMap: &v1alpha1.VolumeConfigMap{
				Type: "ConfigMap",
				Spec: &v1alpha1.ConfigMapSpec{
					Data: data,
				},
				VolumeMount: v1alpha1.VolumeMount{
					MountPath: "export",
				},
			},
		}
	}

	return volumes
}

type volumeEmptyDirConverter struct {
	ResourceName string
	*v1alpha1.VolumeEmptyDir
}

func (volumeEmptyDirConverter) ToResource(kpkg *v1alpha1.KubePkg) (object.Object, error) {
	return nil, nil
}

func (c volumeEmptyDirConverter) MountTo(container *corev1.Container) (*corev1.Volume, error) {
	if c.MountPath != "export" {
		v := &corev1.Volume{
			Name: c.ResourceName,
		}

		v.EmptyDir = c.Opt

		if v.EmptyDir == nil {
			v.EmptyDir = &corev1.EmptyDirVolumeSource{}
		}

		sortedAppendVolumeMount(container, toVolumeMount(c.ResourceName, c.VolumeMount))

		return v, nil
	}

	return nil, nil
}

type volumeHostPathConverter struct {
	ResourceName string
	*v1alpha1.VolumeHostPath
}

func (volumeHostPathConverter) ToResource(kpkg *v1alpha1.KubePkg) (object.Object, error) {
	return nil, nil
}

func (c *volumeHostPathConverter) MountTo(container *corev1.Container) (*corev1.Volume, error) {
	if c.MountPath != "export" {
		v := &corev1.Volume{
			Name: c.ResourceName,
		}

		v.HostPath = c.Opt
		if v.HostPath == nil {
			v.HostPath = &corev1.HostPathVolumeSource{}
		}

		sortedAppendVolumeMount(container, toVolumeMount(c.ResourceName, c.VolumeMount))

		return v, nil
	}

	return nil, nil
}

type volumeImageConverter struct {
	ResourceName string
	*v1alpha1.VolumeImage
}

func (volumeImageConverter) ToResource(kpkg *v1alpha1.KubePkg) (object.Object, error) {
	return nil, nil
}

func (c *volumeImageConverter) MountTo(container *corev1.Container) (*corev1.Volume, error) {
	if c.MountPath != "export" {
		v := &corev1.Volume{
			Name: c.ResourceName,
		}

		v.Image = c.Opt

		sortedAppendVolumeMount(container, toVolumeMount(c.ResourceName, c.VolumeMount))

		return v, nil
	}

	return nil, nil
}

type volumePersistentVolumeClaimConverter struct {
	ResourceName string
	*v1alpha1.VolumePersistentVolumeClaim
}

func (v *volumePersistentVolumeClaimConverter) ToResource(kpkg *v1alpha1.KubePkg) (object.Object, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	pvc.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("PersistentVolumeClaim"))
	pvc.SetNamespace(kpkg.Namespace)
	pvc.SetName(v.ResourceName)

	spec, err := Merge(&pvc.Spec, &v.Spec)
	if err != nil {
		return nil, err
	}
	pvc.Spec = *spec

	return pvc, nil
}

func (c volumePersistentVolumeClaimConverter) MountTo(container *corev1.Container) (*corev1.Volume, error) {
	if c.MountPath != "export" {
		v := &corev1.Volume{
			Name: c.ResourceName,
		}

		v.PersistentVolumeClaim = c.Opt
		if v.PersistentVolumeClaim == nil {
			v.PersistentVolumeClaim = &corev1.PersistentVolumeClaimVolumeSource{}
		}
		v.PersistentVolumeClaim.ClaimName = c.ResourceName

		sortedAppendVolumeMount(container, toVolumeMount(c.ResourceName, c.VolumeMount))

		return v, nil
	}

	return nil, nil
}

type volumeConfigMapConverter struct {
	ResourceName string
	External     bool
	*v1alpha1.VolumeConfigMap
}

func (v *volumeConfigMapConverter) IsZero() bool {
	if v.Spec == nil {
		return true
	}
	return len(v.Spec.Data) == 0
}

func (v *volumeConfigMapConverter) IsOptional() bool {
	if v.Opt != nil && v.Opt.Optional != nil {
		return *v.Opt.Optional
	}
	return false
}

func (c *volumeConfigMapConverter) ToResource(kpkg *v1alpha1.KubePkg) (object.Object, error) {
	if c.External || c.IsZero() {
		return nil, nil
	}

	cm := &corev1.ConfigMap{}
	cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
	cm.SetNamespace(kpkg.Namespace)
	cm.SetName(c.ResourceName)

	if spec := c.Spec; spec != nil {
		cm.Data = spec.Data
	}

	if err := reloader.AnnotateDigestTo(cm, cm.Data); err != nil {
		return nil, err
	}

	return cm, nil
}

func (c *volumeConfigMapConverter) MountTo(container *corev1.Container) (*corev1.Volume, error) {
	if c.MountPath == "export" {
		if c.IsZero() {
			return nil, nil
		}

		source := corev1.EnvFromSource{
			Prefix: c.Prefix,
		}

		source.ConfigMapRef = &corev1.ConfigMapEnvSource{}
		source.ConfigMapRef.Name = c.ResourceName
		source.ConfigMapRef.Optional = c.Optional

		container.EnvFrom = append(container.EnvFrom, source)

		return nil, nil
	}

	if !c.External && c.IsZero() && !c.IsOptional() {
		return nil, nil
	}

	v := &corev1.Volume{
		Name: c.ResourceName,
	}
	v.ConfigMap = c.Opt
	if v.ConfigMap == nil {
		v.ConfigMap = &corev1.ConfigMapVolumeSource{}
	}
	v.ConfigMap.Name = c.ResourceName

	sortedAppendVolumeMount(container, toVolumeMount(c.ResourceName, c.VolumeMount))

	return v, nil
}

type volumeSecretConverter struct {
	ResourceName string
	External     bool
	*v1alpha1.VolumeSecret
}

func (v *volumeSecretConverter) IsZero() bool {
	if v.Spec == nil {
		return true
	}
	return len(v.Spec.Data) == 0
}

func (v *volumeSecretConverter) IsOptional() bool {
	if v.Opt != nil && v.Opt.Optional != nil {
		return *v.Opt.Optional
	}
	return false
}

func (c *volumeSecretConverter) ToResource(kpkg *v1alpha1.KubePkg) (object.Object, error) {
	if c.External || c.IsZero() {
		return nil, nil
	}

	secret := &corev1.Secret{}
	secret.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))
	secret.SetNamespace(kpkg.Namespace)
	secret.SetName(c.ResourceName)

	if spec := c.Spec; spec != nil {
		secret.StringData = spec.Data
	}

	if err := reloader.AnnotateDigestTo(secret, secret.StringData); err != nil {
		return nil, err
	}

	return secret, nil
}

func (c *volumeSecretConverter) MountTo(container *corev1.Container) (*corev1.Volume, error) {
	if c.MountPath == "export" {
		if c.IsZero() {
			return nil, nil
		}

		source := corev1.EnvFromSource{
			Prefix: c.Prefix,
		}

		source.SecretRef = &corev1.SecretEnvSource{}
		source.SecretRef.Name = c.ResourceName
		source.SecretRef.Optional = c.Optional

		container.EnvFrom = append(container.EnvFrom, source)

		return nil, nil
	}

	if !c.External && c.IsZero() && !c.IsOptional() {
		return nil, nil
	}

	v := &corev1.Volume{
		Name: c.ResourceName,
	}
	v.Secret = c.Opt
	if v.Secret == nil {
		v.Secret = &corev1.SecretVolumeSource{}
	}
	v.Secret.SecretName = c.ResourceName

	sortedAppendVolumeMount(container, toVolumeMount(c.ResourceName, c.VolumeMount))

	return v, nil
}

func toVolumeMount(name string, v v1alpha1.VolumeMount) (vv corev1.VolumeMount) {
	vv.Name = name
	vv.MountPath = v.MountPath
	vv.SubPath = v.SubPath
	vv.ReadOnly = v.ReadOnly

	switch corev1.MountPropagationMode(v.MountPropagation) {
	case corev1.MountPropagationBidirectional:
		vv.MountPropagation = ptr.Ptr(corev1.MountPropagationBidirectional)
	case corev1.MountPropagationHostToContainer:
		vv.MountPropagation = ptr.Ptr(corev1.MountPropagationHostToContainer)
	}
	return vv
}

func mountVolume(v *v1alpha1.VolumeMount, vm corev1.VolumeMount) {
	v.MountPath = vm.MountPath
	v.SubPath = vm.SubPath
	v.ReadOnly = vm.ReadOnly

	if mode := vm.MountPropagation; mode != nil {
		v.MountPropagation = string(*mode)
	}

	return
}

func sortedAppendVolumeMount(container *corev1.Container, vm corev1.VolumeMount) {
	container.VolumeMounts = append(container.VolumeMounts, vm)

	if len(container.VolumeMounts) > 1 {
		sort.Slice(container.VolumeMounts, func(i, j int) bool {
			return container.VolumeMounts[i].Name < container.VolumeMounts[j].Name
		})
	}
}
