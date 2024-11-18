package convert

import (
	"errors"
	"fmt"
	"github.com/octohelm/kubepkgspec/pkg/wellknown"
	"maps"
	"sort"
	"strings"

	"github.com/distribution/reference"
	kubepkgv1alpha1 "github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/object"
	"github.com/octohelm/kubepkgspec/pkg/reload"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeployResourceFrom(kpkg *kubepkgv1alpha1.KubePkg) (object.Object, error) {
	if underlying := kpkg.Spec.Deploy.Underlying; underlying != nil {
		switch x := underlying.(type) {
		case *kubepkgv1alpha1.DeployEndpoints:
			endpoints := &corev1.Endpoints{}
			endpoints.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Endpoints"))
			endpoints.SetNamespace(kpkg.Namespace)
			endpoints.SetName(kpkg.Name)
			endpoints.Labels = maps.Clone(x.Labels)
			endpoints.Annotations = maps.Clone(x.Annotations)

			subset := corev1.EndpointSubset{}
			subset.Addresses = x.Addresses

			portNames := make([]string, 0, len(x.Ports))
			for n := range x.Ports {
				portNames = append(portNames, n)
			}
			sort.Strings(portNames)
			for _, n := range portNames {
				p := corev1.EndpointPort{}
				p.Port = x.Ports[n]
				p.Name, p.Protocol, _ = DecodePortName(n)

				subset.Ports = append(subset.Ports, p)
			}

			endpoints.Subsets = []corev1.EndpointSubset{subset}

			return endpoints, nil
		case *kubepkgv1alpha1.DeployDeployment:
			deployment := &appsv1.Deployment{}
			deployment.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("Deployment"))
			deployment.SetNamespace(kpkg.Namespace)
			deployment.SetName(kpkg.Name)
			deployment.Labels = maps.Clone(x.Labels)
			deployment.Annotations = maps.Clone(x.Annotations)

			podTemplateSpec, err := ToPodTemplateSpec(kpkg)
			if err != nil {
				return nil, err
			}

			deployment.Spec.Template = *podTemplateSpec
			deployment.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: maps.Clone(deployment.Spec.Template.Labels),
			}

			maps.Copy(deployment.Spec.Template.Labels, x.Labels)

			spec, err := Merge(&deployment.Spec, (&x.Spec).DeepCopyAs())
			if err != nil {
				return nil, err
			}
			deployment.Spec = *spec

			return deployment, nil
		case *kubepkgv1alpha1.DeployStatefulSet:
			statefulSet := &appsv1.StatefulSet{}
			statefulSet.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("StatefulSet"))
			statefulSet.SetNamespace(kpkg.Namespace)
			statefulSet.SetName(kpkg.Name)
			statefulSet.Labels = maps.Clone(x.Labels)
			statefulSet.Annotations = maps.Clone(x.Annotations)

			podTemplateSpec, err := ToPodTemplateSpec(kpkg)
			if err != nil {
				return nil, err
			}

			statefulSet.Spec.Template = *podTemplateSpec
			statefulSet.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: maps.Clone(statefulSet.Spec.Template.Labels),
			}

			maps.Copy(statefulSet.Spec.Template.Labels, x.Labels)

			spec, err := Merge(&statefulSet.Spec, (&x.Spec).DeepCopyAs())
			if err != nil {
				return nil, err
			}
			statefulSet.Spec = *spec

			return statefulSet, nil
		case *kubepkgv1alpha1.DeployDaemonSet:
			daemonSet := &appsv1.DaemonSet{}
			daemonSet.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("DaemonSet"))
			daemonSet.SetNamespace(kpkg.Namespace)
			daemonSet.SetName(kpkg.Name)
			daemonSet.Labels = maps.Clone(x.Labels)
			daemonSet.Annotations = maps.Clone(x.Annotations)

			podTemplateSpec, err := ToPodTemplateSpec(kpkg)
			if err != nil {
				return nil, err
			}

			daemonSet.Spec.Template = *podTemplateSpec
			daemonSet.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: maps.Clone(daemonSet.Spec.Template.Labels),
			}

			maps.Copy(daemonSet.Spec.Template.Labels, x.Labels)

			spec, err := Merge(&daemonSet.Spec, (&x.Spec).DeepCopyAs())
			if err != nil {
				return nil, err
			}
			daemonSet.Spec = *spec

			return daemonSet, nil
		case *kubepkgv1alpha1.DeployJob:
			job := &batchv1.Job{}
			job.SetGroupVersionKind(batchv1.SchemeGroupVersion.WithKind("Job"))
			job.SetNamespace(kpkg.Namespace)
			job.SetName(kpkg.Name)
			job.Labels = maps.Clone(x.Labels)
			job.Annotations = maps.Clone(x.Annotations)

			podTemplateSpec, err := ToPodTemplateSpec(kpkg)
			if err != nil {
				return nil, err
			}

			job.Spec.Template = *podTemplateSpec
			job.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: maps.Clone(job.Spec.Template.Labels),
			}

			maps.Copy(job.Spec.Template.Labels, x.Labels)

			spec, err := Merge(&job.Spec, (&x.Spec).DeepCopyAs())
			if err != nil {
				return nil, err
			}
			job.Spec = *spec

			return job, nil
		case *kubepkgv1alpha1.DeployCronJob:
			cronJob := &batchv1.CronJob{}
			cronJob.SetGroupVersionKind(batchv1.SchemeGroupVersion.WithKind("CronJob"))
			cronJob.SetNamespace(kpkg.Namespace)
			cronJob.SetName(kpkg.Name)
			cronJob.Labels = maps.Clone(x.Labels)
			cronJob.Annotations = maps.Clone(x.Annotations)

			podTemplateSpec, err := ToPodTemplateSpec(kpkg)
			if err != nil {
				return nil, err
			}

			cronJob.Spec.JobTemplate.Spec.Template = *podTemplateSpec
			cronJob.Spec.JobTemplate.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: maps.Clone(cronJob.Spec.JobTemplate.Spec.Template.Labels),
			}

			maps.Copy(cronJob.Spec.JobTemplate.Spec.Template.Labels, x.Labels)

			spec, err := Merge(&cronJob.Spec, (&x.Spec).DeepCopyAs())
			if err != nil {
				return nil, err
			}
			cronJob.Spec = *spec

			return cronJob, nil
		case *kubepkgv1alpha1.DeployConfigMap:
			cm := &corev1.ConfigMap{}
			cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
			cm.SetNamespace(kpkg.Namespace)
			cm.SetName(kpkg.Name)
			cm.Labels = maps.Clone(x.Labels)
			cm.Annotations = maps.Clone(x.Annotations)

			configurationType := ""
			if x.Annotations != nil {
				if ct, ok := x.Annotations["configuration.octohelm.tech/type"]; ok {
					configurationType = ct
				}
			}

			data, err := patchConfigEndpoint(configurationType, kpkg.Spec.Config)
			if err != nil {
				return nil, err
			}

			if len(data) == 0 {
				return nil, err
			}

			cm.Data = data

			if err := reload.AnnotateDigestTo(cm, reload.ScopeConfigMapDigest, kpkg.Name, cm.Data); err != nil {
				return nil, err
			}

			return cm, nil
		case *kubepkgv1alpha1.DeploySecret:
			secret := &corev1.Secret{}
			secret.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))
			secret.SetNamespace(kpkg.Namespace)
			secret.SetName(kpkg.Name)
			secret.Labels = maps.Clone(x.Labels)
			secret.Annotations = maps.Clone(x.Annotations)

			configurationType := ""
			if x.Annotations != nil {
				if ct, ok := x.Annotations["configuration.octohelm.tech/type"]; ok {
					configurationType = ct
				}
			}

			data, err := patchConfigEndpoint(configurationType, kpkg.Spec.Config)
			if err != nil {
				return nil, err
			}

			if len(data) == 0 {
				return nil, err
			}

			secret.StringData = data

			if err := reload.AnnotateDigestTo(secret, reload.ScopeSecretDigest, kpkg.Name, secret.StringData); err != nil {
				return nil, err
			}

			return secret, nil
		}
	}

	return nil, nil
}

