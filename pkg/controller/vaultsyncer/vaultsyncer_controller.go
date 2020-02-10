/*
 * Copyright 2020 Kulkarni, Ashish <thatInfrastructureGuy@gmail.com>
 * Author: Ashish Kulkarni
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vaultsyncer

import (
	"context"

	operatorv1alpha1 "github.com/thatinfrastructureguy/vaultsync-operator/pkg/apis/operator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
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

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner VaultSyncer
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorv1alpha1.VaultSyncer{},
	})

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

	if len(instance.Spec.ProviderCredsSecret) == 0 {
		instance.Spec.ProviderCredsSecret = "provider-credentials"
	}

	if len(instance.Spec.Image) == 0 {
		instance.Spec.Image = "thatinfrastructureguy/vaultsync:v0.0.14"
	}

	if len(instance.Spec.SecretName) == 0 {
		instance.Spec.SecretName = instance.Spec.VaultName
	}

	if len(instance.Spec.SecretNamespace) == 0 {
		instance.Spec.SecretNamespace = "default"
	}

	// Define a new Pod object
	podObject := newPodForCR(instance)

	// Update VaultSync Status
	instance.Status.SecretName = instance.Spec.SecretName
	instance.Status.SecretNamespace = instance.Spec.SecretNamespace
	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		reqLogger.Error(err, "Failed to update VaultSyncer status")
		return reconcile.Result{}, err
	}

	// Set VaultSyncer instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, podObject, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: podObject.Name, Namespace: podObject.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", podObject.Namespace, "Pod.Name", podObject.Name)
		err = r.client.Create(context.TODO(), podObject)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
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
					Name:  "vaultsync",
					Image: cr.Spec.Image,
					Env: []corev1.EnvVar{
						corev1.EnvVar{Name: "PROVIDER", Value: cr.Spec.Provider},
						corev1.EnvVar{Name: "VAULT_NAME", Value: cr.Spec.VaultName},
						corev1.EnvVar{Name: "CONSUMER", Value: cr.Spec.Consumer},
						corev1.EnvVar{Name: "SECRET_NAMESPACE", Value: cr.Spec.SecretNamespace},
						corev1.EnvVar{Name: "SECRET_NAME", Value: cr.Spec.SecretName},
						corev1.EnvVar{Name: "DEPLOYMENT_LIST", Value: cr.Spec.DeploymentList},
						corev1.EnvVar{Name: "STATEFULSET_LIST", Value: cr.Spec.StatefulsetList},
						corev1.EnvVar{Name: "REFRESH_RATE", Value: cr.Spec.RefreshRate},
						corev1.EnvVar{Name: "CONVERT_HYPHENS_TO_UNDERSCORES", Value: cr.Spec.ConvertHyphensToUnderscores},
					},
					EnvFrom: []corev1.EnvFromSource{
						{
							SecretRef: &corev1.SecretEnvSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: cr.Spec.ProviderCredsSecret,
								},
							},
						},
					},
				},
			},
			ServiceAccountName: "vaultsync-operator",
		},
	}
}
