package manifest

import (
	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/x/ptr"
	corev1 "k8s.io/api/core/v1"
)

type VolumeConverter interface {
	ToResource(kpkg *v1alpha1.KubePkg) (Object, error)
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
			}
		case *v1alpha1.VolumeSecret:
			volumes[n] = &volumeSecretConverter{
				ResourceName: SubResourceName(kpkg, n),
				VolumeSecret: x,
			}
		case *v1alpha1.VolumeEmptyDir:
			volumes[n] = &volumeEmptyDirConverter{
				ResourceName:   SubResourceName(kpkg, n),
				VolumeEmptyDir: x,
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

func (volumeEmptyDirConverter) ToResource(kpkg *v1alpha1.KubePkg) (Object, error) {
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

		container.VolumeMounts = append(container.VolumeMounts, toVolumeMount(c.ResourceName, c.VolumeMount))

		return v, nil
	}

	return nil, nil
}

type volumeSecretConverter struct {
	ResourceName string
	*v1alpha1.VolumeSecret
}

func (v *volumeSecretConverter) ToResource(kpkg *v1alpha1.KubePkg) (Object, error) {
	secret := &corev1.Secret{}
	secret.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))
	secret.SetNamespace(kpkg.Namespace)
	secret.SetName(v.ResourceName)

	if spec := v.Spec; spec != nil {
		secret.StringData = spec.Data
	}

	if err := AnnotateDigestTo(secret, ScopeSecretDigest, secret.Name, secret.StringData); err != nil {
		return nil, err
	}

	return secret, nil
}

func (c *volumeSecretConverter) MountTo(container *corev1.Container) (*corev1.Volume, error) {
	if c.MountPath == "export" {
		source := corev1.EnvFromSource{
			Prefix: c.Prefix,
		}

		source.SecretRef = &corev1.SecretEnvSource{}
		source.SecretRef.Name = c.ResourceName
		source.SecretRef.Optional = c.Optional

		container.EnvFrom = append(container.EnvFrom, source)

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

	container.VolumeMounts = append(container.VolumeMounts, toVolumeMount(c.ResourceName, c.VolumeMount))
	return v, nil
}

type volumeHostPathConverter struct {
	ResourceName string
	*v1alpha1.VolumeHostPath
}

func (volumeHostPathConverter) ToResource(kpkg *v1alpha1.KubePkg) (Object, error) {
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

		container.VolumeMounts = append(container.VolumeMounts, toVolumeMount(c.ResourceName, c.VolumeMount))

		return v, nil
	}

	return nil, nil
}

type volumePersistentVolumeClaimConverter struct {
	ResourceName string
	*v1alpha1.VolumePersistentVolumeClaim
}

func (v *volumePersistentVolumeClaimConverter) ToResource(kpkg *v1alpha1.KubePkg) (Object, error) {
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

		container.VolumeMounts = append(container.VolumeMounts, toVolumeMount(c.ResourceName, c.VolumeMount))

		return v, nil
	}

	return nil, nil
}

type volumeConfigMapConverter struct {
	ResourceName string
	*v1alpha1.VolumeConfigMap
}

func (c *volumeConfigMapConverter) ToResource(kpkg *v1alpha1.KubePkg) (Object, error) {
	cm := &corev1.ConfigMap{}
	cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
	cm.SetNamespace(kpkg.Namespace)
	cm.SetName(c.ResourceName)

	if spec := c.Spec; spec != nil {
		cm.Data = spec.Data
	}

	if err := AnnotateDigestTo(cm, ScopeConfigMapDigest, cm.Name, cm.Data); err != nil {
		return nil, err
	}

	return cm, nil
}

func (c *volumeConfigMapConverter) MountTo(container *corev1.Container) (*corev1.Volume, error) {
	if c.MountPath == "export" {
		source := corev1.EnvFromSource{
			Prefix: c.Prefix,
		}

		source.ConfigMapRef = &corev1.ConfigMapEnvSource{}
		source.ConfigMapRef.Name = c.ResourceName
		source.ConfigMapRef.Optional = c.Optional

		container.EnvFrom = append(container.EnvFrom, source)

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

	container.VolumeMounts = append(container.VolumeMounts, toVolumeMount(c.ResourceName, c.VolumeMount))

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
	return
}
