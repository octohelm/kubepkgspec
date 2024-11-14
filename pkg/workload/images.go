package workload

import (
	"iter"
	"strings"

	kubepkgv1alpha1 "github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/kubepkgspec/pkg/object"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
)

func ParseImage(imageName string) *kubepkgv1alpha1.Image {
	nameDigest := strings.SplitN(imageName, "@", 2)

	nameTags := strings.SplitN(nameDigest[0], ":", 2)

	img := &kubepkgv1alpha1.Image{
		Name: nameTags[0],
	}

	if len(nameTags) == 2 {
		img.Tag = nameTags[1]
	}

	if len(nameDigest) == 2 {
		img.Digest = nameDigest[1]
	}

	return img
}

func Images(obj iter.Seq[object.Object]) iter.Seq[*kubepkgv1alpha1.Image] {
	return func(yield func(img *kubepkgv1alpha1.Image) bool) {
		for o := range obj {
			switch x := o.(type) {
			case *kubepkgv1alpha1.KubePkg:
				for name, c := range x.Spec.Containers {
					if !yield(&c.Image) {
						return
					}
					x.Spec.Containers[name] = c
				}

				for name, sub := range x.Spec.Manifests {
					if o, ok := object.Is(sub); ok {

						o, err := object.Convert(o)
						if err == nil {
							subImages := Images(func(yield func(object.Object) bool) {
								if !yield(o) {
									return
								}
							})

							for i := range subImages {
								if !yield(i) {
									return
								}
							}

							x.Spec.Manifests[name] = o

							continue
						}
					}

					x.Spec.Manifests[name] = sub
				}

				for name, v := range x.Spec.Volumes {
					if vi, ok := v.Underlying.(*kubepkgv1alpha1.VolumeImage); ok {
						if vi.Opt != nil && vi.Opt.Reference != "" {
							img := ParseImage(vi.Opt.Reference)
							if !yield(img) {
								return
							}
							vi.Opt.Reference = img.FullName()
							x.Spec.Volumes[name] = v
						}
					}
				}

			case *appsv1.DaemonSet:
				for i, c := range x.Spec.Template.Spec.Containers {
					img := ParseImage(c.Image)
					if !yield(img) {
						return
					}
					c.Image = img.FullName()
					x.Spec.Template.Spec.Containers[i] = c
				}
			case *appsv1.Deployment:
				for i, c := range x.Spec.Template.Spec.Containers {
					img := ParseImage(c.Image)
					if !yield(img) {
						return
					}
					c.Image = img.FullName()
					x.Spec.Template.Spec.Containers[i] = c
				}
			case *appsv1.StatefulSet:
				for i, c := range x.Spec.Template.Spec.Containers {
					img := ParseImage(c.Image)
					if !yield(img) {
						return
					}
					c.Image = img.FullName()
					x.Spec.Template.Spec.Containers[i] = c
				}
			case *batchv1.Job:
				for i, c := range x.Spec.Template.Spec.Containers {
					img := ParseImage(c.Image)
					if !yield(img) {
						return
					}
					c.Image = img.FullName()
					x.Spec.Template.Spec.Containers[i] = c
				}
			case *batchv1.CronJob:
				for i, c := range x.Spec.JobTemplate.Spec.Template.Spec.Containers {
					img := ParseImage(c.Image)
					if !yield(img) {
						return
					}
					c.Image = img.FullName()
					x.Spec.JobTemplate.Spec.Template.Spec.Containers[i] = c
				}
			}
		}
	}
}
