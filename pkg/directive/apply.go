package directive

import (
	"context"
	"errors"
)

func Apply[V ~string](ctx context.Context, v V) (bool, V, error) {
	expr, err := parseExpr(string(v))
	if err != nil {
		if errors.Is(err, ErrInvalid) {
			return true, v, ErrInvalid
		}
	}

	if expr != nil {
		patched, err := expr.Apply(ctx)

		return true, V(patched), err
	}

	return false, v, nil
}
