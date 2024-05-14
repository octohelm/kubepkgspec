package manifest

import (
	"github.com/octohelm/kubepkgspec/pkg/install"
	"github.com/octohelm/kubepkgspec/pkg/kubepkg"
	"github.com/octohelm/kubepkgspec/pkg/reload"
	"iter"

	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/object"
	corev1 "k8s.io/api/core/v1"
)

func SortedExtract(kpkg *v1alpha1.KubePkg, options ...kubepkg.Option) ([]object.Object, error) {
	manifests, err := Extract(kpkg, options...)
	if err != nil {
		return nil, err
	}

	list := make([]object.Object, 0)

	if namespace := kpkg.Namespace; namespace != "" {
		n := &corev1.Namespace{}
		n.APIVersion = "v1"
		n.Kind = "Namespace"
		n.Name = namespace
		list = append(list, n)
	}

	for m := range manifests {
		list = append(list, m)
	}

	return install.SortByKind(list), nil
}

func Extract(kpkg *v1alpha1.KubePkg, options ...kubepkg.Option) (iter.Seq[object.Object], error) {
	manifests, err := kubepkg.Convert(kpkg, options...)
	if err != nil {
		return nil, err
	}
	if err := reload.Patch(manifests); err != nil {
		return nil, err
	}
	return manifests, nil
}
