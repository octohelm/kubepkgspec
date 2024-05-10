package kubepkg

import (
	"github.com/octohelm/kubekit/pkg/metadata"
)

var (
	AnnotationIngressGateway = metadata.Annotation("ingress.octohelm.tech/gateway")
	AnnotationPlatform       = metadata.Annotation("octohelm.tech/platform")
)
