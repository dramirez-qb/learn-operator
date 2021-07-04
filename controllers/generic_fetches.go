package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//FetchHorizontalPodAutoscaler returns the HorizontalPodAutoscaler resource with the name in the namespace
func FetchHorizontalPodAutoscaler(name, namespace string, client client.Client) (*autoscalingv2beta2.HorizontalPodAutoscaler, error) {
	hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, hpa)
	return hpa, err
}

//FetchServiceAccount returns the ServiceAccount resource with the name in the namespace
func FetchServiceAccount(name, namespace string, client client.Client) (*corev1.ServiceAccount, error) {
	serviceaccount := &corev1.ServiceAccount{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, serviceaccount)
	return serviceaccount, err
}

//FetchRoleBinding returns the RoleBinding resource with the name in the namespace
func FetchRoleBinding(name, namespace string, client client.Client) (*rbacv1.RoleBinding, error) {
	rolebinding := &rbacv1.RoleBinding{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, rolebinding)
	return rolebinding, err
}

//FetchNetworkPolicy returns the NetworkPolicy resource with the name in the namespace
func FetchNetworkPolicy(name, namespace string, client client.Client) (*networkingv1.NetworkPolicy, error) {
	networkpolicy := &networkingv1.NetworkPolicy{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, networkpolicy)
	return networkpolicy, err
}

//FetchService returns the Service resource with the name in the namespace
func FetchService(name, namespace string, client client.Client) (*corev1.Service, error) {
	service := &corev1.Service{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, service)
	return service, err
}

//FetchService returns the Deployment resource with the name in the namespace
func FetchDeployment(name, namespace string, client client.Client) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, deployment)
	return deployment, err
}

//FetchPersistentVolumeClaim returns the PersistentVolumeClaim resource with the name in the namespace
func FetchPersistentVolumeClaim(name, namespace string, client client.Client) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, pvc)
	return pvc, err
}

//FetchCronJob returns the CronJob resource with the name in the namespace
func FetchCronJob(name, namespace string, client client.Client) (*v1beta1.CronJob, error) {
	cronJob := &v1beta1.CronJob{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cronJob)
	return cronJob, err
}

//FetchSecret returns the Secret resource with the name in the namespace
func FetchSecret(name, namespace string, client client.Client) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, secret)
	return secret, err
}

//FetchSecret returns the ConfigMap resource with the name in the namespace
func FetchConfigMap(name, namespace string, client client.Client) (*corev1.ConfigMap, error) {
	cfg := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cfg)
	return cfg, err
}

//FetchSecret returns the ConfigMap resource with the name in the namespace
func FetchStatefulSet(name, namespace string, client client.Client) (*appsv1.StatefulSet, error) {
	cfg := &appsv1.StatefulSet{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cfg)
	return cfg, err
}

//FetchSecret returns the ConfigMap resource with the name in the namespace
func FetchEndpoints(name, namespace string, client client.Client) (*corev1.Endpoints, error) {
	cfg := &corev1.Endpoints{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cfg)
	return cfg, err
}
