package manifest

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/octohelm/kubepkgspec/pkg/manifest/object"
	batchv1 "k8s.io/api/batch/v1"

	"github.com/containerd/containerd/platforms"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/manifest/strfmt"
)

func SortedExtract(kpkg *v1alpha1.KubePkg) ([]object.Object, error) {
	manifests, err := Extract(kpkg)
	if err != nil {
		return nil, err
	}

	list := make([]object.Object, 0, len(manifests)+1)

	if namespace := kpkg.Namespace; namespace != "" {
		n := &corev1.Namespace{}
		n.APIVersion = "v1"
		n.Kind = "Namespace"
		n.Name = namespace
		list = append(list, n)
	}

	for k := range manifests {
		list = append(list, manifests[k])
	}

	return SortByKind(list), nil
}

func Extract(kpkg *v1alpha1.KubePkg) (map[string]object.Object, error) {
	e := &extractor{}
	if err := e.walk(kpkg); err != nil {
		return nil, err
	}
	if err := e.patchConfigMapOrSecretDeps(); err != nil {
		return nil, err
	}
	return e.manifests, nil
}

type extractor struct {
	manifests map[string]object.Object
}

func (e *extractor) register(o object.Object) {
	if o == nil || o.GetObjectKind() == nil {
		return
	}
	if e.manifests == nil {
		e.manifests = map[string]object.Object{}
	}
	e.manifests[objectID(o)] = o
}

func objectID(d object.Object) string {
	return fmt.Sprintf("%s.%s", strings.ToLower(d.GetObjectKind().GroupVersionKind().Kind), d.GetName())
}

func (e *extractor) walk(kpkg *v1alpha1.KubePkg) error {
	if err := e.walkDeploy(kpkg); err != nil {
		return err
	}
	if err := e.walkVolumes(kpkg); err != nil {
		return err
	}
	if err := e.walkNetworks(kpkg); err != nil {
		return err
	}
	if err := e.walkRbac(kpkg); err != nil {
		return err
	}
	if err := e.walkManifests(kpkg); err != nil {
		return err
	}
	return nil
}

func (e *extractor) walkNetworks(kpkg *v1alpha1.KubePkg) error {
	var gatewayTemplates []strfmt.GatewayTemplate
	if err := AnnotationIngressGateway.UnmarshalFrom(kpkg, &gatewayTemplates); err != nil {
		return err
	}

	for n := range kpkg.Spec.Services {
		s := kpkg.Spec.Services[n]

		service := &corev1.Service{}
		service.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Service"))
		service.SetNamespace(kpkg.Namespace)
		service.SetName(SubResourceName(kpkg, n))

		service.Spec.Selector = map[string]string{
			"app": kpkg.Name,
		}

		service.Spec.ClusterIP = s.ClusterIP

		paths := map[string]string{}
		for portName, p := range s.Paths {
			paths[portName] = p
		}

		portNames := make([]string, 0, len(s.Ports))
		for n := range s.Ports {
			portNames = append(portNames, n)
		}
		sort.Strings(portNames)

		for _, n := range portNames {
			p := corev1.ServicePort{}
			p.Protocol = PortProtocol(n)
			p.Port = s.Ports[n]
			p.Name = n
			p.TargetPort = intstr.FromString(n)
			service.Spec.Ports = append(service.Spec.Ports, p)

			if strings.HasPrefix(p.Name, "http") {
				if _, ok := paths[p.Name]; !ok {
					paths[p.Name] = "/"
				}
			}
		}

		endpoints := map[string]string{}

		if n == "#" && len(gatewayTemplates) > 0 {
			if len(paths) > 0 && s.Expose == nil {
				s.Expose = &v1alpha1.Expose{
					Type: "Ingress",
				}
			}
		}

		if s.Expose == nil || s.Expose.Type != "NodePort" {
			endpoints["default"] = fmt.Sprintf("http://%s", service.Name)
		}

		if s.Expose != nil {
			switch s.Expose.Type {
			case "NodePort":
				service.Spec.Type = corev1.ServiceTypeNodePort
				for i, p := range service.Spec.Ports {
					service.Spec.Ports[i].NodePort = p.Port
				}
			case "Ingress":
				if len(gatewayTemplates) > 0 {
					igs := strfmt.From(gatewayTemplates)

					rules := igs.For(service.Name, service.Namespace).IngressRules(paths, s.Expose.Gateway...)

					for name, e := range igs.For(service.Name, service.Namespace).Endpoints() {
						endpoints[name] = e
					}

					if len(rules) > 0 {
						ingress := &networkingv1.Ingress{}
						ingress.SetGroupVersionKind(networkingv1.SchemeGroupVersion.WithKind("Ingress"))
						ingress.SetNamespace(kpkg.Namespace)
						ingress.SetName(SubResourceName(kpkg, n))

						ingress.Spec.Rules = rules

						e.register(ingress)
					}
				}
			}
		}

		e.register(service)

		if len(endpoints) > 0 {
			cmForEndpoints := &corev1.ConfigMap{}
			cmForEndpoints.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
			cmForEndpoints.SetNamespace(service.Namespace)
			cmForEndpoints.SetName(fmt.Sprintf("endpoint-%s", service.Name))
			cmForEndpoints.Data = endpoints

			e.register(cmForEndpoints)
		}
	}

	return nil
}

