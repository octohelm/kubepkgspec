package object

import (
	"fmt"

	"github.com/go-json-experiment/json"
	jsonv1 "github.com/go-json-experiment/json/v1"
)

func Convert(o Object) (Object, error) {
	gvk := o.GetObjectKind().GroupVersionKind()

	typed, err := Scheme.New(gvk)
	if err == nil {
		raw, err := json.Marshal(o, jsonv1.OmitEmptyWithLegacySemantics(true))
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(raw, typed, jsonv1.OmitEmptyWithLegacySemantics(true)); err != nil {
			return nil, fmt.Errorf("convert failed: %w", err)
		}

		stableGV := Scheme.VersionsForGroupKind(gvk.GroupKind())[0]
		if gvk.Version != stableGV.Version {
			return nil, fmt.Errorf("unsupport gvk %s, please upgrade to %s", gvk, stableGV.WithKind(gvk.Kind))
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
