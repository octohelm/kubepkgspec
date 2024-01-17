package manifest

import (
	"encoding/json"
	"github.com/octohelm/kubepkgspec/pkg/kubeutil"
	"github.com/opencontainers/go-digest"
)

var AnnotationIngressGateway = kubeutil.Annotation("ingress.octohelm.tech/gateway")
var AnnotationPlatform = kubeutil.Annotation("octohelm.tech/platform")

var AnnotationReloadConfigMap = kubeutil.Annotation("reload.octohelm.tech/configmap")
var AnnotationReloadSecret = kubeutil.Annotation("reload.octohelm.tech/secret")

var ScopeConfigMapDigest = kubeutil.Scope("digest.configmap.octohelm.tech")
var ScopeSecretDigest = kubeutil.Scope("digest.secret.octohelm.tech")

func AnnotateDigestTo(o kubeutil.AnnotationsAccessor, scope kubeutil.Scope, name string, value any) error {
	a, err := scope.Annotation(name)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(value)
	return a.MarshalTo(o, digest.FromBytes(raw), "")
}
