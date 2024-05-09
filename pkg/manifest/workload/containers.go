package workload

import (
	"iter"

	"github.com/octohelm/kubepkgspec/pkg/manifest/object"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

func Containers(obj iter.Seq[object.Object]) iter.Seq[*corev1.Container] {
	workloads := Workloads(obj)

	return func(yield func(c *corev1.Container) bool) {
		for o := range workloads {
			switch x := o.(type) {
			case *appsv1.DaemonSet:
				for i, c := range x.Spec.Template.Spec.Containers {
					if !yield(&c) {
						return
					}
					x.Spec.Template.Spec.Containers[i] = c
				}
			case *appsv1.Deployment:
				for i, c := range x.Spec.Template.Spec.Containers {
					if !yield(&c) {
						return
					}
					x.Spec.Template.Spec.Containers[i] = c
				}
			case *appsv1.StatefulSet:
				for i, c := range x.Spec.Template.Spec.Containers {
					if !yield(&c) {
						return
					}
					x.Spec.Template.Spec.Containers[i] = c
				}
			case *batchv1.Job:
				for i, c := range x.Spec.Template.Spec.Containers {
					if !yield(&c) {
						return
					}
					x.Spec.Template.Spec.Containers[i] = c
				}
			case *batchv1.CronJob:
				for i, c := range x.Spec.JobTemplate.Spec.Template.Spec.Containers {
					if !yield(&c) {
						return
					}
					x.Spec.JobTemplate.Spec.Template.Spec.Containers[i] = c
				}
			}
		}
	}
}
