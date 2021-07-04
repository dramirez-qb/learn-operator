package controllers

import (
	"context"

	devopsv1alpha1 "github.com/dxas90/learn-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	configMapData    = make(map[string]string)
	envConfigMapData = make(map[string]string)
)

func init() {
	configMapData = map[string]string{
		"REDIS_DSN":   "redis://redis:6379?timeout=0.5",
		"MONGODB_URL": "mongodb://mongodb:27017",
		"MAILER_URL":  "smtp://mail-server:1025",
	}
	envConfigMapData = map[string]string{
		"DEVOPS_APP": "Learn",
	}
}

func (r *LearnReconciler) createResources(ctx context.Context, cr *devopsv1alpha1.Learn, request reconcile.Request) error {
	reqLogger := log.FromContext(ctx)
	reqLogger.Info("Creating Status resources ...")

	// Check if service for the app exist, if not create one
	if err := r.createConfigMapsCR(cr); err != nil {
		reqLogger.Error(err, "Failed to create ConfigMaps")
		return err
	}

	// Check if ServiceAccount for the app exist, if not create one
	if err := r.createServiceAccountCR(cr); err != nil {
		reqLogger.Error(err, "Failed to create ServiceAccount")
		return err
	}

	// Check if Deployment for the app exist, if not create one
	if err := r.createDeploymentCR(cr); err != nil {
		reqLogger.Error(err, "Failed to create Deployment")
		return err
	}

	// Check if createServiceCR for the app exist, if not create one
	if err := r.createServiceCR(cr); err != nil {
		reqLogger.Error(err, "Failed to create Service")
		return err
	}

	// Check if HPA for the app exist, if not create one
	if err := r.createHpaCR(cr); err != nil {
		reqLogger.Error(err, "Failed to create HorizontalPodAutoscaler")
		return err
	}

	return nil
}

// Check if Service for the app exist, if not create one
func (r *LearnReconciler) createServiceCR(cr *devopsv1alpha1.Learn) error {
	ctx := context.Background()
	srv := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}, srv)
	if err != nil {
		if err := r.Create(ctx, NewService(cr, r.Scheme)); err != nil {
			return err
		}
	}
	return nil
}

// Check if ConfigMap for the app exist, if not create one
func (r *LearnReconciler) createConfigMapsCR(cr *devopsv1alpha1.Learn) error {
	ctx := context.Background()
	cm := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      cr.Name + "-conf",
		Namespace: cr.Namespace,
	}, cm)
	if err != nil {
		if err := r.Create(ctx, NewConfigMapCR(cr, "-conf", configMapData, r.Scheme)); err != nil {
			return err
		}
	}
	return nil
}

// Check if ServiceAccount for the app exist, if not create one
func (r *LearnReconciler) createServiceAccountCR(cr *devopsv1alpha1.Learn) error {
	ctx := context.Background()
	sa := &corev1.ServiceAccount{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      cr.Name + "-sa",
		Namespace: cr.Namespace,
	}, sa)
	if err != nil {
		if err := r.Create(ctx, NewServiceAccount(cr, r.Scheme)); err != nil {
			return err
		}
	}
	return nil
}

// Check if Deployment for the app exist, if not create one
func (r *LearnReconciler) createDeploymentCR(cr *devopsv1alpha1.Learn) error {
	ctx := context.Background()
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}, deployment)
	if err != nil {
		if err := r.Create(ctx, NewDeploymentForCR(cr, r.Scheme)); err != nil {
			return err
		}
	}
	return nil
}

// Check if HPA for the app exist, if not create one
func (r *LearnReconciler) createHpaCR(cr *devopsv1alpha1.Learn) error {
	ctx := context.Background()
	hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}, hpa)
	if err != nil {
		if err := r.Create(ctx, NewHorizontalPodAutoscalerForCR(cr, r.Scheme)); err != nil {
			return err
		}
	}
	return nil
}

// newPodForCR returns a configMap with the value of Data the cr
func NewConfigMapCR(cr *devopsv1alpha1.Learn, suffix string, Data map[string]string, scheme *runtime.Scheme) *corev1.ConfigMap {
	labels := map[string]string{
		"app":    cr.Name,
		"devops": cr.Name,
	}
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + suffix,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		Data: Data,
	}
	controllerutil.SetControllerReference(cr, configMap, scheme)
	return configMap
}

// Returns the service object for the Learn app
func NewService(cr *devopsv1alpha1.Learn, scheme *runtime.Scheme) *corev1.Service {
	labels := map[string]string{
		"app":    cr.Name,
		"devops": cr.Name,
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: "web",
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8080,
					},
					Port:     8080,
					Protocol: "TCP",
				},
			},
		},
	}
	// Set Status cr as the owner and controller
	controllerutil.SetControllerReference(cr, service, scheme)
	return service
}

