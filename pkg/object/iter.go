package object

import (
	"iter"
	"sync/atomic"
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

func (w *objectIter) normalizeAndProcess(o Object) (Object, error) {
	typed, err := Convert(o)
	if err != nil {
		return nil, err
	}

	if typed != o {
		o = typed
	}

	for i := range w.progress {
		o, err = w.progress[i](o)
		if err != nil {
			return nil, err
		}
	}
	return o, nil
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
		obj, ok := Is(x)
		if ok {
			return func(yield func(Object) bool) {
				normalized, err := w.normalizeAndProcess(obj)
				if err != nil {
					w.Done(err)
					return
				}
				if !yield(normalized) {
					return
				}
			}
		}

		return func(yield func(Object) bool) {
			// walk child
			for k, value := range x {
				o, ok := Is(value)
				if ok {
					normalized, err := w.normalizeAndProcess(o)
					if err != nil {
						w.Done(err)
						return
					}

					if !yield(normalized) {
						return
					}

					x[k] = normalized
					continue
				}

				for o := range w.iter(value) {
					if !yield(o) {
						return
					}
				}
			}
		}
	case []any:
		return func(yield func(Object) bool) {
			for i, value := range x {
				o, ok := Is(value)
				if ok {
					normalized, err := w.normalizeAndProcess(o)
					if err != nil {
						w.Done(err)
						return
					}

					if !yield(normalized) {
						return
					}

					x[i] = normalized
					continue
				}

				for o := range w.iter(value) {
					if !yield(o) {
						return
					}
				}
			}
		}
	}
	return func(yield func(Object) bool) {
	}
}
