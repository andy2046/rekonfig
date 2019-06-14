package deployment

import (
	"context"

	"github.com/andy2046/rekonfig/pkg/core"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new Deployment Controller and adds it to the Manager.
// The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDeployment{
		scheme:  mgr.GetScheme(),
		handler: core.NewHandler(mgr.GetClient(), mgr.GetRecorder("rekonfig")),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("deployment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Deployment
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch ConfigMaps owned by a Deployment
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(mapObject handler.MapObject) []reconcile.Request {
			ns, req := mapObject.Meta.GetNamespace(), []reconcile.Request{}
			for annoteKey := range mapObject.Meta.GetAnnotations() {
				if name := core.GetNameFromAnnoteKey(annoteKey); name != "" {
					req = append(req, reconcile.Request{
						NamespacedName: types.NamespacedName{
							Name:      name,
							Namespace: ns,
						},
					})
				}
			}
			return req
		}),
	})
	if err != nil {
		return err
	}

	// Watch Secrets owned by a Deployment
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(mapObject handler.MapObject) []reconcile.Request {
			ns, req := mapObject.Meta.GetNamespace(), []reconcile.Request{}
			for annoteKey := range mapObject.Meta.GetAnnotations() {
				if name := core.GetNameFromAnnoteKey(annoteKey); name != "" {
					req = append(req, reconcile.Request{
						NamespacedName: types.NamespacedName{
							Name:      name,
							Namespace: ns,
						},
					})
				}
			}
			return req
		}),
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileDeployment{}

// ReconcileDeployment reconciles a Deployment object
type ReconcileDeployment struct {
	scheme  *runtime.Scheme
	handler *core.Handler
}

// Reconcile reads that state of the cluster for a Deployment object and makes changes based on the state read
// and what is in the StatefulSet.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileDeployment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Deployment instance
	instance := &appsv1.Deployment{}
	err := r.handler.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	return r.handler.HandleDeployment(instance)
}
