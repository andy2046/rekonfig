package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

// removeAnnoteReferences iterates over a list of children and removes the annotation
// references from the child before updating it
func (h *Handler) removeAnnoteReferences(obj podController, children []Object) error {
	for _, child := range children {
		// Filter the existing annotation references
		annoteRefs := make(map[string]string)
		ownerRef := getAnnoteKeyFromName(obj.GetName())
		for annoteKey, v := range child.GetAnnotations() {
			if annoteKey != ownerRef {
				annoteRefs[annoteKey] = v
			}
		}

		// Compare the annotation references and update if changed
		if !reflect.DeepEqual(annoteRefs, child.GetAnnotations()) {
			h.recorder.Eventf(child, corev1.EventTypeNormal, "RemoveWatch", "Removing watch for %s %s",
				kindOf(child), child.GetName())
			child.SetAnnotations(annoteRefs)
			err := h.Update(context.TODO(), child)
			if err != nil {
				return fmt.Errorf("error updating child %s/%s: %v", child.GetNamespace(), child.GetName(), err)
			}
		}
	}
	return nil
}

// updateAnnoteReferences determines which children need to have their
// annotation references added/updated and which need to have their annotation references
// removed and then performs all updates
func (h *Handler) updateAnnoteReferences(owner podController, existing []Object, current []configObject) error {
	// Add an owner annotation reference to each child object
	errChan := make(chan error)
	for _, obj := range current {
		go func(child Object) {
			errChan <- h.updateAnnoteReference(owner, child)
		}(obj.object)
	}

	// Return any errors encountered while updating the child objects
	errs := []string{}
	for range current {
		err := <-errChan
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("error(s) encountered updating children: %s", strings.Join(errs, ", "))
	}

	// Get the orphaned children and remove their owner annotation references
	orphans := getOrphans(existing, current)
	err := h.removeAnnoteReferences(owner, orphans)
	if err != nil {
		return fmt.Errorf("error removing owner annotation references: %v", err)
	}

	return nil
}

// updateAnnoteReference ensures that the child object has an annotation reference
// pointing to the owner
func (h *Handler) updateAnnoteReference(owner podController, child Object) error {
	ownerRef, v := getAnnoteReference(owner)
	for ref := range child.GetAnnotations() {
		// owner annotation reference already exists, do nothing
		if ref == ownerRef {
			return nil
		}
	}

	// Append the new owner annotation reference and update the child
	h.recorder.Eventf(child, corev1.EventTypeNormal, "AddWatch", "Adding watch for %s %s",
		kindOf(child), child.GetName())
	annoteRefs := child.GetAnnotations()
	if annoteRefs == nil {
		annoteRefs = make(map[string]string)
	}
	annoteRefs[ownerRef] = v
	child.SetAnnotations(annoteRefs)
	err := h.Update(context.TODO(), child)
	if err != nil {
		return fmt.Errorf("error updating child: %v", err)
	}
	return nil
}

// getOrphans creates a slice of orphaned child objects that need their
// owner references removing
func getOrphans(existing []Object, current []configObject) []Object {
	orphans := []Object{}
	for _, child := range existing {
		if !isIn(current, child) {
			orphans = append(orphans, child)
		}
	}
	return orphans
}

// getAnnoteReference constructs an annotation reference key/value pointing to the object given
func getAnnoteReference(obj podController) (annoteKey, v string) {
	annoteKey, v = getAnnoteKeyFromName(obj.GetName()), kindOf(obj)
	return
}

// isIn checks whether a child object exists within a slice of objects
func isIn(list []configObject, child Object) bool {
	for _, obj := range list {
		if obj.object.GetUID() == child.GetUID() {
			return true
		}
	}
	return false
}

// kindOf returns the Kind of the given object as a string
func kindOf(obj Object) string {
	switch obj.(type) {
	case *corev1.ConfigMap:
		return "ConfigMap"
	case *corev1.Secret:
		return "Secret"
	case *deployment:
		return "Deployment"
	case *statefulset:
		return "StatefulSet"
	case *daemonset:
		return "DaemonSet"
	default:
		return "Unknown"
	}
}

// hasRequiredAnnotation returns true if the given PodController has
// the required annotation present
func hasRequiredAnnotation(obj podController) bool {
	annotations := obj.GetAnnotations()
	if value, ok := annotations[RequiredAnnotation]; ok && value == requiredAnnotationValue {
		return true
	}
	return false
}
