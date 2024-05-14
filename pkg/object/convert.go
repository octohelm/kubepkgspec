package object

import (
	"encoding/json"
	kubepkgv1alpha1 "github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

var Scheme = runtime.NewScheme()

func init() {
	_ = scheme.AddToScheme(Scheme)
	_ = apiextensionsv1.AddToScheme(Scheme)
	_ = kubepkgv1alpha1.AddToScheme(Scheme)
}

func Convert(o Object) (Object, error) {
	gvk := o.GetObjectKind().GroupVersionKind()

	typed, err := Scheme.New(gvk)
	if err == nil {
		raw, err := json.Marshal(o)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(raw, typed); err != nil {
			return nil, errors.Wrap(err, "convert failed")
		}

		stableGV := Scheme.VersionsForGroupKind(gvk.GroupKind())[0]
		if gvk.Version != stableGV.Version {
			return nil, errors.Errorf("unsupport gvk %s, please upgrade to %s", gvk, stableGV.WithKind(gvk.Kind))
		}

		typedObj, err := FromRuntimeObject(typed)
		if err != nil {
			return nil, err
		}

		typedObj.GetObjectKind().SetGroupVersionKind(gvk)

		return typedObj, nil
	}

	return o, nil
}
