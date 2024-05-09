package kubepkg

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/octohelm/kubekit/pkg/crd"
	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
)

var CRDs = []*apiextensionsv1.CustomResourceDefinition{
	crd.AsKubeCRD(&crd.CustomResourceDefinition{
		GroupVersion: v1alpha1.SchemeGroupVersion,
		KindType:     &v1alpha1.KubePkg{},
		ListKindType: &v1alpha1.KubePkgList{},
		SpecType:     &v1alpha1.Spec{},
	}),
}