func ToPodTemplateSpec(kpkg *kubepkgv1alpha1.KubePkg) (*corev1.PodTemplateSpec, error) {
	if len(kpkg.Spec.Containers) == 0 {
		return nil, errors.New("containers should not empty")
	}
	template := &corev1.PodTemplateSpec{}
	if template.Labels == nil {
		template.Labels = map[string]string{}
	}

	template.Labels["app"] = kpkg.Name

	if kpkg.Spec.Deploy.Underlying != nil {
		if err := wellknown.LabelAppInstance.SetTo(kpkg.Spec.Deploy.Underlying, kpkg.Name); err != nil {
			return nil, err
		}
		if err := wellknown.LabelAppVersion.SetTo(kpkg.Spec.Deploy.Underlying, kpkg.Spec.Version); err != nil {
			return nil, err
		}
	}

	initContainerNames := make([]string, 0, len(kpkg.Spec.Containers))
	containerNames := make([]string, 0, len(kpkg.Spec.Containers))
	finalPlatforms := []string{"linux/amd64", "linux/arm64"}

	for name, c := range kpkg.Spec.Containers {
		if platforms := c.Image.Platforms; len(platforms) > 0 {
			finalPlatforms = Intersection(finalPlatforms, platforms)
		}

		if strings.HasPrefix(name, "init-") {
			initContainerNames = append(initContainerNames, name)
			continue
		}

		containerNames = append(containerNames, name)
	}

	sort.Strings(initContainerNames)
	sort.Strings(containerNames)

	for _, name := range initContainerNames {
		c, err := toContainer(kpkg.Spec.Containers[name], name, template, kpkg)
		if err != nil {
			return nil, err
		}
		template.Spec.InitContainers = append(template.Spec.InitContainers, *c)
	}
	for _, name := range containerNames {
		c, err := toContainer(kpkg.Spec.Containers[name], name, template, kpkg)
		if err != nil {
			return nil, err
		}
		template.Spec.Containers = append(template.Spec.Containers, *c)
	}

	if kpkg.Spec.ServiceAccount != nil {
		template.Spec.ServiceAccountName = kpkg.Name
	}

	if len(finalPlatforms) > 0 {
		arch := make([]string, 0, len(finalPlatforms))

		for _, p := range finalPlatforms {
			parts := strings.Split(p, "/")
			if len(parts) >= 2 {
				arch = append(arch, parts[1])
			}
		}

		template.Spec.Affinity = Must(template.Spec.Affinity)
		template.Spec.Affinity.NodeAffinity = Must(template.Spec.Affinity.NodeAffinity)
		template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = Must(template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution)

		// patch only when empty
		if len(template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms) == 0 {
			template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = append(
				template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms,
				corev1.NodeSelectorTerm{
					MatchExpressions: []corev1.NodeSelectorRequirement{{
						Key:      "kubernetes.io/arch",
						Operator: "In",
						Values:   arch,
					}},
				},
			)
		}
	}

	return template, nil
}

