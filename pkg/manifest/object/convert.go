package object

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes/scheme"
)

func Convert(o Object) (Object, error) {
	gvk := o.GetObjectKind().GroupVersionKind()

	typed, err := scheme.Scheme.New(gvk)
	if err == nil {
		if err := scheme.Scheme.Convert(o, typed, nil); err != nil {
			return nil, err
		}

		stableGV := scheme.Scheme.VersionsForGroupKind(gvk.GroupKind())[0]
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
