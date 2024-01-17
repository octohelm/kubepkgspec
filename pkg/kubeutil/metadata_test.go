package kubeutil

import (
	"testing"

	testingx "github.com/octohelm/x/testing"
	"github.com/opencontainers/go-digest"
)

func TestMetadata(t *testing.T) {
	s := Scope("app.kubernetes.io")

	t.Run("Label", func(t *testing.T) {
		labels := Labels{}

		err := s.MustLabel("name").SetTo(labels, "x")
		testingx.Expect(t, err, testingx.Be[error](nil))

		err = s.MustLabel("component").SetTo(labels, "database")
		testingx.Expect(t, err, testingx.Be[error](nil))

		testingx.Expect(t, labels, testingx.Equal(Labels{
			"app.kubernetes.io/name":      "x",
			"app.kubernetes.io/component": "database",
		}))
	})

	t.Run("annotations", func(t *testing.T) {
		annotations := Annotations{}

		err := s.MustAnnotation("digest").MarshalTo(annotations, digest.FromString(""), "")
		testingx.Expect(t, err, testingx.Be[error](nil))

		err = s.MustAnnotation("platform").MarshalTo(annotations, []string{"linux/arm64", "linux/amd64"}, "")
		testingx.Expect(t, err, testingx.Be[error](nil))

		err = s.MustAnnotation("ports").MarshalTo(annotations, map[string]any{"meter.industai.com": "http-80"}, "toml")
		testingx.Expect(t, err, testingx.Be[error](nil))

		testingx.Expect(t, annotations, testingx.Equal(Annotations{
			"app.kubernetes.io/platform": "linux/arm64,linux/amd64",
			"app.kubernetes.io/digest":   "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			"app.kubernetes.io/ports":    "toml:'meter.industai.com' = 'http-80'\n",
		}))

		t.Run("UnmarshalFrom csv", func(t *testing.T) {
			var platform []string
			err := s.MustAnnotation("platform").UnmarshalFrom(annotations, &platform)
			testingx.Expect(t, err, testingx.Be[error](nil))
			testingx.Expect(t, platform, testingx.Equal([]string{"linux/arm64", "linux/amd64"}))
		})

		t.Run("UnmarshalFrom toml", func(t *testing.T) {
			retPorts := map[string]any{}
			err := s.MustAnnotation("ports").UnmarshalFrom(annotations, &retPorts)
			testingx.Expect(t, err, testingx.Be[error](nil))
			testingx.Expect(t, retPorts, testingx.Equal(map[string]any{"meter.industai.com": "http-80"}))
		})
	})
}

type Labels map[string]string

func (l Labels) GetLabels() map[string]string {
	return l
}

func (l Labels) SetLabels(labels map[string]string) {
	for k, v := range labels {
		l[k] = v
	}
}

type Annotations map[string]string

func (annotations Annotations) GetAnnotations() map[string]string {
	return annotations
}

func (annotations Annotations) SetAnnotations(values map[string]string) {
	for k, v := range values {
		annotations[k] = v
	}
}
