package kubepkg

import (
	"fmt"
	"iter"
	"maps"
	"sort"
	"strings"

	"github.com/containerd/platforms"
	"github.com/go-json-experiment/json"
	"github.com/octohelm/x/ptr"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/kubepkg/convert"
	"github.com/octohelm/kubepkgspec/pkg/object"
	"github.com/octohelm/kubepkgspec/pkg/reloader"
)

type Option = func(c *converter)

func WithRecursive(recursive bool) Option {
	return func(c *converter) {
		c.recursive = recursive
	}
}

func Convert(kpkg *v1alpha1.KubePkg, options ...Option) (iter.Seq[object.Object], error) {
	e := &converter{}
	e.build(options...)

	if err := e.walk(kpkg); err != nil {
		return nil, err
	}
	return func(yield func(object.Object) bool) {
		for _, m := range e.manifests {
			if !yield(m) {
				return
			}
		}
	}, nil
}

type converter struct {
	manifests map[string]object.Object
	recursive bool
}

func (c *converter) build(options ...Option) {
	for _, o := range options {
		o(c)
	}
}

func (e *converter) register(o object.Object, kpkg *v1alpha1.KubePkg) error {
	if o == nil || o.GetObjectKind() == nil {
		return nil
	}
	if e.manifests == nil {
		e.manifests = map[string]object.Object{}
	}

	if err := convert.LabelInstanceAndVersion(kpkg, o); err != nil {
		return err
	}
	e.manifests[objectID(o)] = o
	return nil
}