func (e *extractor) walkVolumes(kpkg *v1alpha1.KubePkg) error {
	vcs := VolumeConvertorsFrom(kpkg)

	for _, c := range vcs {
		r, err := c.ToResource(kpkg)
		if err != nil {
			return err
		}
		e.register(r)
	}

	return nil
}

func (e *extractor) walkRbac(kpkg *v1alpha1.KubePkg) error {
	if sa := kpkg.Spec.ServiceAccount; sa != nil {
		serviceAccount := &corev1.ServiceAccount{}
		serviceAccount.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ServiceAccount"))
		serviceAccount.SetNamespace(kpkg.Namespace)
		serviceAccount.SetName(kpkg.Name)

		e.register(serviceAccount)

		if sa.Scope == v1alpha1.ScopeTypeCluster {
			clusterRole := &rbacv1.ClusterRole{}
			clusterRole.SetGroupVersionKind(rbacv1.SchemeGroupVersion.WithKind("ClusterRole"))
			clusterRole.SetNamespace(kpkg.Namespace)
			clusterRole.SetName(kpkg.Name)
			clusterRole.Rules = sa.Rules

			clusterRoleBinding := &rbacv1.ClusterRoleBinding{}
			clusterRoleBinding.SetGroupVersionKind(rbacv1.SchemeGroupVersion.WithKind("ClusterRoleBinding"))
			clusterRoleBinding.SetNamespace(kpkg.Namespace)
			clusterRoleBinding.SetName(kpkg.Name)

			clusterRoleBinding.RoleRef.Name = clusterRole.Name
			clusterRoleBinding.RoleRef.Kind = clusterRole.Kind
			clusterRoleBinding.RoleRef.APIGroup = rbacv1.SchemeGroupVersion.Group

			clusterRoleBinding.Subjects = []rbacv1.Subject{{
				Kind:      serviceAccount.Kind,
				Name:      serviceAccount.Name,
				Namespace: serviceAccount.Namespace,
			}}

			e.register(clusterRole)
			e.register(clusterRoleBinding)

			return nil
		}
		role := &rbacv1.Role{}
		role.SetGroupVersionKind(rbacv1.SchemeGroupVersion.WithKind("Role"))
		role.SetNamespace(kpkg.Namespace)
		role.SetName(kpkg.Name)
		role.Rules = sa.Rules

		roleBinding := &rbacv1.RoleBinding{}
		roleBinding.SetGroupVersionKind(rbacv1.SchemeGroupVersion.WithKind("RoleBinding"))
		roleBinding.SetNamespace(kpkg.Namespace)
		roleBinding.SetName(kpkg.Name)

		roleBinding.RoleRef.Name = role.Name
		roleBinding.RoleRef.Kind = role.Kind
		roleBinding.RoleRef.APIGroup = rbacv1.SchemeGroupVersion.Group

		roleBinding.Subjects = []rbacv1.Subject{{
			Kind:      serviceAccount.Kind,
			Name:      serviceAccount.Name,
			Namespace: serviceAccount.Namespace,
		}}

		e.register(role)
		e.register(roleBinding)
	}

	return nil
}

func (e *extractor) walkDeploy(kpkg *v1alpha1.KubePkg) error {
	d, err := DeployResourceFrom(kpkg)
	if err != nil {
		return err
	}
	e.register(d)
	return nil
}

func (e *extractor) walkManifests(kpkg *v1alpha1.KubePkg) error {
	c := &completer{
		kpkg: kpkg,
	}

	i := object.NewIter(
		c.patchNamespace,
		c.patchConfigMapOrSecret,
		c.patchNodeAffinityIfNeed,
	)

	for m := range i.Object(context.Background(), kpkg.Spec.Manifests) {
		e.register(m)
	}

	return i.Err()
}

