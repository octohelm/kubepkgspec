package kubepkg

import (
	"cmp"
	"fmt"
	"iter"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	kubepkgv1alpha1 "github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/kubepkg/convert"
	"github.com/octohelm/kubepkgspec/pkg/object"
	"github.com/octohelm/kubepkgspec/pkg/wellknown"
	"github.com/octohelm/kubepkgspec/pkg/workload"
)

func ExtractAsKubePkg(objIter iter.Seq[object.Object]) (*kubepkgv1alpha1.KubePkg, error) {
	ii, err := Extract(objIter)
	if err != nil {
		return nil, err
	}

	kk := &kubepkgv1alpha1.KubePkg{}
	kk.Spec.Manifests = map[string]any{}

	for o := range ii {
		kk.Spec.Manifests[id(o)] = o
	}

	return kk, nil
}

func Extract(objIter iter.Seq[object.Object]) (iter.Seq[object.Object], error) {
	manifests := map[string]object.Object{}
	workloads := make([]object.Object, 0)

	for o := range objIter {
		if workload.IsWorkload(o) {
			workloads = append(workloads, o)
			continue
		}
		manifests[id(o)] = o
	}

	if len(workloads) > 0 {
		s := &extractor{
			manifests: manifests,
			used:      map[string]bool{},
		}

		kubepkgs := make([]*kubepkgv1alpha1.KubePkg, 0)

		for _, w := range workloads {
			kpkg, err := s.SimpleWorkload(w)
			if err != nil {
				return nil, err
			}
			kubepkgs = append(kubepkgs, kpkg)
		}

		return func(yield func(object.Object) bool) {
			for _, x := range kubepkgs {
				if !yield(x) {
					return
				}
			}

			for o := range s.RemainManifests() {
				if !yield(o) {
					return
				}
			}
		}, nil
	}

	return func(yield func(object.Object) bool) {
	}, nil
}

type extractor struct {
	manifests map[string]object.Object
	used      map[string]bool
}

func (e *extractor) RemainManifests() iter.Seq[object.Object] {
	return func(yield func(object.Object) bool) {
		for _, o := range e.manifests {
			switch o.(type) {
			case *rbacv1.Role, *rbacv1.RoleBinding, *rbacv1.ClusterRole, *rbacv1.ClusterRoleBinding, *corev1.ServiceAccount:
				continue
			case *corev1.ConfigMap, *corev1.Secret:
				continue
			case *corev1.Namespace:
				continue
			}

			if !e.used[id(o)] {
				if !yield(o) {
					return
				}
			}
		}
	}
}

