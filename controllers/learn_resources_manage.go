package controllers

import (
	"context"

	devopsv1alpha1 "github.com/dxas90/learn-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
)

// manageResources will ensure that the resource are with the expected values in the cluster
func (r *LearnReconciler) manageResources(cr *devopsv1alpha1.Learn) error {
	ctx := context.Background()
	dp := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}, dp)
	if err != nil {
		return err
	}

	// Ensure the deployment size is the same as the spec
	return r.ensureDepSize(cr, dp)
}

// ensureDepSize will ensure that the quanity of instances in the cluster for the Database deployment is the same defined in the CR
func (r *LearnReconciler) ensureDepSize(cr *devopsv1alpha1.Learn, dep *appsv1.Deployment) error {
	size := cr.Spec.Replicas
	if dep.Spec.Replicas != &size {
		// Set the number of Replicas spec in the CR
		dep.Spec.Replicas = &size
		if err := r.Update(context.TODO(), dep); err != nil {
			return err
		}
	}
	return nil
}
