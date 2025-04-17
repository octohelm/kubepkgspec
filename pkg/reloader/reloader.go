package reloader

import (
	"maps"
	"slices"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/octohelm/kubekit/pkg/metadata"
	"github.com/opencontainers/go-digest"
)

var (
	AnnotationConfigMapReload = metadata.Annotation("configmap.reloader.octohelm.tech/reload")
	AnnotationSecretReload    = metadata.Annotation("secret.reloader.octohelm.tech/reload")
)

const (
	AnnotationSpecDigest     = metadata.Annotation("spec/digest")
	AnnotationRevisionDigest = metadata.Annotation("revision/digest")
)

var (
	ScopeConfigMapDigest = metadata.Scope("digest.configmap.reloader.octohelm.tech")
	ScopeSecretDigest    = metadata.Scope("digest.secret.reloader.octohelm.tech")
)

func DigestFromValues[Values map[string]string](values Values) (digest.Digest, error) {
	hasher := digest.SHA256.Digester()

	enc := jsontext.NewEncoder(hasher.Hash())

	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return "", err
	}

	for _, k := range slices.Sorted(maps.Keys(values)) {
		if err := enc.WriteToken(jsontext.String(k)); err != nil {
			return "", err
		}

		if err := enc.WriteToken(jsontext.String(values[k])); err != nil {
			return "", err
		}
	}

	if err := enc.WriteToken(jsontext.EndObject); err != nil {
		return "", err
	}

	return hasher.Digest(), nil
}

func AnnotateDigestTo[Values map[string]string](o metadata.AnnotationsAccessor, values Values) error {
	dgst, err := DigestFromValues(values)
	if err != nil {
		return err
	}
	return AnnotationSpecDigest.MarshalTo(o, dgst, "")
}
