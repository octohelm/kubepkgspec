package v1alpha1

import "github.com/octohelm/courier/pkg/validator"

type StringOrSlice []string

func (StringOrSlice) OneOf() []any {
	return []any{"", []string{}}
}

func (d StringOrSlice) IsZero() bool {
	return len(d) == 0
}

func (d *StringOrSlice) UnmarshalJSON(data []byte) error {
	if len(data) > 0 {
		switch data[0] {
		case '"':
			value := ""
			if err := validator.Unmarshal(data, &value); err != nil {
				return err
			}
			*d = []string{value}
		case '[':
			var values []string
			if err := validator.Unmarshal(data, &values); err != nil {
				return err
			}
			*d = values
		}
	}
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
