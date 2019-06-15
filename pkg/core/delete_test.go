package core

import (
	"context"
	"testing"
	"time"

	"github.com/andy2046/rekonfig/pkg/testutil"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func init() {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))
}

func TestDelete_handleDelete(t *testing.T) {
	deploymentObject := testutil.ExampleDeployment.DeepCopy()
	deploymentObject.SetAnnotations(map[string]string{RequiredAnnotation: requiredAnnotationValue})
	f := deploymentObject.GetFinalizers()
	f = append(f, FinalizerString)
	f = append(f, "do.not.delete/finalizer")
	deploymentObject.SetFinalizers(f)

	podControllerDeployment := &deployment{deploymentObject}

	cm1 := testutil.ExampleConfigMap1.DeepCopy()
	cm2 := testutil.ExampleConfigMap2.DeepCopy()
	s1 := testutil.ExampleSecret1.DeepCopy()
	s2 := testutil.ExampleSecret2.DeepCopy()
	for _, obj := range []Object{cm1, cm2, s1, s2} {
		obj.SetAnnotations(map[string]string{
			getAnnoteKeyFromName(deploymentObject.GetName()): kindOf(podControllerDeployment),
		})
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		cm1,
		cm2,
		s1,
		s2,
		deploymentObject,
	}

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	s := scheme.Scheme
	recorder := record.NewBroadcasterForTests(5 * time.Second)
	h := NewHandler(cl, recorder.NewRecorder(s, corev1.EventSource{Component: "rekonfig"}))

	_, err := h.handleDelete(podControllerDeployment)
	assert.Nil(t, err, "handleDelete should return nil err")

	// Check if deployment finalizer has been removed.
	dep := &appsv1.Deployment{}
	cl.Get(context.TODO(), types.NamespacedName{
		Name:      deploymentObject.GetName(),
		Namespace: deploymentObject.GetNamespace(),
	}, dep)
	depFz := dep.GetFinalizers()
	assert.Len(t, depFz, 1, "deployment finalizer should be only one left")
	assert.Equal(t, "do.not.delete/finalizer", depFz[0], "deployment finalizer should be this value")

	// Check if configmap annotation has been removed.
	cm := &corev1.ConfigMap{}
	cl.Get(context.TODO(), types.NamespacedName{
		Name:      cm1.GetName(),
		Namespace: cm1.GetNamespace(),
	}, cm)
	assert.Len(t, cm.GetAnnotations(), 0, "configmap annotation should be empty")
}

func TestDelete_toBeDeleted(t *testing.T) {
	deploymentObject := testutil.ExampleDeployment.DeepCopy()
	assert.False(t, toBeDeleted(deploymentObject), "toBeDeleted should return false if the deleteion timestamp is nil")

	tm := metav1.NewTime(time.Now())
	deploymentObject.SetDeletionTimestamp(&tm)
	assert.True(t, toBeDeleted(deploymentObject), "toBeDeleted should return true if deletion timestamp is non-nil")
}
