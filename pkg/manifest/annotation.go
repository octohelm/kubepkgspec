package manifest

import (
	"github.com/octohelm/kubekit/pkg/metadata"
)

var AnnotationIngressGateway = metadata.Annotation("ingress.octohelm.tech/gateway")
var AnnotationPlatform = metadata.Annotation("octohelm.tech/platform")