func (e *converter) walk(kpkg *v1alpha1.KubePkg) error {
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

func (e *converter) walkNetworks(kpkg *v1alpha1.KubePkg) error {
	for n := range kpkg.Spec.Services {
		s := kpkg.Spec.Services[n]

		service := &corev1.Service{}
		service.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Service"))
		service.SetNamespace(kpkg.Namespace)
		service.SetName(convert.SubResourceName(kpkg, n))

		if kpkg.Spec.Deploy.Underlying.GetKind() == (&v1alpha1.DeployEndpointSlice{}).GetKind() {
			service.Spec.ClusterIP = corev1.ClusterIPNone
		} else {
			service.Spec.Selector = map[string]string{
				"app": kpkg.Name,
			}

			service.Spec.ClusterIP = s.ClusterIP
		}

		paths := map[string][]string{}

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
			p.Protocol = convert.PortProtocol(n)
			p.Port = s.Ports[n]
			p.Name = n
			p.TargetPort = intstr.FromString(n)

			service.Spec.Ports = append(service.Spec.Ports, p)

			if strings.HasPrefix(p.Name, "http") {
				if _, ok := paths[p.Name]; !ok {
					paths[p.Name] = []string{"/"}
				}
			}
		}

		endpoints := map[string]string{}

		if s.Expose == nil || s.Expose.GetType() != "NodePort" {
			endpoints["default"] = fmt.Sprintf("http://%s", service.Name)
		}

		if s.Expose != nil {
			switch x := s.Expose.Underlying.(type) {
			case *v1alpha1.ExposeNodePort:
				service.Spec.Type = corev1.ServiceTypeNodePort
				for i, p := range service.Spec.Ports {
					service.Spec.Ports[i].NodePort = p.Port
				}
			case *v1alpha1.ExposeIngress:
				if len(x.Gateway) > 0 {
					for _, gateway := range x.Gateway {
						if len(s.Paths) > 0 {
							httpRoute := &gatewayv1.HTTPRoute{}
							httpRoute.SetGroupVersionKind(gatewayv1.SchemeGroupVersion.WithKind("HTTPRoute"))
							httpRoute.SetNamespace(kpkg.Namespace)
							httpRoute.SetName(convert.SubResourceName(kpkg, n))
							httpRoute.Annotations = map[string]string{}

							parts := strings.Split(gateway, ".")

							if len(parts) > 1 {
								httpRoute.Spec.ParentRefs = []gatewayv1.ParentReference{
									{
										Name:      gatewayv1.ObjectName(parts[0]),
										Namespace: ptr.Ptr(gatewayv1.Namespace(parts[1])),
									},
								}
							} else {
								httpRoute.Spec.ParentRefs = []gatewayv1.ParentReference{
									{
										Name: gatewayv1.ObjectName(parts[0]),
									},
								}
							}

							httpRoute.Spec.Rules = []gatewayv1.HTTPRouteRule{}

							if len(x.Options) > 0 {
								raw, _ := json.Marshal(x.Options)
								if len(raw) > 0 {
									httpRoute.Annotations["path-rule.kubepkg.octohelm.tech/options"] = string(raw)
								}
							}

							for n, matches := range s.Paths {
								if len(matches) == 0 {
									continue
								}

								rule := gatewayv1.HTTPRouteRule{}

								for _, path := range matches {
									rule.Matches = append(rule.Matches, gatewayv1.HTTPRouteMatch{
										Path: &gatewayv1.HTTPPathMatch{
											Type:  ptr.Ptr(gatewayv1.PathMatchPathPrefix),
											Value: ptr.Ptr(path),
										},
									})
								}

								httpBackendRef := gatewayv1.HTTPBackendRef{}
								httpBackendRef.Name = gatewayv1.ObjectName(service.Name)

								if port, ok := s.Ports[n]; ok {
									httpBackendRef.Port = ptr.Ptr(gatewayv1.PortNumber(port))
								} else {
									httpBackendRef.Port = ptr.Ptr(gatewayv1.PortNumber(80))
								}

								rule.BackendRefs = []gatewayv1.HTTPBackendRef{
									httpBackendRef,
								}

								httpRoute.Spec.Rules = append(httpRoute.Spec.Rules, rule)
							}

							if err := e.register(httpRoute, kpkg); err != nil {
								return err
							}
						}
					}
				}
			}
		}

		if err := e.register(service, kpkg); err != nil {
			return err
		}

		if x, ok := kpkg.Spec.Deploy.Underlying.(*v1alpha1.DeployEndpointSlice); ok {
			endpointSlice := &discoveryv1.EndpointSlice{}
			endpointSlice.SetGroupVersionKind(discoveryv1.SchemeGroupVersion.WithKind("EndpointSlice"))
			endpointSlice.SetNamespace(kpkg.Namespace)
			endpointSlice.SetName(service.Name)

			endpointSlice.Labels = maps.Clone(x.Labels)
			endpointSlice.Annotations = maps.Clone(x.Annotations)

			if endpointSlice.Labels == nil {
				endpointSlice.Labels = make(map[string]string)
			}

			endpointSlice.Labels["kubernetes.io/service-name"] = service.Name
			endpointSlice.Labels["endpointslice-controller.k8s.io"] = "manual"

			endpointPortNames := make([]string, 0, len(x.Ports))
			for n := range x.Ports {
				endpointPortNames = append(endpointPortNames, n)
			}
			sort.Strings(endpointPortNames)
			for _, n := range endpointPortNames {
				p := discoveryv1.EndpointPort{}
				name, protocol, _ := convert.DecodePortName(n)

				p.Name = ptr.Ptr(name)
				p.Protocol = ptr.Ptr(protocol)
				p.Port = ptr.Ptr(x.Ports[n])

				endpointSlice.Ports = append(endpointSlice.Ports, p)
			}

			endpoint := discoveryv1.Endpoint{}

			endpointSlice.AddressType = x.AddressType
			endpoint.Addresses = x.Addresses
			endpoint.Conditions.Ready = ptr.Ptr(true)

			endpointSlice.Endpoints = []discoveryv1.Endpoint{endpoint}

			if err := e.register(endpointSlice, kpkg); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *converter) walkVolumes(kpkg *v1alpha1.KubePkg) error {
	vcs := convert.VolumeConvertorsFrom(kpkg)

	for _, c := range vcs {
		r, err := c.ToResource(kpkg)
		if err != nil {
			return err
		}

		if err := e.register(r, kpkg); err != nil {
			return err
		}
	}

	return nil
}

func (e *converter) walkRbac(kpkg *v1alpha1.KubePkg) error {
	if sa := kpkg.Spec.ServiceAccount; sa != nil {
		serviceAccount := &corev1.ServiceAccount{}
		serviceAccount.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ServiceAccount"))
		serviceAccount.SetNamespace(kpkg.Namespace)
		serviceAccount.SetName(kpkg.Name)

		if err := e.register(serviceAccount, kpkg); err != nil {
			return err
		}

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

			if err := e.register(clusterRole, kpkg); err != nil {
				return err
			}

			if err := e.register(clusterRoleBinding, kpkg); err != nil {
				return err
			}

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

		if err := e.register(role, kpkg); err != nil {
			return err
		}

		if err := e.register(roleBinding, kpkg); err != nil {
			return err
		}
	}

	return nil
}

func (e *converter) walkDeploy(kpkg *v1alpha1.KubePkg) error {
	d, err := convert.DeployResourceFrom(kpkg)
	if err != nil {
		return err
	}
	return e.register(d, kpkg)
}

func (e *converter) walkManifests(kpkg *v1alpha1.KubePkg) error {
	c := &completer{
		kpkg: kpkg,
	}

	i := object.NewIter(
		c.patchNamespace,
		c.patchConfigMapOrSecret,
		c.patchNodeAffinityIfNeed,
	)

	for m := range i.Object(kpkg.Spec.Manifests) {
		if e.recursive {
			if k, ok := m.(*v1alpha1.KubePkg); ok {
				if err := e.walk(k); err != nil {
					return err
				}
				continue
			}
		}

		if err := e.register(m, kpkg); err != nil {
			return err
		}
	}

	return i.Err()
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
		if err := reloader.AnnotateDigestTo(x, x.Data); err != nil {
			return nil, err
		}
		return x, nil
	case *corev1.Secret:
		if err := reloader.AnnotateDigestTo(x, x.StringData); err != nil {
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
		pod.Spec.Affinity = convert.Must(pod.Spec.Affinity)
		pod.Spec.Affinity.NodeAffinity = convert.Must(pod.Spec.Affinity.NodeAffinity)
		pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = convert.Must(pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution)

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

func objectID(d object.Object) string {
	return fmt.Sprintf("%s.%s", strings.ToLower(d.GetObjectKind().GroupVersionKind().Kind), d.GetName())
}
