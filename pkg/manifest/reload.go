package manifest

import (
	"encoding/json"
	"sort"

	"github.com/octohelm/kubekit/pkg/metadata"
	"github.com/opencontainers/go-digest"
	corev1 "k8s.io/api/core/v1"
)

var AnnotationReloadConfigMap = metadata.Annotation("reload.octohelm.tech/configmap")
var AnnotationReloadSecret = metadata.Annotation("reload.octohelm.tech/secret")

var ScopeConfigMapDigest = metadata.Scope("digest.configmap.octohelm.tech")
var ScopeSecretDigest = metadata.Scope("digest.secret.octohelm.tech")

func AnnotateDigestTo(o metadata.AnnotationsAccessor, scope metadata.Scope, name string, value any) error {
	a, err := scope.Annotation(name)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return a.MarshalTo(o, digest.FromBytes(raw), "")
}

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

func (r *reload) MarshalTo(o metadata.AnnotationsAccessor) error {
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
