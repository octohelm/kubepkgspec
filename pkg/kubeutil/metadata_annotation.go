package kubeutil

import (
	"bufio"
	"bytes"
	"reflect"
	"strings"

	encodingx "github.com/octohelm/x/encoding"
	"github.com/pelletier/go-toml/v2"
)

type AnnotationsAccessor interface {
	GetAnnotations() map[string]string
	SetAnnotations(annotations map[string]string)
}

func ParseAnnotation(label string) (Annotation, error) {
	if err := ValidateLabel(label); err != nil {
		return "", err
	}
	return Annotation(label), nil
}

type Annotation string

func (l Annotation) Get(o AnnotationsAccessor) (string, bool) {
	annotations := o.GetAnnotations()
	if annotations != nil {
		annotation := string(l)
		value, ok := annotations[annotation]
		return value, ok
	}

	return "", false
}

func (l Annotation) UnmarshalFrom(o AnnotationsAccessor, v any) error {
	annotation := string(l)

	annotations := o.GetAnnotations()
	if annotations != nil {
		if value, ok := annotations[annotation]; ok {
			if strings.HasPrefix(value, "toml:") {
				return toml.Unmarshal([]byte(value[len("toml:"):]), v)
			} else {
				rv := reflect.ValueOf(v)

				slice := rv
				for slice.Kind() == reflect.Ptr {
					slice = slice.Elem()
				}

				if slice.Kind() == reflect.Slice {
					s := bufio.NewScanner(bytes.NewBuffer([]byte(value)))
					s.Split(onComma)

					for s.Scan() {
						item := reflect.New(slice.Type().Elem())

						if err := encodingx.UnmarshalText(item, []byte(s.Text())); err != nil {
							return err
						}

						slice = reflect.Append(slice, item.Elem())
					}

					rv.Elem().Set(slice)
				} else {
					return encodingx.UnmarshalText(v, []byte(value))
				}
			}
		}
	}

	return nil
}

func onComma(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == ',' {
			return i + 1, data[:i], nil
		}
	}
	if !atEOF {
		return 0, nil, nil
	}
	return 0, data, bufio.ErrFinalToken
}

func (l Annotation) MarshalTo(o AnnotationsAccessor, v any, style string) error {
	annotation := string(l)

	annotations := o.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	b := &strings.Builder{}

	switch style {
	case "toml":
		b.WriteString("toml:")

		enc := toml.NewEncoder(b)
		err := enc.Encode(v)
		if err != nil {
			return err
		}
	default:
		switch x := v.(type) {
		case []byte:
			b.Write(x)
		default:
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}

			if rv.Kind() == reflect.Slice {
				for i := 0; i < rv.Len(); i++ {
					if i > 0 {
						b.WriteString(",")
					}
					data, err := encodingx.MarshalText(rv.Index(i))
					if err != nil {
						return err
					}
					b.Write(data)
				}
			} else {
				data, err := encodingx.MarshalText(v)
				if err != nil {
					return err
				}
				b.Write(data)
			}
		}
	}

	if b.Len() > 0 {
		annotations[annotation] = b.String()
	} else {
		delete(annotations, annotation)
	}

	o.SetAnnotations(annotations)

	return nil
}