func (e extractor) SimpleWorkload(o object.Object) (*kubepkgv1alpha1.KubePkg, error) {
	k := &kubepkgv1alpha1.KubePkg{}
	k.SetGroupVersionKind(kubepkgv1alpha1.SchemeGroupVersion.WithKind("KubePkg"))

	k.Name = o.GetName()

	if version := wellknown.LabelAppVersion.GetFrom(o); version != "" {
		k.Spec.Version = version
	}

	switch x := o.(type) {
	case *appsv1.Deployment:
		d := &kubepkgv1alpha1.DeployDeployment{}
		d.Kind = "Deployment"

		if err := e.resolveServiceAccount(k, x.Spec.Template.Spec.ServiceAccountName); err != nil {
			return nil, err
		}
		if err := e.resolveContainersAndVolumes(k, &x.Spec.Template); err != nil {
			return nil, err
		}
		if err := e.resolveNetworks(k, x.Spec.Selector); err != nil {
			return nil, err
		}

		if err := convert.Unmarshal(x.Spec, &d.Spec); err != nil {
			return nil, err
		}
		k.Spec.Deploy.SetUnderlying(d)
	case *appsv1.DaemonSet:
		d := &kubepkgv1alpha1.DeployDaemonSet{}
		d.Kind = "DaemonSet"

		if err := e.resolveServiceAccount(k, x.Spec.Template.Spec.ServiceAccountName); err != nil {
			return nil, err
		}
		if err := e.resolveContainersAndVolumes(k, &x.Spec.Template); err != nil {
			return nil, err
		}
		if err := e.resolveNetworks(k, x.Spec.Selector); err != nil {
			return nil, err
		}

		if err := convert.Unmarshal(x.Spec, &d.Spec); err != nil {
			return nil, err
		}
		k.Spec.Deploy.SetUnderlying(d)
	case *appsv1.StatefulSet:
		d := &kubepkgv1alpha1.DeployStatefulSet{}
		d.Kind = "StatefulSet"

		if err := e.resolveServiceAccount(k, x.Spec.Template.Spec.ServiceAccountName); err != nil {
			return nil, err
		}
		if err := e.resolveContainersAndVolumes(k, &x.Spec.Template); err != nil {
			return nil, err
		}
		if err := e.resolveNetworks(k, x.Spec.Selector); err != nil {
			return nil, err
		}

		if err := convert.Unmarshal(x.Spec, &d.Spec); err != nil {
			return nil, err
		}
		k.Spec.Deploy.SetUnderlying(d)
	case *batchv1.Job:
		d := &kubepkgv1alpha1.DeployJob{}
		d.Kind = "Job"

		if err := e.resolveServiceAccount(k, x.Spec.Template.Spec.ServiceAccountName); err != nil {
			return nil, err
		}
		if err := e.resolveContainersAndVolumes(k, &x.Spec.Template); err != nil {
			return nil, err
		}

		if err := convert.Unmarshal(x.Spec, &d.Spec); err != nil {
			return nil, err
		}
		k.Spec.Deploy.SetUnderlying(d)
	case *batchv1.CronJob:
		d := &kubepkgv1alpha1.DeployCronJob{}
		d.Kind = "CronJob"

		if err := e.resolveServiceAccount(k, x.Spec.JobTemplate.Spec.Template.Spec.ServiceAccountName); err != nil {
			return nil, err
		}
		if err := e.resolveContainersAndVolumes(k, &x.Spec.JobTemplate.Spec.Template); err != nil {
			return nil, err
		}

		if err := convert.Unmarshal(x.Spec, &d.Spec); err != nil {
			return nil, err
		}
		k.Spec.Deploy.SetUnderlying(d)
	}

	if err := e.hositAsConfig(k); err != nil {
		return nil, err
	}

	return k, nil
}

func (e *extractor) hositAsConfig(k *kubepkgv1alpha1.KubePkg) error {
	addConfig := func(key string, value kubepkgv1alpha1.EnvVarValueOrFrom) {
		if k.Spec.Config == nil {
			k.Spec.Config = map[string]kubepkgv1alpha1.EnvVarValueOrFrom{}
		}
		k.Spec.Config[key] = value
	}

	if len(k.Spec.Containers) == 1 {
		for name, c := range k.Spec.Containers {
			for key, v := range c.Env {
				addConfig(key, v)
			}
			c.Env = nil
			k.Spec.Containers[name] = c
		}
	}

	for n, v := range k.Spec.Volumes {
		if vcm, ok := v.Underlying.(*kubepkgv1alpha1.VolumeConfigMap); ok {
			if vcm.MountPath == "export" && vcm.Spec != nil {
				for k, v := range vcm.Spec.Data {
					addConfig(k, kubepkgv1alpha1.EnvVarValueOrFrom{
						Value: v,
					})
				}
				delete(k.Spec.Volumes, n)
			}
		}
	}

	return nil
}

func id(o object.Object) string {
	return fmt.Sprintf("%s.%s", o.GetName(), o.GetObjectKind().GroupVersionKind().Kind)
}

func (e *extractor) markUse(o object.Object) {
	e.used[id(o)] = true
}

