package daemonset

import (
	"context"
	"testing"
	"time"

	"github.com/andy2046/rekonfig/pkg/core"
	"github.com/andy2046/rekonfig/pkg/testutil"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var configAnnotationPrefix = "rekonfig.gitops.in/owner-"

func init() {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))
}

// TestDaemonSetController runs Reconcile() against a fake client.
func TestDaemonSetController(t *testing.T) {
	cm1 := testutil.ExampleConfigMap1.DeepCopy()
	cm2 := testutil.ExampleConfigMap2.DeepCopy()
	cm3 := testutil.ExampleConfigMap3.DeepCopy()
	s1 := testutil.ExampleSecret1.DeepCopy()
	s2 := testutil.ExampleSecret2.DeepCopy()
	s3 := testutil.ExampleSecret3.DeepCopy()

	deploymentObject := testutil.ExampleDaemonSet.DeepCopy()
	// add the required annotation
	annotations := deploymentObject.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations[core.RequiredAnnotation] = "true"
	deploymentObject.SetAnnotations(annotations)

	// Objects to track in the fake client.
	objs := []runtime.Object{
		cm1,
		cm2,
		cm3,
		s1,
		s2,
		s3,
		deploymentObject,
	}

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	s := scheme.Scheme
	recorder := record.NewBroadcasterForTests(5 * time.Second)
	h := core.NewHandler(cl, recorder.NewRecorder(s, corev1.EventSource{Component: "rekonfig"}))

	// Create a ReconcileDaemonSet object with the scheme and fake client.
	r := &ReconcileDaemonSet{handler: h, scheme: s}
	key := types.NamespacedName{Namespace: deploymentObject.GetNamespace(), Name: deploymentObject.GetName()}
	// Mock request to simulate Reconcile() being called on an event for a watched resource.
	req := reconcile.Request{
		NamespacedName: key,
	}
	res, err := r.Reconcile(req)
	assert.Nil(t, err, "no reconcile error")
	assert.False(t, res.Requeue, "reconcile should not requeue request")

	cm := &corev1.ConfigMap{}
	cl.Get(context.TODO(), types.NamespacedName{
		Name:      cm1.GetName(),
		Namespace: cm1.GetNamespace(),
	}, cm)
	cm1 = cm
	cm = &corev1.ConfigMap{}
	cl.Get(context.TODO(), types.NamespacedName{
		Name:      cm2.GetName(),
		Namespace: cm2.GetNamespace(),
	}, cm)
	cm2 = cm
	cm = &corev1.ConfigMap{}
	cl.Get(context.TODO(), types.NamespacedName{
		Name:      cm3.GetName(),
		Namespace: cm3.GetNamespace(),
	}, cm)
	cm3 = cm

	sc := &corev1.Secret{}
	cl.Get(context.TODO(), types.NamespacedName{
		Name:      s1.GetName(),
		Namespace: s1.GetNamespace(),
	}, sc)
	s1 = sc
	sc = &corev1.Secret{}
	cl.Get(context.TODO(), types.NamespacedName{
		Name:      s2.GetName(),
		Namespace: s2.GetNamespace(),
	}, sc)
	s2 = sc
	sc = &corev1.Secret{}
	cl.Get(context.TODO(), types.NamespacedName{
		Name:      s3.GetName(),
		Namespace: s3.GetNamespace(),
	}, sc)
	s3 = sc

	// Check if configmap annotation has been added.
	ownerRef := configAnnotationPrefix + deploymentObject.GetName()
	for _, obj := range []core.Object{cm1, cm2, cm3, s1, s2, s3} {
		annotes, exists := obj.GetAnnotations(), false
		for annoteKey := range annotes {
			if annoteKey == ownerRef {
				exists = true
			}
		}
		assert.True(t, exists, "should add owner annotation reference to all children")
	}

	// Check if deployment finalizer has been added.
	dep := &appsv1.DaemonSet{}
	cl.Get(context.TODO(), types.NamespacedName{
		Name:      deploymentObject.GetName(),
		Namespace: deploymentObject.GetNamespace(),
	}, dep)
	depFz := dep.GetFinalizers()
	assert.Len(t, depFz, 1, "deployment finalizer should be added")
	assert.Equal(t, core.FinalizerString, depFz[0], "deployment finalizer should be this value")

	// Check if a config hash has been added to the Pod Template
	podAnnotes, exists, configHash := dep.Spec.Template.GetAnnotations(), false, ""
	for annoteKey, v := range podAnnotes {
		if annoteKey == core.ConfigHashAnnotation {
			exists = true
			configHash = v
		}
	}
	assert.True(t, exists, "should add a config hash to the Pod Template")

	// Check if deployment config hash updated when a child is updated
	cm1.Data["key1"] = "modified"
	err = cl.Update(context.TODO(), cm1)
	assert.Nil(t, err, "no error")

	res, err = r.Reconcile(req)
	assert.Nil(t, err, "no reconcile error")
	assert.False(t, res.Requeue, "reconcile should not requeue request")

	dep = &appsv1.DaemonSet{}
	cl.Get(context.TODO(), types.NamespacedName{
		Name:      deploymentObject.GetName(),
		Namespace: deploymentObject.GetNamespace(),
	}, dep)
	podAnnotes, exists = dep.Spec.Template.GetAnnotations(), false
	for annoteKey, v := range podAnnotes {
		if annoteKey == core.ConfigHashAnnotation {
			exists = true
			assert.NotEqual(t, configHash, v, "should update the config hash in the Pod Template")
		}
	}
	assert.True(t, exists, "should update the config hash in the Pod Template")
}
