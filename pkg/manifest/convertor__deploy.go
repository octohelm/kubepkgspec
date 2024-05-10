package manifest

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/manifest/object"
)

func DeployResourceFrom(kpkg *v1alpha1.KubePkg) (object.Object, error) {
	if underlying := kpkg.Spec.Deploy.Underlying; underlying != nil {
		switch x := underlying.(type) {
		case *v1alpha1.DeployDeployment:
			deployment := &appsv1.Deployment{}
			deployment.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("Deployment"))
			deployment.SetNamespace(kpkg.Namespace)
			deployment.SetName(kpkg.Name)

			podTemplateSpec, err := toPodTemplateSpec(kpkg)
			if err != nil {
				return nil, err
			}

			deployment.Spec.Template = *podTemplateSpec
			deployment.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: deployment.Spec.Template.Labels,
			}

			spec, err := Merge(&deployment.Spec, (&x.Spec).DeepCopyAs())
			if err != nil {
				return nil, err
			}
			deployment.Spec = *spec

			return deployment, nil
		case *v1alpha1.DeployStatefulSet:
			statefulSet := &appsv1.StatefulSet{}
			statefulSet.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("StatefulSet"))
			statefulSet.SetNamespace(kpkg.Namespace)
			statefulSet.SetName(kpkg.Name)

			podTemplateSpec, err := toPodTemplateSpec(kpkg)
			if err != nil {
				return nil, err
			}

			statefulSet.Spec.Template = *podTemplateSpec
			statefulSet.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: statefulSet.Spec.Template.Labels,
			}

			spec, err := Merge(&statefulSet.Spec, (&x.Spec).DeepCopyAs())
			if err != nil {
				return nil, err
			}
			statefulSet.Spec = *spec

			return statefulSet, nil
		case *v1alpha1.DeployDaemonSet:
			daemonSet := &appsv1.DaemonSet{}
			daemonSet.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("DaemonSet"))
			daemonSet.SetNamespace(kpkg.Namespace)
			daemonSet.SetName(kpkg.Name)

			podTemplateSpec, err := toPodTemplateSpec(kpkg)
			if err != nil {
				return nil, err
			}

			daemonSet.Spec.Template = *podTemplateSpec
			daemonSet.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: daemonSet.Spec.Template.Labels,
			}

			spec, err := Merge(&daemonSet.Spec, (&x.Spec).DeepCopyAs())
			if err != nil {
				return nil, err
			}
			daemonSet.Spec = *spec

			return daemonSet, nil
		case *v1alpha1.DeployJob:
			job := &batchv1.Job{}
			job.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("Job"))
			job.SetNamespace(kpkg.Namespace)
			job.SetName(kpkg.Name)

			podTemplateSpec, err := toPodTemplateSpec(kpkg)
			if err != nil {
				return nil, err
			}

			job.Spec.Template = *podTemplateSpec
			job.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: job.Spec.Template.Labels,
			}

			spec, err := Merge(&job.Spec, (&x.Spec).DeepCopyAs())
			if err != nil {
				return nil, err
			}
			job.Spec = *spec

			return job, nil
		case *v1alpha1.DeployCronJob:
			cronJob := &batchv1.CronJob{}
			cronJob.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("CronJob"))
			cronJob.SetNamespace(kpkg.Namespace)
			cronJob.SetName(kpkg.Name)

			podTemplateSpec, err := toPodTemplateSpec(kpkg)
			if err != nil {
				return nil, err
			}

			cronJob.Spec.JobTemplate.Spec.Template = *podTemplateSpec
			cronJob.Spec.JobTemplate.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: cronJob.Spec.JobTemplate.Spec.Template.Labels,
			}

			spec, err := Merge(&cronJob.Spec, (&x.Spec).DeepCopyAs())
			if err != nil {
				return nil, err
			}
			cronJob.Spec = *spec

			return cronJob, nil
		case *v1alpha1.DeployConfigMap:
			cm := &corev1.ConfigMap{}
			cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
			cm.SetNamespace(kpkg.Namespace)
			cm.SetName(kpkg.Name)

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
			cm.Data = data

			if err := AnnotateDigestTo(cm, ScopeConfigMapDigest, kpkg.Name, cm.Data); err != nil {
				return nil, err
			}

			return cm, nil
		case *v1alpha1.DeploySecret:
			secret := &corev1.Secret{}
			secret.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))
			secret.SetNamespace(kpkg.Namespace)
			secret.SetName(kpkg.Name)

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
			secret.StringData = data

			if err := AnnotateDigestTo(secret, ScopeSecretDigest, kpkg.Name, secret.StringData); err != nil {
				return nil, err
			}

			return secret, nil
		}
	}

	return nil, nil
}

func toPodTemplateSpec(kpkg *v1alpha1.KubePkg) (*corev1.PodTemplateSpec, error) {
	if len(kpkg.Spec.Containers) == 0 {
		return nil, errors.New("containers should not empty")
	}
	template := &corev1.PodTemplateSpec{}
	if template.Labels == nil {
		template.Labels = map[string]string{}
	}
	template.Labels["app"] = kpkg.Name

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

		template.Spec.Affinity = must(template.Spec.Affinity)
		template.Spec.Affinity.NodeAffinity = must(template.Spec.Affinity.NodeAffinity)
		template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = must(template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution)

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

func toContainer(c v1alpha1.Container, name string, podTemplateSpec *corev1.PodTemplateSpec, kpkg *v1alpha1.KubePkg) (*corev1.Container, error) {
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

	envs := make(map[string]v1alpha1.EnvVarValueOrFrom)

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
		p.Protocol = PortProtocol(n)
		p.ContainerPort = c.Ports[n]
		p.Name, p.HostPort = PortNameAndHostPort(n)

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
	}
}

func must[T any](v *T) *T {
	if v == nil {
		v = new(T)
	}
	return v
}
