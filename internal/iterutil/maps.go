package iterutil

import "iter"

func Values[O any](list map[string]O) iter.Seq[O] {
	return func(yield func(O) bool) {
		for _, x := range list {
			if !yield(x) {
				return
			}
		}
	}
}
