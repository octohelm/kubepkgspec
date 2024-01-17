package manifest

import (
	"io"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/stretchr/objx"
	corev1 "k8s.io/api/core/v1"
)

func SubResourceName(kpkg *v1alpha1.KubePkg, name string) string {
	if name == "#" || name == "" {
		return kpkg.Name
	}
	if name[0] == '~' || name[0] == '/' {
		return name[1:]
	}
	return strings.Join([]string{kpkg.Name, name}, "-")
}

func PortProtocol(n string) corev1.Protocol {
	if strings.HasPrefix(n, "udp-") {
		return corev1.ProtocolUDP
	} else if strings.HasPrefix(n, "sctp-") {
		return corev1.ProtocolSCTP
	} else {
		return corev1.ProtocolTCP
	}
}

func PortNameAndHostPort(n string) (string, int32) {
	// <portName>:80
	if i := strings.Index(n, ":"); i > 0 {
		hostPort, _ := strconv.ParseInt(n[i+1:], 10, 32)
		return n[0:i], int32(hostPort)
	}
	return n, 0
}

func Intersection[E comparable](a []E, b []E) (c []E) {

	includes := map[E]bool{}
	for i := range a {
		includes[a[i]] = true
	}

	c = make([]E, 0, len(a)+len(b))
	for i := range b {
		x := b[i]

		if _, ok := includes[x]; ok {
			c = append(c, x)
		}
	}

	return
}

func Merge[X any](from *X, overwrites *X) (*X, error) {
	var srcObj objx.Map
	if err := decodeTo(from, &srcObj); err != nil {
		return nil, err
	}
	var overwritesObj objx.Map
	if err := decodeTo(overwrites, &overwritesObj); err != nil {
		return nil, err
	}

	merged := MergeObj(srcObj, overwritesObj)
	m := new(X)
	if err := decodeTo(merged, m); err != nil {
		return nil, err
	}
	return m, nil
}

func MergeObj(from objx.Map, patch objx.Map) objx.Map {
	if patch == nil {
		return from
	}

	mergedKeys := map[string]bool{}

	merged := from.Transform(func(key string, currValue any) (string, any) {
		mergedKeys[key] = true

		if patchValue, ok := patch[key]; ok {
			if m, ok := patchValue.(objx.Map); ok {
				patchValue = map[string]any(m)
			}

			switch p := patchValue.(type) {
			case map[string]any:
				switch x := currValue.(type) {
				case objx.Map:
					return key, MergeObj(x, p)
				case map[string]any:
					return key, MergeObj(x, p)
				}
			}

			// don't merge nil value
			if patchValue != nil {
				return key, patchValue
			}
		}

		return key, currValue
	})

	for key := range patch {
		if _, ok := mergedKeys[key]; !ok {
			merged[key] = patch[key]
		}
	}

	return merged
}

func decodeTo(from any, to any) error {
	r, w := io.Pipe()
	e := json.NewEncoder(w)
	go func() {
		defer w.Close()
		_ = e.Encode(from)
	}()
	d := json.NewDecoder(r)
	if err := d.Decode(&to); err != nil {
		return err
	}
	return nil
}
