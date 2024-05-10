package iterutil

import "iter"

func Items[O any](list []O) iter.Seq[O] {
	return func(yield func(O) bool) {
		for _, x := range list {
			if !yield(x) {
				return
			}
		}
	}
}

func ToSlice[O any](s iter.Seq[O]) []O {
	list := make([]O, 0)
	for x := range s {
		list = append(list, x)
	}
	return list
}
