package object

import (
	"github.com/octohelm/x/anyjson"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	Object = client.Object
	List   = client.ObjectList
)

func FromRuntimeObject(ro runtime.Object) (Object, error) {
	o, err := meta.Accessor(ro)
	if err != nil {
		return nil, err
	}
	return o.(Object), nil
}

func ListFromRuntimeObject(ro runtime.Object) (List, error) {
	o, err := meta.ListAccessor(ro)
	if err != nil {
		return nil, err
	}
	return o.(List), nil
}

func Is(value any) (Object, bool) {
	switch x := value.(type) {
	case Object:
		return x, true
	case map[string]any:
		if isKubernetesManifest(x) {
			o, err := FromRuntimeObject(&unstructured.Unstructured{Object: x})
			if err != nil {
				return nil, false
			}
			return o, true
		}
	}

	return nil, false
}

func isKubernetesManifest(obj anyjson.Obj) bool {
	if _, ok := obj["apiVersion"].(string); ok {
		if _, ok := obj["kind"].(string); ok {
			return true
		}
	}
	return false
}
