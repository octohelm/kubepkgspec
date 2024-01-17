package kubepkg

import (
	"os"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestCustomResourceDefinition(t *testing.T) {
	data, _ := yaml.Marshal(CRDs)
	_ = os.WriteFile("crd.yaml", data, 0644)
}
