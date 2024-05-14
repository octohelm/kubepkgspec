package manifest

import (
	_ "embed"
	"testing"

	"github.com/octohelm/kubepkgspec/pkg/kubepkg/convert"
	"github.com/octohelm/x/anyjson"

	"github.com/octohelm/kubepkgspec/pkg/install"

	"github.com/octohelm/kubepkgspec/internal/iterutil"
	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/kubepkg"
	"github.com/octohelm/kubepkgspec/pkg/object"
	"github.com/octohelm/kubepkgspec/pkg/workload"
	testingx "github.com/octohelm/x/testing"
	"golang.org/x/tools/txtar"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/yaml"
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

			s := &txtar.Archive{}
			for _, o := range list {
				s.Files = append(s.Files, asTxtTarFile(o))
			}
			testingx.Expect(t, s, testingx.MatchSnapshot("example.kubepkg"))

			t.Run("should extract as kubepkg from manifests", func(t *testing.T) {
				kpkg, err := kubepkg.ExtractAsKubePkg(iterutil.Items(list))
				testingx.Expect(t, err, testingx.BeNil[error]())

				list, err := SortedExtract(kpkg)
				s := &txtar.Archive{}
				for _, o := range list {
					s.Files = append(s.Files, asTxtTarFile(o.(object.Object)))
				}
				testingx.Expect(t, s, testingx.MatchSnapshot("example-simpled.kubepkg"))
			})

			t.Run("could change image names", func(t *testing.T) {
				renames := map[string]string{
					"docker.io/library/nginx": "library/nginx",
				}

				for img := range workload.Images(iterutil.Items(list)) {
					img.Name = renames[img.Name]
				}

				for c := range workload.Workloads(iterutil.Items(list)) {
					switch x := c.(type) {
					case *appsv1.Deployment:
						testingx.Expect(t,
							x.Spec.Template.Spec.Containers[0].Image,
							testingx.Be("library/nginx:1.25.0-alpine"),
						)
					}
				}
			})
		})
	})

	t.Run("from manifests", func(t *testing.T) {
		rawKPkg := &v1alpha1.KubePkg{}
		testingx.Expect(t, yaml.Unmarshal(kubepkgFromManifests, rawKPkg), testingx.BeNil[error]())

		manifests := iterutil.Items(install.SortByKind(iterutil.ToSlice(object.NewIter().Object(rawKPkg.Spec.Manifests))))

		kpkg, err := kubepkg.ExtractAsKubePkg(manifests)
		testingx.Expect(t, err, testingx.BeNil[error]())

		list, err := SortedExtract(kpkg)
		testingx.Expect(t, err, testingx.BeNil[error]())

		s := &txtar.Archive{}
		for _, o := range list {
			s.Files = append(s.Files, asTxtTarFile(o))
		}

		testingx.Expect(t, s, testingx.MatchSnapshot("from-manifests.kubepkg"))
	})
}

func asTxtTarFile(o object.Object) txtar.File {
	x := txtar.File{
		Name: o.GetName(),
	}

	if namespace := o.GetNamespace(); namespace != "" {
		x.Name += "." + namespace
	}

	x.Name += "." + o.GetObjectKind().GroupVersionKind().Kind + "." + "yaml"

	obj := &anyjson.Object{}
	if err := convert.Unmarshal(o, obj); err != nil {
		panic(err)
	}

	data, err := yaml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	x.Data = data
	return x
}
