package core

import (
	"context"
	"fmt"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// Handler performs main controller business logic
type Handler struct {
	client.Client
	recorder record.EventRecorder
}

// NewHandler creates a new Handler instance
func NewHandler(c client.Client, r record.EventRecorder) *Handler {
	return &Handler{Client: c, recorder: r}
}

// HandleDeployment is called by the deployment controller to reconcile deployments
func (h *Handler) HandleDeployment(instance *appsv1.Deployment) (reconcile.Result, error) {
	return h.handlePodController(&deployment{Deployment: instance})
}

// HandleStatefulSet is called by the StatefulSet controller to reconcile StatefulSets
func (h *Handler) HandleStatefulSet(instance *appsv1.StatefulSet) (reconcile.Result, error) {
	return h.handlePodController(&statefulset{StatefulSet: instance})
}

// HandleDaemonSet is called by the DaemonSet controller to reconcile DaemonSets
func (h *Handler) HandleDaemonSet(instance *appsv1.DaemonSet) (reconcile.Result, error) {
	return h.handlePodController(&daemonset{DaemonSet: instance})
}

// handlePodController reconciles the state of a podController
func (h *Handler) handlePodController(instance podController) (reconcile.Result, error) {
	log := logf.Log.WithName("rekonfig")

	// If the required annotation isn't present, ignore the instance
	if !hasRequiredAnnotation(instance) {
		// Perform deletion logic if the finalizer is present on the object
		if hasFinalizer(instance) {
			log.V(0).Info("Required annotation removed from instance, cleaning up orphans", "namespace",
				instance.GetNamespace(), "name", instance.GetName())
			return h.handleDelete(instance)
		}
		return reconcile.Result{}, nil
	}

	// If the instance is marked for deletion, run cleanup process
	if toBeDeleted(instance) {
		log.V(0).Info("Instance marked for deletion, cleaning up orphans", "namespace",
			instance.GetNamespace(), "name", instance.GetName())
		return h.handleDelete(instance)
	}

	// Get all children that have an owner annotation reference pointing to this instance
	existing, err := h.getExistingChildren(instance)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("error fetching existing children: %v", err)
	}

	// Get all children that the instance currently references
	current, err := h.getCurrentChildren(instance)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("error fetching current children: %v", err)
	}

	// Reconcile the owner annotation references on the existing and current children
	err = h.updateAnnoteReferences(instance, existing, current)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("error updating owner annotation references: %v", err)
	}

	hash, err := calculateConfigHash(current)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("error calculating configuration hash: %v", err)
	}

	// Update the desired state of the Deployment in a DeepCopy
	copy := instance.DeepCopy()
	setConfigHash(copy, hash)
	addFinalizer(copy)

	// If the desired state doesn't match the existing state, update it
	if !reflect.DeepEqual(instance, copy) {
		log.V(0).Info("Updating instance hash", "namespace", instance.GetNamespace(), "name", instance.GetName(), "hash", hash)
		h.recorder.Eventf(copy.GetObject(), corev1.EventTypeNormal, "ConfigChanged", "Configuration hash updated to %s", hash)
		err := h.Update(context.TODO(), copy.GetObject())
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("error updating instance %s/%s: %v", instance.GetNamespace(), instance.GetName(), err)
		}
	}

	return reconcile.Result{}, nil
}
