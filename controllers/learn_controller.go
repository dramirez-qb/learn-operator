/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	devopsv1alpha1 "github.com/dxas90/learn-operator/api/v1alpha1"
)

const statusFinalizer = "finalizer.status.devops.dxas90"

// LearnReconciler reconciles a Status object
type LearnReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=devops.dxas90,resources=learns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=devops.dxas90,resources=learns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=devops.dxas90,resources=learns/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments;statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services;configmaps;serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses;networkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Status object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *LearnReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	// Fetch the OptimalLeadNats instance
	instance := &devopsv1alpha1.Learn{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}
	// TODO Check number of nodes
	var replicas int32 = 2
	if instance.Spec.Replicas <= 0 {
		instance.Spec.Replicas = replicas
	}
	if instance.Spec.Image == "" {
		instance.Spec.Image = "dxas90/learn:latest"
	}

	if err := r.createResources(ctx, instance, req); err != nil {
		reqLogger.Error(err, "Failed to create the resource required for the Learn CR")
		return reconcile.Result{}, err
	}

	// if err := r.manageResources(instance); err != nil {
	// 	reqLogger.Error(err, "Failed to manage resource required for the Learn CR")
	// 	return reconcile.Result{}, err
	// }

	if err := r.createUpdateCRStatus(ctx, req); err != nil {
		reqLogger.Error(err, "Failed to create and update the status in the Learn CR")
		return reconcile.Result{}, err
	}

	isMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isMarkedToBeDeleted {
		if contains(instance.GetFinalizers(), statusFinalizer) {
			// Run finalization logic for statusFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeLearn(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}

			// Remove statusFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(instance, statusFinalizer)
			err := r.Update(ctx, instance)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if !contains(instance.GetFinalizers(), statusFinalizer) {
		if err := r.addFinalizer(ctx, instance); err != nil {
			return ctrl.Result{}, err
		}
	}

	reqLogger.Info("Skip reconcile: Status already exists", "Namespace", req.Namespace)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LearnReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devopsv1alpha1.Learn{}).
		Owns(&appsv1.Deployment{}).
		Owns(&autoscalingv2beta2.HorizontalPodAutoscaler{}).
		Owns(&v1.Service{}).
		Owns(&v1.ConfigMap{}).
		Owns(&v1.ServiceAccount{}).
		Owns(&networkingv1.NetworkPolicy{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&rbacv1.RoleBinding{}).
		Complete(r)
}

func (r *LearnReconciler) finalizeLearn(ctx context.Context, m *devopsv1alpha1.Learn) error {
	reqLogger := log.FromContext(ctx)
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.
	reqLogger.Info("Successfully finalized Status")
	return nil
}

func (r *LearnReconciler) addFinalizer(ctx context.Context, m *devopsv1alpha1.Learn) error {
	reqLogger := log.FromContext(ctx)
	reqLogger.Info("Adding Finalizer for the Status")
	controllerutil.AddFinalizer(m, statusFinalizer)

	// Update CR
	err := r.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update Learn with finalizer")
		return err
	}
	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