// NewDeploymentForCR returns a deployment name/namespace as the cr
func NewDeploymentForCR(cr *devopsv1alpha1.Learn, scheme *runtime.Scheme) *appsv1.Deployment {
	labels := map[string]string{
		"app":    cr.Name,
		"devops": cr.Name,
	}
	var defaultMode int32 = 0755
	var defaultFSGroup int64 = 65534
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      cr.Name,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Replicas: &cr.Spec.Replicas,
			Strategy: appsv1.DeploymentStrategy{
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 0,
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 2,
					},
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: cr.Name + "-conf",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: cr.Name + "-conf",
									},
									DefaultMode: &defaultMode,
								},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:            "pull-secrets",
							Image:           "busybox",
							ImagePullPolicy: corev1.PullIfNotPresent,
							EnvFrom: []corev1.EnvFromSource{
								{
									ConfigMapRef: &corev1.ConfigMapEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: cr.Name + "-conf",
										},
									},
								},
							},
							Env: []corev1.EnvVar{
								{
									Name: "POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "status.podIP",
										},
									},
								},
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.namespace",
										},
									},
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("5m"),
									corev1.ResourceMemory: resource.MustParse("16Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("5m"),
									corev1.ResourceMemory: resource.MustParse("16Mi"),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            cr.Name,
							Image:           cr.Spec.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "web",
									ContainerPort: 8080,
									Protocol:      "TCP",
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									ConfigMapRef: &corev1.ConfigMapEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: cr.Name + "-conf",
										},
									},
								},
							},
							Env: []corev1.EnvVar{
								{
									Name: "POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "status.podIP",
										},
									},
								},
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.name",
										},
									},
								},
								{
									Name: "MY_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.namespace",
										},
									},
								},
								{
									Name: "USER",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.name",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      cr.Name + "-conf",
									MountPath: "/conf",
									ReadOnly:  true,
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/healthz",
										Port: intstr.IntOrString{
											Type:   intstr.String,
											StrVal: "web",
										},
									},
								},
								InitialDelaySeconds: 3,
								TimeoutSeconds:      2,
								FailureThreshold:    5,
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("10m"),
									corev1.ResourceMemory: resource.MustParse("48Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("10m"),
									corev1.ResourceMemory: resource.MustParse("48Mi"),
								},
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: "File",
						},
					},
					DNSPolicy:     corev1.DNSClusterFirst,
					RestartPolicy: corev1.RestartPolicyAlways,
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup: &defaultFSGroup,
					},
					ServiceAccountName: cr.Name + "-sa",
				},
			},
		},
	}
	controllerutil.SetControllerReference(cr, deployment, scheme)
	return deployment
}

// Returns the ServiceAccount object for the Learn app
func NewServiceAccount(cr *devopsv1alpha1.Learn, scheme *runtime.Scheme) *corev1.ServiceAccount {
	labels := map[string]string{
		"app":    cr.Name,
		"devops": cr.Name,
	}
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-sa",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
	}
	// Set Status cr as the owner and controller
	controllerutil.SetControllerReference(cr, serviceAccount, scheme)
	return serviceAccount
}

// Returns the ServiceAccount object for the Learn app
func NewHorizontalPodAutoscalerForCR(cr *devopsv1alpha1.Learn, scheme *runtime.Scheme) *autoscalingv2beta2.HorizontalPodAutoscaler {
	labels := map[string]string{
		"app":    cr.Name,
		"devops": cr.Name,
	}
	// TODO change from config
	var MinReplicas int32 = 1
	var MaxReplicas int32 = 5
	var averageCPUUtilization int32 = 80
	var averageMemoryValue int64 = 50
	hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		}, Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
				Kind:       "Deployment",
				APIVersion: "apps/v1",
				Name:       cr.Name,
			},
			MinReplicas: &MinReplicas,
			MaxReplicas: MaxReplicas,
			Metrics: []autoscalingv2beta2.MetricSpec{
				{
					Type: "Resource",
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: "cpu",
						Target: autoscalingv2beta2.MetricTarget{
							Type:               autoscalingv2beta2.UtilizationMetricType,
							AverageUtilization: &averageCPUUtilization,
						},
					},
				},
				{
					Type: "Resource",
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: "memory",
						Target: autoscalingv2beta2.MetricTarget{
							Type:         autoscalingv2beta2.AverageValueMetricType,
							AverageValue: resource.NewQuantity(averageMemoryValue, resource.MustParse("Mi").Format),
						},
					},
				},
			},
		},
	}
	// Set Status cr as the owner and controller
	controllerutil.SetControllerReference(cr, hpa, scheme)
	return hpa
}
