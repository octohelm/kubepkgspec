package directive

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

var (
	ErrUnsupported = errors.New("directive is unsupported")
	ErrInvalid     = errors.New("directive is invalid")
)

type Directive interface {
	Group() string
	Name() string
	Apply(ctx context.Context, args Args) (string, error)
}

type Args map[string]any

func (args Args) String() string {
	if len(args) == 0 {
		return ""
	}

	s := &strings.Builder{}
	s.WriteString("{")

	for i, k := range slices.Sorted(maps.Keys(args)) {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(k)
		s.WriteString("=")

		switch v := args[k].(type) {
		case string:
			s.WriteString(fmt.Sprintf("%q", v))
		default:
			if v == nil {
				s.WriteString("null")
			} else {
				s.WriteString(fmt.Sprintf("%v", v))
			}
		}
	}

	s.WriteString("}")

	return s.String()
}

var directives = map[string]map[string]Directive{}

func Register(d Directive) {
	if directives[d.Group()] == nil {
		directives[d.Group()] = map[string]Directive{}
	}
	directives[d.Group()][d.Name()] = d
}

type Expr struct {
	Group string
	Name  string
	Args  Args
}

func (e *Expr) Apply(ctx context.Context) (string, error) {
	if subs, ok := directives[e.Group]; ok {
		if d, ok := subs[e.Name]; ok {
			return d.Apply(ctx, e.Args)
		}
	}
	return "", fmt.Errorf("%w: group=%s, name=%s", ErrUnsupported, e.Group, e.Name)
}

const (
	argFlagNone = iota
	argFlagKey
	argFlagValue
)

func parseExpr(directive string) (*Expr, error) {
	if strings.HasPrefix(directive, "/") {
		ss := &scanner.Scanner{}
		ss.Init(bytes.NewBufferString(directive[1:]))
		ss.IsIdentRune = func(ch rune, i int) bool {
			return unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0 || (ch == '-' || ch == '_') && i > 0
		}

		expr := &Expr{
			Args: Args{},
		}

		argFlag := argFlagNone

		key := ""

		commitArg := func(keyOrValue string) error {
			switch argFlag {
			case argFlagKey:
				key = keyOrValue
			case argFlagValue:
				jsonText := jsontext.Value(keyOrValue)
				if !jsonText.IsValid() {
					return fmt.Errorf("%w: invalid value: %s", ErrInvalid, keyOrValue)
				}

				switch jsonText.Kind() {
				case '"':
					var v string
					if err := json.Unmarshal(jsonText, &v); err != nil {
						return err
					}
					expr.Args[key] = v
				case '0':
					if bytes.Contains(jsonText, []byte(".")) {
						var v float64
						if err := json.Unmarshal(jsonText, &v); err != nil {
							return err
						}
						expr.Args[key] = v
						return nil
					}

					var v int64
					if err := json.Unmarshal(jsonText, &v); err != nil {
						return err
					}
					expr.Args[key] = v
				case 'f':
					expr.Args[key] = false
				case 't':
					expr.Args[key] = true
				default:
					expr.Args[key] = nil
				}

				argFlag = argFlagKey
			default:
			}

			return nil
		}

		for t := ss.Scan(); t != scanner.EOF; t = ss.Scan() {
			switch t {
			case '{':
				argFlag = argFlagKey
			case '=':
				argFlag = argFlagValue
			case ',':
				argFlag = argFlagKey
			case '}':
				argFlag = argFlagNone
			default:
				tokenText := ss.TokenText()

				if argFlag != argFlagNone {
					if err := commitArg(tokenText); err != nil {
						return nil, err
					}
					continue
				}

				if expr.Group == "" {
					expr.Group = tokenText
					continue
				}

				if expr.Name == "" {
					expr.Name = tokenText
				}
			}
		}

		if expr.Group == "" || expr.Name == "" {
			return nil, fmt.Errorf("%w: %s", ErrInvalid, directive)
		}

		return expr, nil
	}

	return nil, nil
}
