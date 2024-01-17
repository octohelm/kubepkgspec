package manifest

import (
	"fmt"
	"testing"

	"sigs.k8s.io/yaml"

	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/x/ptr"
	testingx "github.com/octohelm/x/testing"
)

func TestExtract(t *testing.T) {
	t.Run("external", func(t *testing.T) {
		kpkg := &v1alpha1.KubePkg{}
		kpkg.Name = "demo"
		kpkg.Namespace = "default"
		kpkg.Annotations = map[string]string{
			"ingress.octohelm.tech/gateway": "public+https://{{ .Name }}---{{ .Namespace }}.public,internal+https://{{ .Name }}---{{ .Namespace }}.internal?always=true",
		}
		kpkg.Spec.Version = "1.0.0"
		kpkg.Spec.Config = map[string]v1alpha1.EnvVarValueOrFrom{
			"X": {Value: "x"},
		}

		d := &v1alpha1.DeployDeployment{
			Kind: "Deployment",
		}
		d.Spec.Replicas = ptr.Ptr(int32(1))
		kpkg.Spec.Deploy.SetUnderlying(d)

		kpkg.Spec.Containers = map[string]v1alpha1.Container{
			"web": {
				Image: v1alpha1.Image{
					Name: "docker.io/library/nginx",
					Tag:  "1.24.0-alpine",
					Platforms: []string{
						"linux/amd64",
						"linux/arm64",
					},
				},
				Ports: map[string]int32{
					"http": 80,
				},
			},
		}

		kpkg.Spec.Services = map[string]v1alpha1.Service{
			"#": {
				Ports: map[string]int32{
					"http": 80,
				},
				Paths: map[string]string{
					"http": "/",
				},
			},
		}

		kpkg.Spec.Volumes = map[string]v1alpha1.Volume{
			"html": {
				Underlying: &v1alpha1.VolumeConfigMap{
					Type: "ConfigMap",
					Spec: &v1alpha1.ConfigMapSpec{
						Data: map[string]string{
							"index.html": "<div>hello</div>",
						},
					},
					VolumeMount: v1alpha1.VolumeMount{
						MountPath: "/usr/share/nginx/html",
					},
				},
			},
		}

		list, err := SortedExtract(kpkg)
		testingx.Expect(t, err, testingx.Be[error](nil))

		data, _ := yaml.Marshal(list)
		fmt.Println(string(data))
	})
}
