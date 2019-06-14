package core

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// configMetadata contains information about ConfigMaps/Secrets referenced within PodTemplates
type configMetadata struct {
	required bool
	allKeys  bool
	keys     map[string]struct{}
}

// getResult is a helper struct returned from the getObject method
type getResult struct {
	err      error
	obj      Object
	metadata configMetadata
}

// getCurrentChildren returns a list of all Secrets and ConfigMaps that are
// referenced in the Deployment's spec.
// Any reference to a whole ConfigMap or Secret (i.e. via an EnvFrom or a Volume) will
// result in one entry in the list,
// irrespective of whether individual elements are also references (i.e. via an Env entry)
func (h *Handler) getCurrentChildren(obj podController) ([]configObject, error) {
	configMaps, secrets := getChildNamesByType(obj)

	// get all the ConfigMaps and Secrets from API server
	resultsChan := make(chan getResult)
	for name, metadata := range configMaps {
		go func(name string, metadata configMetadata) {
			resultsChan <- h.getConfigMap(obj.GetNamespace(), name, metadata)
		}(name, metadata)
	}
	for name, metadata := range secrets {
		go func(name string, metadata configMetadata) {
			resultsChan <- h.getSecret(obj.GetNamespace(), name, metadata)
		}(name, metadata)
	}

	// collect results from the API server gets
	var (
		errs     []string
		children []configObject
	)
	for i := 0; i < len(configMaps)+len(secrets); i++ {
		result := <-resultsChan
		if result.err != nil {
			errs = append(errs, result.err.Error())
		}
		if result.obj != nil {
			children = append(children, configObject{
				object:   result.obj,
				required: result.metadata.required,
				allKeys:  result.metadata.allKeys,
				keys:     result.metadata.keys,
			})
		}
	}

	// If there was any error, don't return any children
	if len(errs) > 0 {
		return []configObject{}, fmt.Errorf("error(s) encountered when geting children: %s", strings.Join(errs, ", "))
	}

	// No error, return the list of children
	return children, nil
}

// getChildNamesByType parses the Deployment object and returns two maps,
// the first containing ConfigMap metadata for all referenced ConfigMaps, key is the name of the ConfigMap,
// the second containing Secret metadata for all referenced Secrets, key is the name of the Secrets
func getChildNamesByType(obj podController) (map[string]configMetadata, map[string]configMetadata) {
	configMaps := make(map[string]configMetadata)
	secrets := make(map[string]configMetadata)

	// Loop all Volumes and check the VolumeSources for ConfigMaps and Secrets
	for _, vol := range obj.GetPodTemplate().Spec.Volumes {
		if cm := vol.VolumeSource.ConfigMap; cm != nil {
			configMaps[cm.Name] = configMetadata{required: true, allKeys: true}
		}
		if s := vol.VolumeSource.Secret; s != nil {
			secrets[s.SecretName] = configMetadata{required: true, allKeys: true}
		}
	}

	// Loop all Containers and the respective EnvFrom,
	// check the EnvFromSources for ConfigMaps and Secrets
	for _, container := range obj.GetPodTemplate().Spec.Containers {
		for _, env := range container.EnvFrom {
			if cm := env.ConfigMapRef; cm != nil {
				configMaps[cm.Name] = configMetadata{required: true, allKeys: true}
			}
			if s := env.SecretRef; s != nil {
				secrets[s.Name] = configMetadata{required: true, allKeys: true}
			}
		}
	}

	// Loop all Containers and the respective Env
	for _, container := range obj.GetPodTemplate().Spec.Containers {
		for _, env := range container.Env {
			if valFrom := env.ValueFrom; valFrom != nil {
				if cm := valFrom.ConfigMapKeyRef; cm != nil {
					configMaps[cm.Name] = parseConfigMapKeyRef(configMaps[cm.Name], cm)
				}
				if s := valFrom.SecretKeyRef; s != nil {
					secrets[s.Name] = parseSecretKeyRef(secrets[s.Name], s)
				}
			}
		}
	}

	return configMaps, secrets
}

// parseConfigMapKeyRef updates the metadata for a ConfigMap
// to include the keys specified in the ConfigMapKeySelector
func parseConfigMapKeyRef(metadata configMetadata, cm *corev1.ConfigMapKeySelector) configMetadata {
	if !metadata.allKeys {
		if metadata.keys == nil {
			metadata.keys = make(map[string]struct{})
		}
		if cm.Optional == nil || !*cm.Optional {
			metadata.required = true
		}
		metadata.keys[cm.Key] = struct{}{}
	}
	return metadata
}

// parseSecretKeyRef updates the metadata for a Secret
// to include the keys specified in the SecretKeySelector
func parseSecretKeyRef(metadata configMetadata, s *corev1.SecretKeySelector) configMetadata {
	if !metadata.allKeys {
		if metadata.keys == nil {
			metadata.keys = make(map[string]struct{})
		}
		if s.Optional == nil || !*s.Optional {
			metadata.required = true
		}
		metadata.keys[s.Key] = struct{}{}
	}
	return metadata
}

// getConfigMap gets a ConfigMap with the given name and namespace from API server
func (h *Handler) getConfigMap(namespace, name string, metadata configMetadata) getResult {
	return h.getObject(namespace, name, metadata, &corev1.ConfigMap{})
}

// getSecret gets a Secret with the given name and namespace from API server
func (h *Handler) getSecret(namespace, name string, metadata configMetadata) getResult {
	return h.getObject(namespace, name, metadata, &corev1.Secret{})
}

// getObject gets the Object with the given name and namespace from API server
func (h *Handler) getObject(namespace, name string, metadata configMetadata, obj Object) getResult {
	objectName := types.NamespacedName{Namespace: namespace, Name: name}
	err := h.Get(context.TODO(), objectName, obj)
	if err != nil {
		if metadata.required {
			return getResult{err: err}
		}
		return getResult{metadata: metadata}
	}
	return getResult{obj: obj, metadata: metadata}
}

// getExistingChildren returns a list of Secrets and ConfigMaps owned by the Deployment instance
func (h *Handler) getExistingChildren(obj podController) ([]Object, error) {
	inNamespace := client.InNamespace(obj.GetNamespace())

	// List all ConfigMaps in the Deployment's namespace
	configMaps := &corev1.ConfigMapList{}
	if err := h.List(context.TODO(), inNamespace, configMaps); err != nil {
		return []Object{}, fmt.Errorf("error listing ConfigMaps: %v", err)
	}

	// List all Secrets in the Deployment's namespcae
	secrets := &corev1.SecretList{}
	if err := h.List(context.TODO(), inNamespace, secrets); err != nil {
		return []Object{}, fmt.Errorf("error listing Secrets: %v", err)
	}

	// Loop all ConfigMaps/Secrets and add to the output children list if it's owned by the Deployment
	children := []Object{}
	for _, cm := range configMaps.Items {
		if isOwnedBy(&cm, obj) {
			children = append(children, cm.DeepCopy())
		}
	}
	for _, s := range secrets.Items {
		if isOwnedBy(&s, obj) {
			children = append(children, s.DeepCopy())
		}
	}

	return children, nil
}

// isOwnedBy returns true if the child has an annotation reference pointing to the owner object
func isOwnedBy(child, owner metav1.Object) bool {
	ownerRef := getAnnoteKeyFromName(owner.GetName())
	for annoteKey := range child.GetAnnotations() {
		if annoteKey == ownerRef {
			return true
		}
	}
	return false
}
