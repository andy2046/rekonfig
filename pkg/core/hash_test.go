package core

import (
	"testing"

	"github.com/andy2046/rekonfig/pkg/testutil"
	"github.com/stretchr/testify/assert"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func init() {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))
}

func TestHash_setConfigHash(t *testing.T) {
	deploymentObject := testutil.ExampleDeployment.DeepCopy()
	podControllerDeployment := &deployment{deploymentObject}

	podAnnotations := deploymentObject.Spec.Template.GetAnnotations()
	if podAnnotations == nil {
		podAnnotations = make(map[string]string)
	}
	podAnnotations["existing"] = "annotation"
	deploymentObject.Spec.Template.SetAnnotations(podAnnotations)

	setConfigHash(podControllerDeployment, "1234")

	podAnnotations = deploymentObject.Spec.Template.GetAnnotations()
	assert.NotNil(t, podAnnotations, "deployment annotation should not be nil")

	hash, ok := podAnnotations[ConfigHashAnnotation]
	assert.True(t, ok, "deployment annotation should be set by setConfigHash")
	assert.Equal(t, "1234", hash, "deployment annotation should be set by setConfigHash")

	hash, ok = podAnnotations["existing"]
	assert.True(t, ok, "existing deployment annotation should still be there")
	assert.Equal(t, "annotation", hash, "existing deployment annotation should still be there")
}

func TestHash_calculateConfigHash(t *testing.T) {
	var modified = "modified"
	cm1 := testutil.ExampleConfigMap1.DeepCopy()
	cm2 := testutil.ExampleConfigMap2.DeepCopy()
	cm3 := testutil.ExampleConfigMap3.DeepCopy()
	s1 := testutil.ExampleSecret1.DeepCopy()
	s2 := testutil.ExampleSecret2.DeepCopy()
	s3 := testutil.ExampleSecret3.DeepCopy()

	// different hash when an allKeys child's data is updated
	c := []configObject{
		{object: cm1, allKeys: true},
		{object: cm2, allKeys: true},
		{object: s1, allKeys: true},
		{object: s2, allKeys: true},
	}
	h1, err := calculateConfigHash(c)
	assert.Nil(t, err, "no error")
	cm1.Data["key1"] = modified
	h2, err := calculateConfigHash(c)
	assert.Nil(t, err, "no error")
	assert.NotEqual(t, h2, h1, "should not be equal")

	// different hash when a single-field child's data is updated
	c = []configObject{
		{object: cm1, allKeys: false, keys: map[string]struct{}{
			"key1": {},
		},
		},
		{object: cm2, allKeys: true},
		{object: s1, allKeys: false, keys: map[string]struct{}{
			"key1": {},
		},
		},
		{object: s2, allKeys: true},
	}
	h1, err = calculateConfigHash(c)
	assert.Nil(t, err, "no error")
	cm1.Data["key1"] = modified
	s1.Data["key1"] = []byte("modified")
	h2, err = calculateConfigHash(c)
	assert.Nil(t, err, "no error")
	assert.NotEqual(t, h2, h1, "should not be equal")

	// same hash when a single-field child's data is updated but not for that field
	c = []configObject{
		{object: cm1, allKeys: false, keys: map[string]struct{}{
			"key1": {},
		},
		},
		{object: cm2, allKeys: true},
		{object: s1, allKeys: false, keys: map[string]struct{}{
			"key1": {},
		},
		},
		{object: s2, allKeys: true},
	}
	h1, err = calculateConfigHash(c)
	assert.Nil(t, err, "no error")
	cm1.Data["key3"] = modified
	s1.Data["key3"] = []byte("modified")
	h2, err = calculateConfigHash(c)
	assert.Nil(t, err, "no error")
	assert.Equal(t, h2, h1, "should be equal")

	// same hash when a child's metadata is updated
	c = []configObject{
		{object: cm1, allKeys: true},
		{object: cm2, allKeys: true},
		{object: s1, allKeys: true},
		{object: s2, allKeys: true},
	}
	h1, err = calculateConfigHash(c)
	assert.Nil(t, err, "no error")
	s1.Annotations = map[string]string{"new": "annotations"}
	h2, err = calculateConfigHash(c)
	assert.Nil(t, err, "no error")
	assert.Equal(t, h2, h1, "should be equal")

	// same hash independent of child ordering
	c1 := []configObject{
		{object: cm1, allKeys: true},
		{object: cm2, allKeys: true},
		{object: cm3, allKeys: false, keys: map[string]struct{}{
			"key1": {},
			"key2": {},
		},
		},
		{object: s1, allKeys: true},
		{object: s2, allKeys: true},
		{object: s3, allKeys: false, keys: map[string]struct{}{
			"key1": {},
			"key2": {},
		},
		},
	}
	c2 := []configObject{
		{object: cm1, allKeys: true},
		{object: s2, allKeys: true},
		{object: s3, allKeys: false, keys: map[string]struct{}{
			"key1": {},
			"key2": {},
		},
		},
		{object: cm2, allKeys: true},
		{object: s1, allKeys: true},
		{object: cm3, allKeys: false, keys: map[string]struct{}{
			"key2": {},
			"key1": {},
		},
		},
	}

	h1, err = calculateConfigHash(c1)
	assert.Nil(t, err, "no error")
	h2, err = calculateConfigHash(c2)
	assert.Nil(t, err, "no error")
	assert.Equal(t, h2, h1, "should be equal")
}
