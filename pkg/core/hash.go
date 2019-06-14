package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
)

// calculateConfigHash returns a sha256 hash for the configuration within the children objects as a string
func calculateConfigHash(children []configObject) (string, error) {
	// hashSource contains all the data to be hashed
	hashSource := struct {
		ConfigMaps map[string]map[string]string `json:"configMaps"`
		Secrets    map[string]map[string][]byte `json:"secrets"`
	}{
		ConfigMaps: make(map[string]map[string]string),
		Secrets:    make(map[string]map[string][]byte),
	}

	// Add the data from each child to the hashSource
	// All children should be in the same namespace so child name should be unique
	for _, child := range children {
		if child.object != nil {
			switch child.object.(type) {
			case *corev1.ConfigMap:
				hashSource.ConfigMaps[child.object.GetName()] = getConfigMapData(child)
			case *corev1.Secret:
				hashSource.Secrets[child.object.GetName()] = getSecretData(child)
			default:
				return "", fmt.Errorf("passed unknown type: %v", reflect.TypeOf(child))
			}
		}
	}

	// Convert the hashSource to a byte slice so that it can be hashed
	hashSourceBytes, err := json.Marshal(hashSource)
	if err != nil {
		return "", fmt.Errorf("unable to marshal JSON: %v", err)
	}

	hashBytes := sha256.Sum256(hashSourceBytes)
	return hex.EncodeToString(hashBytes[:]), nil
}

// getConfigMapData extracts all the relevant data from the ConfigMap, whether it is
// the whole ConfigMap or only the specified keys.
func getConfigMapData(child configObject) map[string]string {
	cm := *child.object.(*corev1.ConfigMap)
	if child.allKeys {
		return cm.Data
	}
	keyData := make(map[string]string)
	for key := range child.keys {
		if value, exists := cm.Data[key]; exists {
			keyData[key] = value
		}
	}
	return keyData
}

// getSecretData extracts all the relevant data from the Secret, whether it is
// the whole Secret or only the specified keys.
func getSecretData(child configObject) map[string][]byte {
	s := *child.object.(*corev1.Secret)
	if child.allKeys {
		return s.Data
	}
	keyData := make(map[string][]byte)
	for key := range child.keys {
		if value, exists := s.Data[key]; exists {
			keyData[key] = value
		}
	}
	return keyData
}

// setConfigHash upates the configuration hash of the given Deployment to the given string
func setConfigHash(obj podController, hash string) {
	// Get the existing annotations
	podTemplate := obj.GetPodTemplate()
	annotations := podTemplate.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	// Update the annotations
	annotations[ConfigHashAnnotation] = hash
	podTemplate.SetAnnotations(annotations)
	obj.SetPodTemplate(podTemplate)
}