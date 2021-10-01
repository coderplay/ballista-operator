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

package v1

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BallistaClusterSpec defines the desired state of BallistaCluster
type BallistaClusterSpec struct {

	// BallistaVersion is the version of Ballista the cluster uses.
	BallistaVersion string `json:"ballistaVersion"`

	// Image is the container image for the scheduler, executor, and init-container. Any custom container images for the
	// scheduler, executor, or init-container takes precedence over this.
	// +optional
	Image *string `json:"image,omitempty"`

	// Scheduler is the scheduler specification.
	Scheduler SchedulerSpec `json:"scheduler"`

	// Executor is the executor specification.
	Executor ExecutorSpec `json:"executor"`
}

// SchedulerSpec is specification of the scheduler.
type SchedulerSpec struct {
	apiv1.PodSpec `json:",inline"`
	// PodName is the name of the scheduler pod that the user creates. This is used for the
	// in-cluster client mode in which the user creates a client pod where the scheduler of
	// the user cluster runs. It's an error to set this field if Mode is not
	// in-cluster-client.
	// +optional
	// +kubebuilder:validation:Pattern=[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*
	PodName *string `json:"podName,omitempty"`
	// Lifecycle for running preStop or postStart commands
	// +optional
	Lifecycle *apiv1.Lifecycle `json:"lifecycle,omitempty"`
	// KubernetesMaster is the URL of the Kubernetes master used by the scheduler to manage executor pods and
	// other Kubernetes resources. Default to https://kubernetes.default.svc.
	// +optional
	KubernetesMaster *string `json:"kubernetesMaster,omitempty"`
	// ServiceAnnotations defines the annotations to be added to the Kubernetes headless service used by
	// executors to connect to the scheduler.
	// +optional
	ServiceAnnotations map[string]string `json:"serviceAnnotations,omitempty"`
	// Ports settings for the pods, following the Kubernetes specifications.
	// +optional
	Ports []Port `json:"ports,omitempty"`
}

// ExecutorSpec is specification of the executor.
type ExecutorSpec struct {
	apiv1.PodSpec `json:",inline"`
	// Instances is the number of executor instances.
	// +optional
	// +kubebuilder:validation:Minimum=1
	Instances *int32 `json:"instances,omitempty"`
	// Ports settings for the pods, following the Kubernetes specifications.
	// +optional
	Ports []Port `json:"ports,omitempty"`
}

// Port represents the port definition in the pods objects.
type Port struct {
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	ContainerPort int32  `json:"containerPort"`
}

// BallistaClusterStatus defines the observed state of BallistaCluster
type BallistaClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BallistaCluster is the Schema for the ballistaclusters API
// BallistaCluster represents a Ballista cluster running on and using Kubernetes as a cluster manager.
type BallistaCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BallistaClusterSpec   `json:"spec,omitempty"`
	Status BallistaClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BallistaClusterList contains a list of BallistaCluster
type BallistaClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BallistaCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BallistaCluster{}, &BallistaClusterList{})
}
