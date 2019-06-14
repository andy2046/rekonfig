package core

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	// ConfigHashAnnotation is the annotation key in PodTemplate for config hash
	ConfigHashAnnotation = "rekonfig.gitops.in/konfig-hash"

	// FinalizerString is the finalizer added to deployments to allow rekonfig to
	// perform advanced deletion logic
	FinalizerString = "rekonfig.gitops.in/finalizer"

	// RequiredAnnotation is the annotation key in the Deployment that
	// we check for before processing the deployment
	RequiredAnnotation = "rekonfig.gitops.in/update-on-konfig-change"

	// requiredAnnotationValue is the annotation value in the Deployment that
	// we check for before processing the deployment
	requiredAnnotationValue = "true"

	// configAnnotationPrefix is the prefix of annotation key in ConfigMap/Secret for Deployment owner
	configAnnotationPrefix = "rekonfig.gitops.in/owner-"
)

type (
	// Object is a helper interface when passing Kubernetes resources between methods.
	// All Kubernetes resources should implement both of the interfaces
	Object interface {
		runtime.Object
		metav1.Object
	}

	// configObject is a container of an "Object" along with metadata
	// that we use to determine what to use from that Object
	configObject struct {
		object   Object
		required bool
		allKeys  bool
		keys     map[string]struct{}
	}

	podController interface {
		runtime.Object
		metav1.Object
		GetObject() runtime.Object
		GetPodTemplate() *corev1.PodTemplateSpec
		SetPodTemplate(*corev1.PodTemplateSpec)
		DeepCopy() podController
	}

	deployment struct {
		*appsv1.Deployment
	}
)

func (d *deployment) GetObject() runtime.Object {
	return d.Deployment
}

func (d *deployment) GetPodTemplate() *corev1.PodTemplateSpec {
	return &d.Deployment.Spec.Template
}

func (d *deployment) SetPodTemplate(template *corev1.PodTemplateSpec) {
	d.Deployment.Spec.Template = *template
}

func (d *deployment) DeepCopy() podController {
	return &deployment{d.Deployment.DeepCopy()}
}

type statefulset struct {
	*appsv1.StatefulSet
}

func (d *statefulset) GetObject() runtime.Object {
	return d.StatefulSet
}

func (d *statefulset) GetPodTemplate() *corev1.PodTemplateSpec {
	return &d.StatefulSet.Spec.Template
}

func (d *statefulset) SetPodTemplate(template *corev1.PodTemplateSpec) {
	d.StatefulSet.Spec.Template = *template
}

func (d *statefulset) DeepCopy() podController {
	return &statefulset{d.StatefulSet.DeepCopy()}
}

type daemonset struct {
	*appsv1.DaemonSet
}

func (d *daemonset) GetObject() runtime.Object {
	return d.DaemonSet
}

func (d *daemonset) GetPodTemplate() *corev1.PodTemplateSpec {
	return &d.DaemonSet.Spec.Template
}

func (d *daemonset) SetPodTemplate(template *corev1.PodTemplateSpec) {
	d.DaemonSet.Spec.Template = *template
}

func (d *daemonset) DeepCopy() podController {
	return &daemonset{d.DaemonSet.DeepCopy()}
}

// GetNameFromAnnoteKey extracts name from Config annotation key
func GetNameFromAnnoteKey(annoteKey string) string {
	if strings.HasPrefix(annoteKey, configAnnotationPrefix) {
		return strings.TrimSpace(annoteKey[len(configAnnotationPrefix):])
	}
	return ""
}

// getAnnoteKeyFromName constructs Config annotation key with given name
func getAnnoteKeyFromName(name string) string {
	return fmt.Sprint(configAnnotationPrefix, name)
}
