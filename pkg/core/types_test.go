package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func init() {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))
}

func TestTypes_GetNameFromAnnoteKey(t *testing.T) {
	owner := "myowner"
	annoteKey := configAnnotationPrefix + owner
	name := GetNameFromAnnoteKey(annoteKey)
	assert.Equal(t, owner, name, "should be equal")

	annoteKey = "non-sense-" + configAnnotationPrefix + owner
	name = GetNameFromAnnoteKey(annoteKey)
	assert.Equal(t, "", name, "should be equal")
}

func TestTypes_getAnnoteKeyFromName(t *testing.T) {
	name := "mydeploy"
	annoteKey := getAnnoteKeyFromName(name)
	assert.Equal(t, configAnnotationPrefix+name, annoteKey, "should be equal")
}
