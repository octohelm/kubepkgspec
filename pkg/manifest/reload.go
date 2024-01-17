package manifest

import (
	"github.com/octohelm/kubepkgspec/pkg/kubeutil"
	corev1 "k8s.io/api/core/v1"
	"sort"
)

type reload struct {
	configMaps       map[string]bool
	secrets          map[string]bool
	configMapDigests map[string]string
	secretDigests    map[string]string
}

func (r *reload) recordSecret(name string) {
	if r.secrets == nil {
		r.secrets = map[string]bool{}
	}
	r.secrets[name] = true
}

func (r *reload) walk(podSpec *corev1.PodSpec) {
	for _, c := range podSpec.Containers {
		for _, envFrom := range c.EnvFrom {
			if envFrom.ConfigMapRef != nil {
				r.recordConfigMap(envFrom.ConfigMapRef.Name)
			}
			if envFrom.SecretRef != nil {
				r.recordSecret(envFrom.SecretRef.Name)
			}
		}
	}

	for _, v := range podSpec.Volumes {
		if v.ConfigMap != nil {
			r.recordConfigMap(v.ConfigMap.Name)
		}
		if v.Secret != nil {
			r.recordSecret(v.ConfigMap.Name)
		}
	}
}

func (r *reload) recordConfigMap(name string) {
	if r.configMaps == nil {
		r.configMaps = map[string]bool{}
	}
	r.configMaps[name] = true
}

func (r *reload) MarshalTo(o kubeutil.AnnotationsAccessor) error {
	if len(r.configMaps) > 0 {
		if err := AnnotationReloadConfigMap.MarshalTo(o, keys(r.configMaps), ""); err != nil {
			return err
		}

		if r.configMapDigests != nil {
			for name := range r.configMaps {
				if digest, ok := r.configMapDigests[name]; ok {
					if err := ScopeConfigMapDigest.MustAnnotation(name).MarshalTo(o, digest, ""); err != nil {
						return err
					}
				}
			}
		}
	}
	if len(r.secrets) > 0 {
		if err := AnnotationReloadSecret.MarshalTo(o, keys(r.secrets), ""); err != nil {
			return err
		}

		if r.secretDigests != nil {
			for name := range r.secrets {
				if digest, ok := r.secretDigests[name]; ok {
					if err := ScopeSecretDigest.MustAnnotation(name).MarshalTo(o, digest, ""); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func keys[V any](m map[string]V) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}
