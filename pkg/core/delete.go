package core

import (
	"context"
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// handleDelete removes all existing owner annotation references pointing to the object
func (h *Handler) handleDelete(obj podController) (reconcile.Result, error) {
	// Fetch all children with an owner annotation reference pointing to the object
	existing, err := h.getExistingChildren(obj)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("error fetching children: %v", err)
	}

	// Remove the owner annotation references from the children
	err = h.removeAnnoteReferences(obj, existing)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("error removing owner annotation references from children: %v", err)
	}

	// Remove the object's Finalizer and update if necessary
	copy := obj.DeepCopy()
	removeFinalizer(copy)
	if !reflect.DeepEqual(obj, copy) {
		err := h.Update(context.TODO(), copy.GetObject())
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("error updating Deployment: %v", err)
		}
	}

	return reconcile.Result{}, nil
}

// toBeDeleted checks whether the object has been marked for deletion
func toBeDeleted(obj metav1.Object) bool {
	// IsZero means the object hasn't been marked for deletion
	return !obj.GetDeletionTimestamp().IsZero()
}
