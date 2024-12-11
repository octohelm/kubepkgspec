package reloader

import (
	"iter"
	"maps"
	"slices"

	"github.com/octohelm/kubekit/pkg/metadata"
	"github.com/octohelm/kubepkgspec/pkg/object"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func Patch(objIter iter.Seq[object.Object]) error {
	secretDigests := map[string]string{}
	configMapDigests := map[string]string{}

	for m := range objIter {
		switch x := m.(type) {
		case *corev1.ConfigMap:
			if v, ok := AnnotationSpecDigest.Get(x); ok {
				configMapDigests[x.Name] = v
			}
		case *corev1.Secret:
			if v, ok := AnnotationSpecDigest.Get(x); ok {
				secretDigests[x.Name] = v
			}
		}
	}

	for m := range objIter {
		switch x := m.(type) {
		case *appsv1.Deployment:
			r := &patcher{
				configMapDigests: configMapDigests,
				secretDigests:    secretDigests,
			}
			r.walk(&x.Spec.Template.Spec)
			if err := r.MarshalTo(&x.Spec.Template); err != nil {
				return err
			}
		case *appsv1.DaemonSet:
			r := &patcher{
				configMapDigests: configMapDigests,
				secretDigests:    secretDigests,
			}
			r.walk(&x.Spec.Template.Spec)
			if err := r.MarshalTo(&x.Spec.Template); err != nil {
				return err
			}
		case *appsv1.StatefulSet:
			r := &patcher{
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

type patcher struct {
	configMaps       map[string]bool
	secrets          map[string]bool
	configMapDigests map[string]string
	secretDigests    map[string]string
}

func (r *patcher) recordSecret(name string) {
	if r.secrets == nil {
		r.secrets = map[string]bool{}
	}
	r.secrets[name] = true
}

func (r *patcher) walk(podSpec *corev1.PodSpec) {
	for _, c := range podSpec.Containers {
		for _, envFrom := range c.EnvFrom {
			if envFrom.ConfigMapRef != nil {
				r.recordConfigMap(envFrom.ConfigMapRef.Name)
			}
			if envFrom.SecretRef != nil {
				r.recordSecret(envFrom.SecretRef.Name)
			}
		}

		for _, env := range c.Env {
			if env.ValueFrom != nil {
				if env.ValueFrom.ConfigMapKeyRef != nil {
					r.recordSecret(env.ValueFrom.ConfigMapKeyRef.Name)
				}

				if env.ValueFrom.SecretKeyRef != nil {
					r.recordSecret(env.ValueFrom.SecretKeyRef.Name)
				}
			}

		}
	}

	for _, v := range podSpec.Volumes {
		if cm := v.ConfigMap; cm != nil {
			r.recordConfigMap(cm.Name)
		}
		if s := v.Secret; s != nil {
			r.recordSecret(s.SecretName)
		}
	}
}

func (r *patcher) recordConfigMap(name string) {
	if r.configMaps == nil {
		r.configMaps = map[string]bool{}
	}
	r.configMaps[name] = true
}

func (r *patcher) MarshalTo(o metadata.AnnotationsAccessor) error {
	if len(r.configMaps) > 0 {
		if err := AnnotationConfigMapReload.MarshalTo(o, slices.Collect(maps.Keys(r.configMaps)), ""); err != nil {
			return err
		}

		if r.configMapDigests != nil {
			for name := range r.configMaps {
				if dgst, ok := r.configMapDigests[name]; ok {
					if err := ScopeConfigMapDigest.MustAnnotation(name).MarshalTo(o, dgst, ""); err != nil {
						return err
					}
				}
			}
		}
	}
	if len(r.secrets) > 0 {
		if err := AnnotationSecretReload.MarshalTo(o, slices.Collect(maps.Keys(r.secrets)), ""); err != nil {
			return err
		}

		if r.secretDigests != nil {
			for name := range r.secrets {
				if dgst, ok := r.secretDigests[name]; ok {
					if err := ScopeSecretDigest.MustAnnotation(name).MarshalTo(o, dgst, ""); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
