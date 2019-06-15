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

func TestFinalizer(t *testing.T) {
	deploymentObject := testutil.ExampleDeployment.DeepCopy()
	f := deploymentObject.GetFinalizers()
	f = append(f, "kubernetes")
	deploymentObject.SetFinalizers(f)
	podControllerDeployment := &deployment{deploymentObject}

	addFinalizer(podControllerDeployment)
	assert.Contains(t, deploymentObject.GetFinalizers(), FinalizerString, "should be added")
	assert.Contains(t, deploymentObject.GetFinalizers(), "kubernetes", "should leave existing finalizers in place")

	assert.True(t, hasFinalizer(podControllerDeployment), "should be there")

	removeFinalizer(podControllerDeployment)
	assert.NotContains(t, deploymentObject.GetFinalizers(), FinalizerString, "should be removed")
	assert.Contains(t, deploymentObject.GetFinalizers(), "kubernetes", "should leave existing finalizers in place")

	assert.False(t, hasFinalizer(podControllerDeployment), "should not be there")
}
