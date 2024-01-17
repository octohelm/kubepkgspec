package kubeutil

import (
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation"
	"strings"
)

var (
	ErrInvalidKey = errors.New("invalid key [<prefix>/]<name>")
)

func ValidateLabel(label string) error {
	name := label
	prefix := ""
	if i := strings.Index(label, "/"); i > 0 {
		prefix = label[0:i]
		name = name[i+1:]
	}

	if errs := validation.IsDNS1123Subdomain(prefix); len(errs) > 0 {
		return errors.Wrapf(ErrInvalidKey, "prefix: %s: %s", errs[0], prefix)
	}

	if errs := validation.IsValidLabelValue(name); len(errs) > 0 {
		return errors.Wrapf(ErrInvalidKey, "name: %s: %s", errs[0], name)
	}

	return nil
}

type Scope string

func (s Scope) MustLabel(name string) Label {
	l, err := ParseLabel(fmt.Sprintf("%s/%s", s, name))
	if err != nil {
		panic(err)
	}
	return l
}

func (s Scope) MustAnnotation(name string) Annotation {
	annotation, err := ParseAnnotation(fmt.Sprintf("%s/%s", s, name))
	if err != nil {
		panic(err)
	}
	return annotation
}

func (s Scope) Label(name string) (Label, error) {
	return ParseLabel(fmt.Sprintf("%s/%s", s, name))
}

func (s Scope) Annotation(name string) (Annotation, error) {
	return ParseAnnotation(fmt.Sprintf("%s/%s", s, name))
}
