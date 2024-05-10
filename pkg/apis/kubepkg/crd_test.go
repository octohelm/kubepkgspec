package kubepkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	cueformat "cuelang.org/go/cue/format"
	testingx "github.com/octohelm/x/testing"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"

	"sigs.k8s.io/yaml"

	"github.com/octohelm/courier/pkg/openapi/jsonschema"
	"github.com/octohelm/courier/pkg/openapi/jsonschema/extractors"
	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
)

func TestCustomResourceDefinition(t *testing.T) {
	data, _ := yaml.Marshal(CRDs)
	_ = os.WriteFile("crd.yaml", data, 0o644)
}

func TestCueify(t *testing.T) {
	scanner := &jsonSchemaScanner{}

	schema := scanner.ExtractSchema(context.Background(), &v1alpha1.KubePkg{})

	t.Run("should write json schema", func(t *testing.T) {
		data, err := json.MarshalIndent(schema, "", "  ")
		testingx.Expect(t, err, testingx.BeNil[error]())

		err = os.WriteFile("schema.json", data, 0o644)
		testingx.Expect(t, err, testingx.BeNil[error]())
	})

	t.Run("should write cue", func(t *testing.T) {
		b := bytes.NewBuffer(nil)
		_, _ = fmt.Fprintln(b, "package kubepkg")

		names := make([]string, 0, len(scanner.Defs))
		for n := range scanner.Defs {
			names = append(names, n)
		}
		sort.Strings(names)

		for _, n := range names {
			def := scanner.Defs[n]
			_, _ = fmt.Fprintln(b)
			_, _ = fmt.Fprintf(b, "#%s: ", n)
			def.PrintTo(b, jsonschema.PrintWithDoc())
			_, _ = fmt.Fprintln(b)
		}

		data, err := cueformat.Source(b.Bytes(), cueformat.Simplify())
		testingx.Expect(t, err, testingx.BeNil[error]())
		_ = os.WriteFile("../../../cuepkg/kubepkg/spec.cue", data, 0o644)
	})
}

type jsonSchemaScanner struct {
	Defs map[string]jsonschema.Schema
	m    sync.Map
}

func (s *jsonSchemaScanner) ExtractSchema(ctx context.Context, target any) *jsonschema.Payload {
	schema := extractors.SchemaFrom(extractors.SchemaRegisterContext.Inject(ctx, s), target, true)

	definitions := map[string]jsonschema.Schema{}

	for n, sub := range s.Defs {
		if schema != sub {
			definitions[n] = s.Defs[n]
		}
	}
	schema.GetMetadata().AddExtension("$defs", definitions)

	return &jsonschema.Payload{
		Schema: schema,
	}
}

const schemaPathPrefix = "#/$defs/"

func (s *jsonSchemaScanner) Record(typeRef string) bool {
	_, ok := s.m.Load(typeRef)
	defer s.m.Store(typeRef, true)
	return ok
}

func (s *jsonSchemaScanner) RefString(ref string) string {
	parts := strings.Split(ref, ".")
	return fmt.Sprintf("%s%s", schemaPathPrefix, parts[len(parts)-1])
}

func (s *jsonSchemaScanner) RegisterSchema(ref string, schema jsonschema.Schema) {
	if s.Defs == nil {
		s.Defs = map[string]jsonschema.Schema{}
	}

	refName := strings.TrimPrefix(ref, schemaPathPrefix)

	if _, ok := s.Defs[refName]; ok {
		fmt.Printf("%s already defined\n", refName)
	} else {
		s.Defs[refName] = schema
	}
}