func (e *extractor) resolveContainersAndVolumes(k *kubepkgv1alpha1.KubePkg, podSpec *corev1.PodTemplateSpec) error {
	err := convert.ContainersAndVolumesFromPodTemplateSpec(k, podSpec)
	if err != nil {
		return err
	}

	if platforms := resolvePlatformsFromSpec(podSpec); len(platforms) > 0 {
		for n, c := range k.Spec.Containers {
			c.Image.Platforms = platforms
			k.Spec.Containers[n] = c
		}
	}

	if k.Spec.Volumes != nil {
		for _, v := range k.Spec.Volumes {
			switch x := v.Underlying.(type) {
			case *kubepkgv1alpha1.VolumeConfigMap:
				data := e.resolveConfigMapData(x.Opt.Name)
				x.Spec = &kubepkgv1alpha1.ConfigMapSpec{
					Data: data,
				}
				x.Opt.Name = ""
			case *kubepkgv1alpha1.VolumeSecret:
				data := e.resolveSecret(x.Opt.SecretName)
				x.Spec = &kubepkgv1alpha1.ConfigMapSpec{
					Data: data,
				}
				x.Opt.SecretName = ""
			}
		}
	}

	return nil
}

func resolvePlatformsFromSpec(t *corev1.PodTemplateSpec) (platforms []string) {
	if affinity := t.Spec.Affinity; affinity != nil {
		if nodeAffinity := t.Spec.Affinity.NodeAffinity; nodeAffinity != nil {
			if requiredDuringSchedulingIgnoredDuringExecution := nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution; requiredDuringSchedulingIgnoredDuringExecution != nil {
				nodeSelectorTerms := make([]corev1.NodeSelectorTerm, 0, len(requiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms))

				for _, term := range requiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms {
					if len(term.MatchExpressions) == 1 {
						if term.MatchExpressions[0].Operator == corev1.NodeSelectorOpIn && term.MatchExpressions[0].Key == "kubernetes.io/arch" {

							archs := term.MatchExpressions[0].Values

							platforms = make([]string, len(archs))

							for i, arch := range archs {
								platforms[i] = fmt.Sprintf("linux/%s", arch)
							}

							continue
						}
					}
					nodeSelectorTerms = append(nodeSelectorTerms, term)
				}

				if len(nodeSelectorTerms) > 0 {
					requiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = nodeSelectorTerms
				} else {
					nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = nil
				}
			}
		}
	}

	return
}

func (e *extractor) resolveNetworks(k *kubepkgv1alpha1.KubePkg, selector *metav1.LabelSelector) error {
	if selector == nil || len(selector.MatchLabels) == 0 {
		return nil
	}

	for _, o := range e.manifests {
		if s, ok := o.(*corev1.Service); ok {
			if mapContains(selector.MatchLabels, s.Spec.Selector) {
				e.markUse(s)

				if k.Spec.Services == nil {
					k.Spec.Services = map[string]kubepkgv1alpha1.Service{}
				}

				svc := kubepkgv1alpha1.Service{
					Ports: map[string]int32{},
				}

				serviceName := k.Name

				for _, p := range s.Spec.Ports {
					name := convert.FormatPortName(cmp.Or(p.TargetPort.StrVal, p.Name), p.Protocol, 0)

					svc.Ports[name] = p.Port
				}

				if s.Name == k.Name {
					serviceName = "#"
				} else if strings.HasPrefix(s.Name, k.Name) {
					serviceName = strings.Trim(serviceName[len(k.Name)+1:], "-")
				}

				if svc.Expose == nil {
					if err := e.resolvePathsFromIngressOrGateway(&svc, s); err != nil {
						return err
					}
				}

				k.Spec.Services[serviceName] = svc

				return nil
			}
		}
	}

	return nil
}