func (e *extractor) patchConfigMapOrSecretDeps() error {
	secretDigests := map[string]string{}
	configMapDigests := map[string]string{}

	for _, m := range e.manifests {
		switch x := m.(type) {
		case *corev1.ConfigMap:
			a := ScopeConfigMapDigest.MustAnnotation(x.Name)
			if v, ok := a.Get(x); ok {
				configMapDigests[x.Name] = v
			}
		case *corev1.Secret:
			a := ScopeSecretDigest.MustAnnotation(x.Name)
			if v, ok := a.Get(x); ok {
				secretDigests[x.Name] = v
			}
		}
	}

	for _, m := range e.manifests {
		switch x := m.(type) {
		case *appsv1.Deployment:
			r := &reload{
				configMapDigests: configMapDigests,
				secretDigests:    secretDigests,
			}
			r.walk(&x.Spec.Template.Spec)
			if err := r.MarshalTo(&x.Spec.Template); err != nil {
				return err
			}
		case *appsv1.DaemonSet:
			r := &reload{
				configMapDigests: configMapDigests,
				secretDigests:    secretDigests,
			}
			r.walk(&x.Spec.Template.Spec)
			if err := r.MarshalTo(&x.Spec.Template); err != nil {
				return err
			}
		case *appsv1.StatefulSet:
			r := &reload{
				configMapDigests: configMapDigests,
				secretDigests:    secretDigests,
			}
			r.walk(&x.Spec.Template.Spec)
			if err := r.MarshalTo(&x.Spec.Template); err != nil {
				return err
			}
		}
	}

	return nil
}

type completer struct {
	kpkg *v1alpha1.KubePkg
}

func (c *completer) patchNamespace(o object.Object) (object.Object, error) {
	o.SetNamespace(c.kpkg.Namespace)

	switch x := o.(type) {
	case *rbacv1.RoleBinding:
		for i := range x.Subjects {
			s := x.Subjects[i]
			s.Namespace = c.kpkg.Namespace
			x.Subjects[i] = s
		}
		return x, nil
	case *rbacv1.ClusterRoleBinding:
		for i := range x.Subjects {
			s := x.Subjects[i]
			s.Namespace = c.kpkg.Namespace
			x.Subjects[i] = s
		}
		return x, nil
	}

	return o, nil
}

func (c *completer) patchConfigMapOrSecret(o object.Object) (object.Object, error) {
	switch x := o.(type) {
	case *corev1.ConfigMap:
		if err := AnnotateDigestTo(x, ScopeConfigMapDigest, x.Name, x.Data); err != nil {
			return nil, err
		}
		return x, nil
	case *corev1.Secret:
		if err := AnnotateDigestTo(x, ScopeSecretDigest, x.Name, x.StringData); err != nil {
			return nil, err
		}
		return x, nil
	}
	return o, nil
}

func (c *completer) patchNodeAffinityIfNeed(o object.Object) (object.Object, error) {
	var platformList []string

	if err := AnnotationPlatform.UnmarshalFrom(o, &platformList); err != nil {
		return nil, err
	}

	if len(platformList) > 0 {
		switch d := o.(type) {
		case *appsv1.Deployment:
			c.patchPodNodeAffinity(platformList, &d.Spec.Template)
		case *appsv1.DaemonSet:
			c.patchPodNodeAffinity(platformList, &d.Spec.Template)
		case *appsv1.StatefulSet:
			c.patchPodNodeAffinity(platformList, &d.Spec.Template)
		case *batchv1.Job:
			c.patchPodNodeAffinity(platformList, &d.Spec.Template)
		case *batchv1.CronJob:
			c.patchPodNodeAffinity(platformList, &d.Spec.JobTemplate.Spec.Template)
		}
	}

	return o, nil
}

func (c *completer) patchPodNodeAffinity(pls []string, pod *corev1.PodTemplateSpec) {
	archs := make([]string, 0, len(pls))

	for _, p := range pls {
		pl, err := platforms.Parse(p)
		if err == nil {
			archs = append(archs, pl.Architecture)
		}
	}

	if len(archs) > 0 {
		pod.Spec.Affinity = must(pod.Spec.Affinity)
		pod.Spec.Affinity.NodeAffinity = must(pod.Spec.Affinity.NodeAffinity)
		pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = must(pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution)

		if len(pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms) == 0 {
			pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = append(
				pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms,
				corev1.NodeSelectorTerm{
					MatchExpressions: []corev1.NodeSelectorRequirement{{
						Key:      "kubernetes.io/arch",
						Operator: "In",
						Values:   archs,
					}},
				},
			)
		}
	}
}
