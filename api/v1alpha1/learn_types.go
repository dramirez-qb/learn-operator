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

package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LearnSpec defines the desired state of Learn
type LearnSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Learn. Edit learn_types.go to remove/update
	Foo string `json:"foo,omitempty"`
	// Image to deploy
	// image is the container image to run.  Image must have a tag.
	// +kubebuilder:validation:Pattern=".+:.+"
	// +kubebuilder:default:="dxas90/learn:latest"
	Image string `json:"image,omitempty"`
	// Replicas that we need
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Max=15
	// +kubebuilder:validation:default:=2
	Replicas int32 `json:"replicas,omitempty"`
}

// LearnStatus defines the observed state of Learn
type LearnStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Status of the Status Deployment created and managed by it
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Deployment Status"
	DeploymentStatus appsv1.DeploymentStatus `json:"deploymentStatus"`

	// Status of the Status Service created and managed by it
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Service Status"
	ServiceStatus corev1.ServiceStatus `json:"serviceStatus"`

	// Status of the Status HorizontalPodAutoscaler created and managed by it
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Horizontal Pod Autoscaler Status"
	// HorizontalPodAutoscalerStatus autoscalingv2beta2.HorizontalPodAutoscalerStatus `json:"hpaStatus"`

	Status string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas",description="The number of jobs launched"
//kubebuilder:printcolumn:name="Deployment Status",type="string",JSONPath=".status.deploymentStatus",description="Learn Deployment Status"
//kubebuilder:printcolumn:name="Service Status",type="string",JSONPath=".status.serviceStatus",description="Learn Service Status"
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="Learn Status"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
//Learn is the Schema for the learns API
type Learn struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LearnSpec   `json:"spec,omitempty"`
	Status LearnStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LearnList contains a list of Learn
type LearnList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Learn `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Learn{}, &LearnList{})
}
