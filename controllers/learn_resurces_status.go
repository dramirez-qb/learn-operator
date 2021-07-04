package controllers

import (
	"context"

	"fmt"

	devopsv1alpha1 "github.com/dxas90/learn-operator/api/v1alpha1"

	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const statusOk = "OK"

//createUpdateCRStatus will create and update the status in the CR applied in the cluster
func (r *LearnReconciler) createUpdateCRStatus(ctx context.Context, request reconcile.Request) error {
	reqLogger := log.FromContext(ctx)
	reqLogger.Info("Create/Update Status status ...")

	if err := r.updateDeploymentStatus(request); err != nil {
		reqLogger.Error(err, "Failed to create Deployment Status")
		return err
	}

	if err := r.updateServiceStatus(request); err != nil {
		reqLogger.Error(err, "Failed to create Service Status")
		return err
	}

	// if err := r.updateHorizontalPodAutoscalerStatus(request); err != nil {
	// 	reqLogger.Error(err, "Failed to create HorizontalPodAutoscaler Status")
	// 	return err
	// }

	if err := r.updateStatus(request); err != nil {
		reqLogger.Error(err, "Failed to create Status")
		return err
	}
	return nil
}

//updateStatusStatus returns error when status regards the all required resource could not be updated
func (r *LearnReconciler) updateStatus(request reconcile.Request) error {
	ctx := context.TODO()
	Status := &devopsv1alpha1.Learn{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}, Status)
	if err != nil {
		return err
	}

	statusMsgUpdate := statusOk
	// Check if all required resource were created and found
	if err := r.isAllCreated(Status); err != nil {
		statusMsgUpdate = err.Error()
	}

	// Check if BackupStatus was changed, if yes update it
	if err := r.insertUpdateGeneralStatus(Status, statusMsgUpdate); err != nil {
		return err
	}
	return nil
}

// Check if General Status was changed, if yes update it
func (r *LearnReconciler) insertUpdateGeneralStatus(cr *devopsv1alpha1.Learn, statusMsgUpdate string) error {
	ctx := context.TODO()
	if statusMsgUpdate != cr.Status.Status {
		cr.Status.Status = statusMsgUpdate
		if err := r.Status().Update(ctx, cr); err != nil {
			return err
		}
	}
	return nil
}

//updateDeploymentStatus returns error when status regards the deployment resource could not be updated
func (r *LearnReconciler) updateDeploymentStatus(request reconcile.Request) error {
	ctx := context.TODO()
	Status := &devopsv1alpha1.Learn{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}, Status)
	if err != nil {
		return err
	}

	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}, deployment)

	if err != nil {
		return err
	}

	// Check if Deployment Status was changed, if yes update it
	if err := r.insertUpdateDeploymentStatus(deployment, Status); err != nil {
		return err
	}

	return nil
}

// insertUpdateDeploymentStatus will check if Deployment status changed, if yes then and update it
func (r *LearnReconciler) insertUpdateDeploymentStatus(deployment *appsv1.Deployment, cr *devopsv1alpha1.Learn) error {
	ctx := context.TODO()
	if !reflect.DeepEqual(deployment.Status, cr.Status.DeploymentStatus) {
		cr.Status.DeploymentStatus = deployment.Status
		if err := r.Status().Update(ctx, cr); err != nil {
			return err
		}
	}
	return nil
}

//updateHorizontalPodAutoscalerStatus returns error when status regards the HorizontalPodAutoscaler resource could not be updated
func (r *LearnReconciler) updateHorizontalPodAutoscalerStatus(request reconcile.Request) error {
	ctx := context.TODO()
	learn := &devopsv1alpha1.Learn{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}, learn)
	if err != nil {
		return err
	}
	hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}, hpa)
	if err != nil {
		return err
	}

	// Check if HorizontalPodAutoscaler Status was changed, if yes update it
	if err := r.insertUpdateHorizontalPodAutoscalerStatus(hpa, learn); err != nil {
		return err
	}

	return nil
}

// insertUpdateDeploymentStatus will check if HorizontalPodAutoscaler status changed, if yes then and update it
func (r *LearnReconciler) insertUpdateHorizontalPodAutoscalerStatus(hpaStatus *autoscalingv2beta2.HorizontalPodAutoscaler, cr *devopsv1alpha1.Learn) error {
	ctx := context.TODO()
	if !reflect.DeepEqual(hpaStatus.Status, cr.Status.ServiceStatus) {
		// cr.Status.HorizontalPodAutoscalerStatus = hpaStatus.DeepCopy().Status
		if err := r.Status().Update(ctx, cr); err != nil {
			return fmt.Errorf("%v\n", hpaStatus.Status.Conditions)
		}
	}
	return nil
}

//updateServiceStatus returns error when status regards the service resource could not be updated
func (r *LearnReconciler) updateServiceStatus(request reconcile.Request) error {
	ctx := context.TODO()
	learn := &devopsv1alpha1.Learn{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}, learn)
	if err != nil {
		return err
	}
	srv := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}, srv)
	if err != nil {
		return err
	}

	// Check if Service Status was changed, if yes update it
	if err := r.insertUpdateServiceStatus(srv, learn); err != nil {
		return err
	}

	return nil
}

// insertUpdateDeploymentStatus will check if Service status changed, if yes then and update it
func (r *LearnReconciler) insertUpdateServiceStatus(serviceStatus *corev1.Service, cr *devopsv1alpha1.Learn) error {
	ctx := context.TODO()
	if !reflect.DeepEqual(serviceStatus.Status, cr.Status.ServiceStatus) {
		cr.Status.ServiceStatus = serviceStatus.Status
		if err := r.Status().Update(ctx, cr); err != nil {
			return err
		}
	}
	return nil
}

//validateBackup returns error when some requirement is missing
func (r *LearnReconciler) isAllCreated(cr *devopsv1alpha1.Learn) error {
	// Check if the Deployment was created
	ctx := context.TODO()
	err := r.Get(ctx, types.NamespacedName{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}, &devopsv1alpha1.Learn{})
	if err != nil {
		return err
	}

	err = r.Get(ctx, types.NamespacedName{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}, &appsv1.Deployment{})

	if err != nil {
		return fmt.Errorf("Error: Deployment is missing.")
	}

	err = r.Get(ctx, types.NamespacedName{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}, &autoscalingv2beta2.HorizontalPodAutoscaler{})

	if err != nil {
		return fmt.Errorf("Error: HorizontalPodAutoscaler is missing.")
	}

	err = r.Get(ctx, types.NamespacedName{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}, &corev1.Service{})

	if err != nil {
		return fmt.Errorf("Error: Service is missing.")
	}
	return err
}
