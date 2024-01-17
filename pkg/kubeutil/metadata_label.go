package kubeutil

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type LabelsAccessor interface {
	GetLabels() map[string]string
	SetLabels(labels map[string]string)
}

func ParseLabel(label string) (Label, error) {
	if err := ValidateLabel(label); err != nil {
		return "", err
	}
	return Label(label), nil
}

type Label string

func (l Label) Requirement(op selection.Operator, values ...string) *labels.Requirement {
	r, _ := labels.NewRequirement(string(l), op, values)
	return r
}

func (l Label) GetFrom(o LabelsAccessor) string {
	label := string(l)

	labels := o.GetLabels()
	if labels != nil {
		if a, ok := labels[label]; ok {
			return a
		}
	}
	return ""
}

func (l Label) SetTo(o LabelsAccessor, value string) error {
	label := string(l)

	labels := o.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	if value != "" {
		labels[label] = value
	} else {
		delete(labels, label)
	}
	o.SetLabels(labels)

	return nil
}
