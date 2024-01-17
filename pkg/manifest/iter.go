package manifest

import (
	"context"
	"sync/atomic"

	"github.com/stretchr/objx"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Iter interface {
	Manifest(ctx context.Context, m any) <-chan Object
	Err() error
}

type ObjectProcess = func(o Object) (Object, error)

func NewIter(progress ...ObjectProcess) Iter {
	return &iter{progress: progress}
}

type iter struct {
	progress []ObjectProcess
	err      atomic.Value
}

func (w *iter) Done(err error) {
	if err != nil {
		w.err.Store(err)
	}
}

func (w *iter) Err() error {
	v := w.err.Load()
	if v == nil {
		return nil
	}
	return v.(error)
}

func (w *iter) Manifest(ctx context.Context, m any) <-chan Object {
	ch := make(chan Object)

	go func() {
		defer close(ch)

		for o := range w.iter(ctx, m) {
			if w.Err() != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case ch <- o:
			}
		}
	}()

	return ch
}

func (w *iter) iter(ctx context.Context, v any) <-chan Object {
	switch x := v.(type) {
	case map[string]any:
		return w.iterObj(ctx, x)
	case []any:
		return w.iterList(ctx, x)
	}
	ch := make(chan Object)
	close(ch)
	return ch
}

func (w *iter) iterList(ctx context.Context, list []any) <-chan Object {
	ch := make(chan Object)

	go func() {
		defer close(ch)
		for _, value := range list {
			for o := range w.iter(ctx, value) {
				ch <- o
			}
		}
	}()

	return ch
}

func (w *iter) iterObj(ctx context.Context, obj objx.Map) <-chan Object {
	ch := make(chan Object)

	go func() {
		defer close(ch)

		if isKubernetesManifest(obj) {
			co, err := ObjectFromRuntimeObject(&unstructured.Unstructured{Object: obj})
			if err != nil {
				w.Done(err)
				return
			}

			for i := range w.progress {
				co, err = w.progress[i](co)
				if err != nil {
					w.Done(err)
				}
			}

			ch <- co

			return
		}

		for key := range obj {
			for o := range w.iter(ctx, obj[key]) {
				ch <- o
			}
		}
	}()

	return ch
}

func isKubernetesManifest(obj objx.Map) bool {
	return obj.Get("apiVersion").IsStr() &&
		obj.Get("apiVersion").Str() != "" &&
		obj.Get("kind").IsStr() &&
		obj.Get("kind").Str() != ""
}