func (e *extractor) resolvePathsFromIngressOrGateway(ks *kubepkgv1alpha1.Service, s *corev1.Service) error {
	for _, o := range e.manifests {
		switch x := o.(type) {
		case *gatewayv1.HTTPRoute:
			for _, r := range x.Spec.Rules {
				for _, m := range r.Matches {
					if p := m.Path; p != nil {
						if path := p.Value; path != nil {
							if len(r.BackendRefs) > 0 {
								b := r.BackendRefs[0]

								portName := ""
								for _, p := range s.Spec.Ports {
									if b.Port != nil && int32(*b.Port) == p.Port {
										portName = p.Name
										break
									}
								}

								if portName != "" {
									e.markUse(o)

									if ks.Paths == nil {
										ks.Paths = map[string]kubepkgv1alpha1.StringOrSlice{}
									}

									ks.Paths[portName] = append(ks.Paths[portName], *path)
								}
							}
						}
					}
				}
			}
		case *networkingv1.Ingress:
			for _, r := range x.Spec.Rules {
				if r.HTTP != nil {
					for _, p := range r.HTTP.Paths {
						if backendSvc := p.Backend.Service; backendSvc != nil {
							if backendSvc.Name == s.Name {

								portName := backendSvc.Port.Name
								if portName == "" {
									for _, p := range s.Spec.Ports {
										if backendSvc.Port.Number == p.Port {
											portName = p.Name
											break
										}
									}
								}

								if portName != "" {
									e.markUse(o)

									if ks.Paths == nil {
										ks.Paths = map[string]kubepkgv1alpha1.StringOrSlice{}
									}

									ks.Paths[portName] = append(ks.Paths[portName], p.Path)
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func mapContains[K comparable, V comparable](scope, target map[K]V) bool {
	for k, v := range target {
		if w, ok := scope[k]; !ok || v != w {
			return false
		}
	}
	return true
}

func (e *extractor) resolveConfigMapData(configMapName string) map[string]string {
	for _, o := range e.manifests {
		if x, ok := o.(*corev1.ConfigMap); ok {
			if x.Name == configMapName {
				e.markUse(x)
				return x.Data
			}
		}
	}
	return nil
}

func (e *extractor) resolveSecret(secretName string) map[string]string {
	for _, o := range e.manifests {
		if x, ok := o.(*corev1.Secret); ok {
			if x.Name == secretName {
				e.markUse(x)
				return x.StringData
			}
		}
	}
	return nil
}

func (e *extractor) resolveServiceAccount(k *kubepkgv1alpha1.KubePkg, serviceAccountName string) error {
	if serviceAccountName == "" {
		return nil
	}

	for _, o := range e.manifests {
		if x, ok := o.(*corev1.ServiceAccount); ok {
			if x.Name == serviceAccountName {
				e.markUse(x)
				k.Spec.ServiceAccount = e.resolveRoleBinding(x.Name)
				return nil
			}
		}
	}

	return nil
}

func (e *extractor) resolveRoleBinding(serviceAccountName string) *kubepkgv1alpha1.ServiceAccount {
	sa := &kubepkgv1alpha1.ServiceAccount{
		Scope: kubepkgv1alpha1.ScopeTypeNamespace,
	}

	for _, o := range e.manifests {
		switch x := o.(type) {
		case *rbacv1.RoleBinding:
			for _, sub := range x.Subjects {
				if sub.Kind == "ServiceAccount" {
					if sub.Name == serviceAccountName {
						e.markUse(x)

						sa.Rules = append(sa.Rules, e.resolveRules(x.RoleRef)...)
					}
				}
			}
		case *rbacv1.ClusterRoleBinding:
			for _, sub := range x.Subjects {
				if sub.Kind == "ServiceAccount" {
					if sub.Name == serviceAccountName {
						e.markUse(x)

						sa.Rules = append(sa.Rules, e.resolveRules(x.RoleRef)...)
						sa.Scope = kubepkgv1alpha1.ScopeTypeCluster
					}
				}
			}
		}
	}

	if len(sa.Rules) > 0 {
		return sa
	}

	return nil
}

func (e *extractor) resolveRules(ref rbacv1.RoleRef) []rbacv1.PolicyRule {
	for _, o := range e.manifests {
		switch x := o.(type) {
		case *rbacv1.Role:
			if ref.Kind == "Role" && x.Name == ref.Name {
				e.markUse(x)
				return x.Rules
			}
		case *rbacv1.ClusterRole:
			if ref.Kind == "ClusterRole" && x.Name == ref.Name {
				e.markUse(x)
				return x.Rules
			}
		}
	}
	return nil
}
