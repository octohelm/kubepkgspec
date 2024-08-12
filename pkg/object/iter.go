package object

import (
	"iter"
	"sync/atomic"

	"github.com/octohelm/x/anyjson"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Process = func(o Object) (Object, error)

func NewIter(progress ...Process) Iter {
	return &objectIter{progress: progress}
}

type Iter interface {
	Object(m any) iter.Seq[Object]
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

func (w *objectIter) Object(m any) iter.Seq[Object] {
	return func(yield func(Object) bool) {
		for o := range w.iter(m) {
			if w.Err() != nil {
				return
			}

			if !yield(o) {
				return
			}
		}
	}
}

func (w *objectIter) iter(v any) iter.Seq[Object] {
	switch x := v.(type) {
	case Object:
		return func(yield func(Object) bool) {
			if !yield(x) {
				return
			}
		}
	case map[string]any:
		return w.iterObj(x)
	case []any:
		return w.iterList(x)
	}
	return func(yield func(Object) bool) {
	}
}

func (w *objectIter) iterList(list []any) iter.Seq[Object] {
	return func(yield func(Object) bool) {
		for _, value := range list {
			for o := range w.iter(value) {
				if !yield(o) {
				}
			}
		}
	}
}

func (w *objectIter) iterObj(obj anyjson.Obj) iter.Seq[Object] {
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
					return
				}
			}

			if !yield(o) {
				return
			}

			return
		}

		// walk child
		for _, value := range obj {
			for o := range w.iter(value) {
				if !yield(o) {
					return
				}
			}
		}
	}
}
