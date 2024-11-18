package wellknown

import "github.com/octohelm/kubekit/pkg/metadata"

var (
	LabelAppName     = metadata.Label("app.kubernetes.io/name")
	LabelAppInstance = metadata.Label("app.kubernetes.io/instance")
	LabelAppVersion  = metadata.Label("app.kubernetes.io/version")
	LabelAppPartOf   = metadata.Label("app.kubernetes.io/part-of")
)
