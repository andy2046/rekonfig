package core

import (
	"testing"
	"time"

	"github.com/andy2046/rekonfig/pkg/testutil"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func init() {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))
}

func TestRef_removeAnnoteReferences(t *testing.T) {
	cm1 := testutil.ExampleConfigMap1.DeepCopy()
	cm2 := testutil.ExampleConfigMap2.DeepCopy()
	cm3 := testutil.ExampleConfigMap3.DeepCopy()
	cm4 := testutil.ExampleConfigMap4.DeepCopy()
	s1 := testutil.ExampleSecret1.DeepCopy()
	s2 := testutil.ExampleSecret2.DeepCopy()
	s3 := testutil.ExampleSecret3.DeepCopy()
	s4 := testutil.ExampleSecret4.DeepCopy()

	deploymentObject := testutil.ExampleDeployment.DeepCopy()
	podControllerDeployment := &deployment{deploymentObject}

	ownerRef, v := getAnnoteReference(podControllerDeployment)
	for _, obj := range []Object{cm1, cm2, s1, s2} {
		obj.SetAnnotations(map[string]string{
			ownerRef: v,
		})
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		cm1,
		cm2,
		cm3,
		cm4,
		s1,
		s2,
		s3,
		s4,
		deploymentObject,
	}

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	s := scheme.Scheme
	recorder := record.NewBroadcasterForTests(5 * time.Second)
	h := NewHandler(cl, recorder.NewRecorder(s, corev1.EventSource{Component: "rekonfig"}))

	children := []Object{cm1, s1}
	err := h.removeAnnoteReferences(podControllerDeployment, children)
	assert.Nil(t, err, "no error")
}

func TestRef_updateAnnoteReferences(t *testing.T) {
	cm1 := testutil.ExampleConfigMap1.DeepCopy()
	cm2 := testutil.ExampleConfigMap2.DeepCopy()
	cm3 := testutil.ExampleConfigMap3.DeepCopy()
	cm4 := testutil.ExampleConfigMap4.DeepCopy()
	s1 := testutil.ExampleSecret1.DeepCopy()
	s2 := testutil.ExampleSecret2.DeepCopy()
	s3 := testutil.ExampleSecret3.DeepCopy()
	s4 := testutil.ExampleSecret4.DeepCopy()

	deploymentObject := testutil.ExampleDeployment.DeepCopy()
	podControllerDeployment := &deployment{deploymentObject}

	ownerRef, v := getAnnoteReference(podControllerDeployment)
	for _, obj := range []Object{cm2, s1, s2} {
		obj.SetAnnotations(map[string]string{
			ownerRef: v,
		})
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{
		cm1,
		cm2,
		cm3,
		cm4,
		s1,
		s2,
		s3,
		s4,
		deploymentObject,
	}

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	s := scheme.Scheme
	recorder := record.NewBroadcasterForTests(5 * time.Second)
	h := NewHandler(cl, recorder.NewRecorder(s, corev1.EventSource{Component: "rekonfig"}))

	existing := []Object{cm2, cm3, s1, s2}
	current := []configObject{
		{object: cm1, allKeys: true},
		{object: s1, allKeys: true},
		{object: s3, allKeys: false, keys: map[string]struct{}{
			"key1": {},
			"key2": {},
		},
		},
	}
	err := h.updateAnnoteReferences(podControllerDeployment, existing, current)
	assert.Nil(t, err, "no error")
}