func ContainersAndVolumesFromPodTemplateSpec(k *kubepkgv1alpha1.KubePkg, podTemplateSpec *corev1.PodTemplateSpec) error {
	k.Spec.Containers = map[string]kubepkgv1alpha1.Container{}

	k.Spec.Volumes = map[string]kubepkgv1alpha1.Volume{}

	for _, c := range podTemplateSpec.Spec.Volumes {
		v := kubepkgv1alpha1.Volume{}

		if vs := c.EmptyDir; vs != nil {
			v.SetUnderlying(&kubepkgv1alpha1.VolumeEmptyDir{
				Type: "EmptyDir",
				Opt:  vs,
			})

			k.Spec.Volumes[c.Name] = v
			continue
		}

		if vs := c.HostPath; vs != nil {
			v.SetUnderlying(&kubepkgv1alpha1.VolumeHostPath{
				Type: "HostPath",
				Opt:  vs,
			})

			k.Spec.Volumes[c.Name] = v
			continue
		}

		if vs := c.PersistentVolumeClaim; vs != nil {
			v.SetUnderlying(&kubepkgv1alpha1.VolumePersistentVolumeClaim{
				Type: "PersistentVolumeClaim",
				Opt:  vs,
			})
			continue
		}

		if vs := c.ConfigMap; vs != nil {
			v.SetUnderlying(&kubepkgv1alpha1.VolumeConfigMap{
				Type: "ConfigMap",
				Opt:  vs,
			})

			k.Spec.Volumes[c.Name] = v
			continue
		}

		if vs := c.Secret; vs != nil {
			v.SetUnderlying(&kubepkgv1alpha1.VolumeSecret{
				Type: "Secret",
				Opt:  vs,
			})

			k.Spec.Volumes[c.Name] = v
			continue
		}
	}

	resolveEnvFrom := func(c corev1.Container) {
		for _, envSource := range c.EnvFrom {
			if ref := envSource.SecretRef; ref != nil {
				v := kubepkgv1alpha1.Volume{}

				v.SetUnderlying(&kubepkgv1alpha1.VolumeSecret{
					Type: "Secret",
					Opt: &corev1.SecretVolumeSource{
						SecretName: ref.Name,
					},
				})

				k.Spec.Volumes[ref.Name] = v

				v.Underlying.GetVolumeMount().MountPath = "export"
				v.Underlying.GetVolumeMount().Prefix = envSource.Prefix

				continue
			}

			if ref := envSource.ConfigMapRef; ref != nil {
				v := kubepkgv1alpha1.Volume{}

				v.SetUnderlying(&kubepkgv1alpha1.VolumeConfigMap{
					Type: "ConfigMap",
					Opt: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: ref.Name,
						},
					},
				})

				k.Spec.Volumes[ref.Name] = v

				v.Underlying.GetVolumeMount().MountPath = "export"
				v.Underlying.GetVolumeMount().Prefix = envSource.Prefix

				continue
			}
		}
	}

	for _, c := range podTemplateSpec.Spec.InitContainers {
		cc, err := fromContainer(&c)
		if err != nil {
			return err
		}
		k.Spec.Containers[fmt.Sprintf("init-%s", c.Name)] = *cc

		for _, m := range c.VolumeMounts {
			if v, ok := k.Spec.Volumes[m.Name]; ok {
				mountVolume(v.Underlying.GetVolumeMount(), m)
			}
		}

		resolveEnvFrom(c)
	}

	for _, c := range podTemplateSpec.Spec.Containers {
		cc, err := fromContainer(&c)
		if err != nil {
			return err
		}
		k.Spec.Containers[c.Name] = *cc

		for _, m := range c.VolumeMounts {
			if v, ok := k.Spec.Volumes[m.Name]; ok {
				mountVolume(v.Underlying.GetVolumeMount(), m)
			}
		}

		resolveEnvFrom(c)
	}

	return nil
}

