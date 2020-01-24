package vaultsyncer

import (
	"context"
	"time"

	operatorv1alpha1 "github.com/thatinfrastructureguy/vaultsync-operator/pkg/apis/operator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vaultsync"
)

var log = logf.Log.WithName("controller_vaultsyncer")

// Add creates a new VaultSyncer Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileVaultSyncer{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("vaultsyncer-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource VaultSyncer
	err = c.Watch(&source.Kind{Type: &operatorv1alpha1.VaultSyncer{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}

// blank assignment to verify that ReconcileVaultSyncer implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileVaultSyncer{}

// ReconcileVaultSyncer reconciles a VaultSyncer object
type ReconcileVaultSyncer struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a VaultSyncer object and makes changes based on the state read
// and what is in the VaultSyncer.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileVaultSyncer) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling VaultSyncer")

	// Fetch the VaultSyncer instance
	instance := &operatorv1alpha1.VaultSyncer{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	deploymentList = instance.Spec.DeploymentList

	reqLogger.Info("DestinationUpdated:", destinationUpdated)
	return reconcile.Result{RequeueAfter: time.Second * instance.Spec.RefreshRate}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *operatorv1alpha1.VaultSyncer) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
