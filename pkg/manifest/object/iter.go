package object

import (
	"context"
	"iter"
	"sync/atomic"

	"github.com/stretchr/objx"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Process = func(o Object) (Object, error)

func NewIter(progress ...Process) Iter {
	return &objectIter{progress: progress}
}

type Iter interface {
	Object(ctx context.Context, m any) iter.Seq[Object]
	Err() error
}

type objectIter struct {
	progress []Process
	err      atomic.Value
}

func (w *objectIter) Done(err error) {
	if err != nil {
		w.err.Store(err)
	}
}

func (w *objectIter) Err() error {
	v := w.err.Load()
	if v == nil {
		return nil
	}
	return v.(error)
}

func (w *objectIter) Object(ctx context.Context, m any) iter.Seq[Object] {
	return func(yield func(Object) bool) {
		for o := range w.iter(ctx, m) {
			if w.Err() != nil {
				return
			}

			if !yield(o) {
				return
			}
		}
	}
}

func (w *objectIter) iter(ctx context.Context, v any) iter.Seq[Object] {
	switch x := v.(type) {
	case map[string]any:
		return w.iterObj(ctx, x)
	case []any:
		return w.iterList(ctx, x)
	}
	return func(yield func(Object) bool) {
	}
}

func (w *objectIter) iterList(ctx context.Context, list []any) iter.Seq[Object] {
	return func(yield func(Object) bool) {
		for _, value := range list {
			for o := range w.iter(ctx, value) {
				if !yield(o) {
				}
			}
		}
	}
}

func (w *objectIter) iterObj(ctx context.Context, obj objx.Map) iter.Seq[Object] {
	return func(yield func(Object) bool) {
		if isKubernetesManifest(obj) {
			o, err := FromRuntimeObject(&unstructured.Unstructured{Object: obj})
			if err != nil {
				w.Done(err)
				return
			}

			typed, err := Convert(o)
			if err != nil {
				w.Done(err)
				return
			}
			if typed != o {
				o = typed
			}

			for i := range w.progress {
				o, err = w.progress[i](o)
				if err != nil {
					w.Done(err)
				}
			}

			if !yield(o) {
				return
			}

			return
		}

		// walk child
		for _, value := range obj {
			for o := range w.iter(ctx, value) {
				if !yield(o) {
					return
				}
			}
		}
	}
}

func isKubernetesManifest(obj objx.Map) bool {
	return obj.Get("apiVersion").IsStr() &&
		obj.Get("apiVersion").Str() != "" &&
		obj.Get("kind").IsStr() &&
		obj.Get("kind").Str() != ""
}
