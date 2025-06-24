package v1alpha1

import (
	"github.com/go-json-experiment/json"
	"github.com/octohelm/courier/pkg/validator"
)

type StringOrSlice []string

func (StringOrSlice) OneOf() []any {
	return []any{"", []string{}}
}

func (d StringOrSlice) IsZero() bool {
	return len(d) == 0
}

func (d *StringOrSlice) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		b := ""
		if err := json.Unmarshal(data, &b); err != nil {
			return err
		}
		*d = []string{b}
		return nil
	}
	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	*d = list
	return nil
}

func (d StringOrSlice) MarshalJSON() ([]byte, error) {
	if len(d) == 0 {
		return nil, nil
	}
	if len(d) == 1 {
		return validator.Marshal(d[0])
	}
	return validator.Marshal([]string(d))
}

type StringSlice []string

func (d StringSlice) IsZero() bool {
	return len(d) == 0
}

func (d *StringSlice) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		b := ""
		if err := json.Unmarshal(data, &b); err != nil {
			return err
		}
		*d = []string{b}
		return nil
	}
	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	*d = list
	return nil
}

func (d StringSlice) MarshalJSON() ([]byte, error) {
	return validator.Marshal([]string(d))
}
