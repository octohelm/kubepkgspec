package v1alpha1

import (
	"bytes"
	"strings"

	"github.com/octohelm/x/ptr"
	corev1 "k8s.io/api/core/v1"
)

type Image struct {
	// 镜像名
	Name string `json:"name"`
	// 镜像标签
	Tag string `json:"tag,omitzero"`
	// 镜像摘要
	Digest string `json:"digest,omitzero"`
	// 镜像支持的平台
	Platforms []string `json:"platforms,omitzero"`
	// 镜像拉取策略
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitzero"`
}

func (v Image) FullName() string {
	s := strings.Builder{}
	s.WriteString(v.Name)
	if tag := v.Tag; tag != "" {
		s.WriteString(":")
		s.WriteString(tag)
	}
	if digest := v.Digest; digest != "" {
		s.WriteString("@")
		s.WriteString(digest)
	}
	return s.String()
}

type Container struct {
	// 镜像
	Image Image `json:"image"`
	// 运行目录
	WorkingDir string `json:"workingDir,omitzero"`
	// 命令
	Command []string `json:"command,omitzero"`
	// 参数
	Args []string `json:"args,omitzero"`
	// 环境变量
	Env map[string]EnvVarValueOrFrom `json:"env,omitzero"`
	// 暴露端口
	Ports map[string]int32 `json:"ports,omitzero"`

	Stdin     bool `json:"stdin,omitzero"`
	StdinOnce bool `json:"stdinOnce,omitzero"`
	TTY       bool `json:"tty,omitzero"`

	Resources      *corev1.ResourceRequirements `json:"resources,omitzero"`
	LivenessProbe  *corev1.Probe                `json:"livenessProbe,omitzero"`
	ReadinessProbe *corev1.Probe                `json:"readinessProbe,omitzero"`
	StartupProbe   *corev1.Probe                `json:"startupProbe,omitzero"`
	Lifecycle      *corev1.Lifecycle            `json:"lifecycle,omitzero"`

	SecurityContext          *corev1.SecurityContext         `json:"securityContext,omitzero"`
	TerminationMessagePath   string                          `json:"terminationMessagePath,omitzero"`
	TerminationMessagePolicy corev1.TerminationMessagePolicy `json:"terminationMessagePolicy,omitzero"`
}

type EnvVarValueOrFrom struct {
	Value     string
	ValueFrom *corev1.EnvVarSource
}

func (envVar EnvVarValueOrFrom) MarshalText() ([]byte, error) {
	if envVar.ValueFrom == nil {
		return []byte(envVar.Value), nil
	}

	buf := bytes.NewBufferString("@")

	if ref := envVar.ValueFrom.FieldRef; ref != nil {
		buf.WriteString("field/")
		buf.WriteString(ref.FieldPath)
	} else if ref := envVar.ValueFrom.ResourceFieldRef; ref != nil {
		buf.WriteString("resource/")
		buf.WriteString(ref.Resource)
	} else if ref := envVar.ValueFrom.ConfigMapKeyRef; ref != nil {
		buf.WriteString("configMap/")
		buf.WriteString(ref.Name)
		buf.WriteString("/")
		buf.WriteString(ref.Key)
		if ref.Optional != nil && *ref.Optional {
			buf.WriteString("?")
		}
	} else if ref := envVar.ValueFrom.SecretKeyRef; ref != nil {
		buf.WriteString("secret/")
		buf.WriteString(ref.Name)
		buf.WriteString("/")
		buf.WriteString(ref.Key)
		if ref.Optional != nil && *ref.Optional {
			buf.WriteString("?")
		}
	}

	return buf.Bytes(), nil
}

func (envVar *EnvVarValueOrFrom) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 {
		return nil
	}

	if text[0] == '@' {
		if idx := bytes.LastIndex(text, []byte("/")); idx > -1 {
			src := string(text[1:idx])
			key := string(text[idx+1:])

			switch {
			case strings.HasPrefix(src, "field"):
				envVar.ValueFrom = &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: key,
					},
				}
			case strings.HasPrefix(src, "resource"):
				envVar.ValueFrom = &corev1.EnvVarSource{
					ResourceFieldRef: &corev1.ResourceFieldSelector{
						Resource: key,
					},
				}
				return nil
			case strings.HasPrefix(src, "configMap"):
				r := &corev1.ConfigMapKeySelector{}
				r.Name = src[len("configMap/"):]
				r.Key = strings.TrimRight(key, "?")
				if strings.HasSuffix(key, "?") {
					r.Optional = ptr.Ptr(true)
				}
				envVar.ValueFrom = &corev1.EnvVarSource{
					ConfigMapKeyRef: r,
				}
				return nil
			case strings.HasPrefix(src, "secret"):
				r := &corev1.SecretKeySelector{}
				r.Name = src[len("secret/"):]
				r.Key = strings.TrimRight(key, "?")
				if strings.HasSuffix(key, "?") {
					r.Optional = ptr.Ptr(true)
				}
				envVar.ValueFrom = &corev1.EnvVarSource{
					SecretKeyRef: r,
				}
				return nil
			}
		}
	}

	envVar.Value = string(text)

	return nil
}
