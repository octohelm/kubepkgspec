package workload

import (
	"iter"
	"reflect"

	batchv1 "k8s.io/api/batch/v1"

	"github.com/octohelm/kubepkgspec/pkg/object"
	appsv1 "k8s.io/api/apps/v1"
)

var (
	KindDeployment  = reflect.TypeOf(appsv1.Deployment{}).Name()
	KindStatefulSet = reflect.TypeOf(appsv1.StatefulSet{}).Name()
	KindDaemonSet   = reflect.TypeOf(appsv1.DaemonSet{}).Name()

	KindJob     = reflect.TypeOf(batchv1.Job{}).Name()
	KindCronJob = reflect.TypeOf(batchv1.CronJob{}).Name()
)

func IsWorkload(o object.Object) bool {
	gvk := o.GetObjectKind().GroupVersionKind()

	switch gvk.Kind {
	case KindDeployment, KindStatefulSet, KindDaemonSet:
		return gvk.Group == appsv1.GroupName
	case KindJob, KindCronJob:
		return gvk.Group == batchv1.GroupName
	}

	return false
}

func Workloads(obj iter.Seq[object.Object]) iter.Seq[object.Object] {
	return func(yield func(object.Object) bool) {
		for o := range obj {
			if !IsWorkload(o) {
				continue
			}

			if !yield(o) {
				return
			}
		}
	}
}