func fromContainer(c *corev1.Container) (*kubepkgv1alpha1.Container, error) {
	n, err := reference.Parse(c.Image)
	if err != nil {
		return nil, err
	}

	container := &kubepkgv1alpha1.Container{}

	switch x := n.(type) {
	case reference.NamedTagged:
		container.Image.Name = x.Name()
		container.Image.Tag = x.Tag()
	case reference.Named:
		container.Image.Name = n.String()
	}

	container.Image.PullPolicy = c.ImagePullPolicy

	container.WorkingDir = c.WorkingDir
	container.Command = c.Command
	container.Args = c.Args

	container.Stdin = c.Stdin
	container.StdinOnce = c.StdinOnce
	container.TTY = c.TTY

	container.Resources = &c.Resources

	container.LivenessProbe = c.LivenessProbe
	container.ReadinessProbe = c.ReadinessProbe
	container.StartupProbe = c.StartupProbe
	container.Lifecycle = c.Lifecycle

	container.SecurityContext = c.SecurityContext

	container.TerminationMessagePath = c.TerminationMessagePath
	container.TerminationMessagePolicy = c.TerminationMessagePolicy

	container.Ports = map[string]int32{}

	for _, port := range c.Ports {
		name := FormatPortName(port.Name, port.Protocol, port.HostPort)

		container.Ports[name] = port.ContainerPort
	}

	container.Env = map[string]kubepkgv1alpha1.EnvVarValueOrFrom{}

	for _, env := range c.Env {
		if from := env.ValueFrom; from != nil {
			container.Env[env.Name] = kubepkgv1alpha1.EnvVarValueOrFrom{
				ValueFrom: from,
			}
		} else {
			container.Env[env.Name] = kubepkgv1alpha1.EnvVarValueOrFrom{
				Value: env.Value,
			}
		}
	}

	return container, nil
}

