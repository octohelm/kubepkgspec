package manifest

import (
	"github.com/octohelm/kubepkgspec/internal/iterutil"
	"github.com/octohelm/kubepkgspec/pkg/install"
	"github.com/octohelm/kubepkgspec/pkg/kubepkg"
	"github.com/octohelm/kubepkgspec/pkg/reload"

	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/object"
	corev1 "k8s.io/api/core/v1"
)

func SortedExtract(kpkg *v1alpha1.KubePkg) ([]object.Object, error) {
	manifests, err := Extract(kpkg)
	if err != nil {
		return nil, err
	}

	list := make([]object.Object, 0, len(manifests)+1)

	if namespace := kpkg.Namespace; namespace != "" {
		n := &corev1.Namespace{}
		n.APIVersion = "v1"
		n.Kind = "Namespace"
		n.Name = namespace
		list = append(list, n)
	}

	for k := range manifests {
		list = append(list, manifests[k])
	}

	return install.SortByKind(list), nil
}

func Extract(kpkg *v1alpha1.KubePkg) (map[string]object.Object, error) {
	manifests, err := kubepkg.Convert(kpkg)
	if err != nil {
		return nil, err
	}
	if err := reload.Patch(iterutil.Values(manifests)); err != nil {
		return nil, err
	}
	return manifests, nil
}
