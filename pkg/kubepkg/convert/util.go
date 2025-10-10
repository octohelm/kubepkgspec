package convert

import (
	"io"
	"strings"

	"github.com/go-json-experiment/json"
	jsonv1 "github.com/go-json-experiment/json/v1"
	"github.com/octohelm/kubekit/pkg/metadata"
	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/wellknown"
	"github.com/octohelm/x/anyjson"
	"golang.org/x/sync/errgroup"
)

func LabelInstanceAndVersion(kpkg *v1alpha1.KubePkg, o metadata.LabelsAccessor) error {
	if err := wellknown.LabelAppInstance.SetTo(o, kpkg.Name); err != nil {
		return err
	}
	if err := wellknown.LabelAppVersion.SetTo(o, kpkg.Spec.Version); err != nil {
		return err
	}
	return nil
}

func Must[T any](v *T) *T {
	if v == nil {
		v = new(T)
	}
	return v
}

func IsGlobalRef(name string) bool {
	return name != "" && (name[0] == '~' || name[0] == '/' || name[0] == '@')
}

func SubResourceName(kpkg *v1alpha1.KubePkg, name string) string {
	if name == "#" || name == "" {
		return kpkg.Name
	}
	if IsGlobalRef(name) {
		return name[1:]
	}
	return strings.Join([]string{kpkg.Name, name}, "-")
}

func Intersection[E comparable](a []E, b []E) (c []E) {
	includes := map[E]bool{}
	for i := range a {
		includes[a[i]] = true
	}

	c = make([]E, 0, len(a)+len(b))
	for i := range b {
		x := b[i]

		if _, ok := includes[x]; ok {
			c = append(c, x)
		}
	}

	return c
}

func Merge[X any](from *X, overwrites *X) (*X, error) {
	srcObj, err := anyjson.FromValue(from)
	if err != nil {
		return nil, err
	}

	overwritesObj, err := anyjson.FromValue(overwrites)
	if err != nil {
		return nil, err
	}

	merged := anyjson.Merge(srcObj, overwritesObj, anyjson.WithEmptyObjectAsNull(), anyjson.WithArrayMergeKey("name"))

	m := new(X)
	if err := unmarshal(merged, m); err != nil {
		return nil, err
	}
	return m, nil
}

func Unmarshal[X any](src any, target *X) error {
	srcObj, err := anyjson.FromValue(src)
	if err != nil {
		return err
	}
	merged := anyjson.Merge(
		anyjson.Valuer(&anyjson.Object{}),
		srcObj,
		anyjson.WithEmptyObjectAsNull(),
		anyjson.WithArrayMergeKey("name"),
	)

	return unmarshal(merged, target)
}

func unmarshal(src any, target any) error {
	r, w := io.Pipe()
	eg := &errgroup.Group{}

	eg.Go(func() error {
		defer w.Close()
		return json.MarshalWrite(w, src, jsonv1.OmitEmptyWithLegacySemantics(true))
	})

	eg.Go(func() error {
		defer r.Close()
		return json.UnmarshalRead(r, target, jsonv1.OmitEmptyWithLegacySemantics(true))
	})

	return eg.Wait()
}