func toContainer(c kubepkgv1alpha1.Container, name string, podTemplateSpec *corev1.PodTemplateSpec, kpkg *kubepkgv1alpha1.KubePkg) (*corev1.Container, error) {
	container := &corev1.Container{}
	container.Name = name

	container.Image = c.Image.FullName()
	container.ImagePullPolicy = c.Image.PullPolicy

	container.WorkingDir = c.WorkingDir
	container.Command = c.Command
	container.Args = c.Args

	container.Stdin = c.Stdin
	container.StdinOnce = c.StdinOnce
	container.TTY = c.TTY

	if resources := c.Resources; resources != nil {
		container.Resources = *resources
	}

	container.LivenessProbe = c.LivenessProbe
	container.ReadinessProbe = c.ReadinessProbe
	container.StartupProbe = c.StartupProbe
	container.Lifecycle = c.Lifecycle

	container.SecurityContext = c.SecurityContext

	container.TerminationMessagePath = c.TerminationMessagePath
	container.TerminationMessagePolicy = c.TerminationMessagePolicy

	envs := make(map[string]kubepkgv1alpha1.EnvVarValueOrFrom)

	for k := range kpkg.Spec.Config {
		if kpkg.Spec.Config[k].ValueFrom != nil {
			envs[k] = kpkg.Spec.Config[k]
		}
	}

	for k := range c.Env {
		envs[k] = c.Env[k]
	}

	envNames := make([]string, 0, len(envs))
	for n := range envs {
		envNames = append(envNames, n)
	}
	sort.Strings(envNames)

	for _, n := range envNames {
		envVar := corev1.EnvVar{}
		envVar.Name = n
		envVarValueOrFrom := envs[n]

		if envVarValueOrFrom.ValueFrom != nil {
			envVar.ValueFrom = envVarValueOrFrom.ValueFrom
		} else {
			envVar.Value = envVarValueOrFrom.Value
		}

		container.Env = append(container.Env, envVar)
	}

	portNames := make([]string, 0, len(c.Ports))
	for n := range c.Ports {
		portNames = append(portNames, n)
	}
	sort.Strings(portNames)

	for _, n := range portNames {
		p := corev1.ContainerPort{}
		p.ContainerPort = c.Ports[n]
		p.Name, p.Protocol, p.HostPort = DecodePortName(n)

		container.Ports = append(container.Ports, p)
	}

	volumeConvertors := VolumeConvertorsFrom(kpkg)

	for n := range volumeConvertors {
		vc := volumeConvertors[n]

		v, err := vc.MountTo(container)
		if err != nil {
			return nil, err
		}

		appendVolumeToPodSpec(podTemplateSpec, v)
	}

	return container, nil
}

func appendVolumeToPodSpec(podTemplateSpec *corev1.PodTemplateSpec, volume *corev1.Volume) {
	if volume == nil {
		return
	}

	added := false
	for _, vol := range podTemplateSpec.Spec.Volumes {
		if vol.Name == volume.Name {
			added = true
			break
		}
	}

	if !added {
		podTemplateSpec.Spec.Volumes = append(podTemplateSpec.Spec.Volumes, *volume)

		if len(podTemplateSpec.Spec.Volumes) > 1 {
			sort.Slice(podTemplateSpec.Spec.Volumes, func(i, j int) bool {
				return podTemplateSpec.Spec.Volumes[i].Name < podTemplateSpec.Spec.Volumes[j].Name
			})
		}
	}
}
