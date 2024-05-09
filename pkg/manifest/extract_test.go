package manifest

import (
	_ "embed"
	"fmt"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/yaml"

	testingx "github.com/octohelm/x/testing"

	"github.com/octohelm/kubepkgspec/internal/iterutil"
	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/manifest/workload"
)

//go:embed testdata/example.kubepkg.json
var kubepkgExample []byte

//go:embed testdata/from-manifests.kubepkg.json
var kubepkgFromManifests []byte

func TestExtract(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		t.Run("should extract raw manifests", func(t *testing.T) {
			kpkg := &v1alpha1.KubePkg{}
			testingx.Expect(t, yaml.Unmarshal(kubepkgExample, kpkg), testingx.BeNil[error]())

			list, err := SortedExtract(kpkg)
			testingx.Expect(t, err, testingx.BeNil[error]())

			names := make([]string, len(list))
			for i, o := range list {
				names[i] = fmt.Sprintf("%s.%s %s", o.GetName(), o.GetNamespace(), o.GetObjectKind().GroupVersionKind())
			}

			testingx.Expect(t, names, testingx.Equal([]string{
				"default. /v1, Kind=Namespace",
				"demo.default /v1, Kind=ConfigMap",
				"demo-html.default /v1, Kind=ConfigMap",
				"endpoint-demo.default /v1, Kind=ConfigMap",
				"demo.default /v1, Kind=Service",
				"demo.default apps/v1, Kind=Deployment",
				"demo.default networking.k8s.io/v1, Kind=Ingress",
			}))

			t.Run("could change image names", func(t *testing.T) {
				renames := map[string]string{
					"docker.io/library/nginx:1.25.0-alpine": "library/nginx:1.25.0-alpine",
				}

				for c := range workload.Containers(iterutil.Items(list)) {
					c.Image = renames[c.Image]
				}

				for c := range workload.Workloads(iterutil.Items(list)) {
					switch x := c.(type) {
					case *appsv1.Deployment:
						testingx.Expect(t, x.Spec.Template.Spec.Containers[0].Image, testingx.Be("library/nginx:1.25.0-alpine"))
					}
				}
			})
		})
	})

	t.Run("from manifests", func(t *testing.T) {
		kpkg := &v1alpha1.KubePkg{}
		testingx.Expect(t, yaml.Unmarshal(kubepkgFromManifests, kpkg), testingx.BeNil[error]())

		list, err := SortedExtract(kpkg)
		testingx.Expect(t, err, testingx.BeNil[error]())

		names := make([]string, len(list))
		for i, o := range list {
			names[i] = fmt.Sprintf("%s.%s %s", o.GetName(), o.GetNamespace(), o.GetObjectKind().GroupVersionKind())
		}

		testingx.Expect(t, names, testingx.Equal([]string{
			"device-system. /v1, Kind=Namespace",
			"gpu-feature-discovery.device-system /v1, Kind=ConfigMap",
			"gpu-feature-discovery-gpu-feature-discovery.device-system apps/v1, Kind=DaemonSet",
		}))
	})
}
