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

// BallistaApplicationSpec defines the desired state of BallistaApplication
type BallistaApplicationSpec struct {

	// BallistaVersion is the version of Ballista the application uses.
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
	BallistaPodSpec `json:",inline"`
	// PodName is the name of the scheduler pod that the user creates. This is used for the
	// in-cluster client mode in which the user creates a client pod where the scheduler of
	// the user application runs. It's an error to set this field if Mode is not
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
	BallistaPodSpec `json:",inline"`
	// Instances is the number of executor instances.
	// +optional
	// +kubebuilder:validation:Minimum=1
	Instances *int32 `json:"instances,omitempty"`
	// Ports settings for the pods, following the Kubernetes specifications.
	// +optional
	Ports []Port `json:"ports,omitempty"`
}

type BallistaPodSpec struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	Cores *int32 `json:"cores,omitempty"`
	// CoreLimit specifies a hard limit on CPU cores for the pod.
	// Optional
	CoreLimit *string `json:"coreLimit,omitempty"`
	// Memory is the amount of memory to request for the pod.
	// +optional
	Memory *string `json:"memory,omitempty"`
	// MemoryOverhead is the amount of off-heap memory to allocate in cluster mode, in MiB unless otherwise specified.
	// +optional
	MemoryOverhead *string `json:"memoryOverhead,omitempty"`
	// Image is the container image to use. Overrides Spec.Image if set.
	// +optional
	Image *string `json:"image,omitempty"`
	// Env carries the environment variables to add to the pod.
	// +optional
	Env []apiv1.EnvVar `json:"env,omitempty"`
	// EnvVars carries the environment variables to add to the pod.
	// Deprecated. Consider using `env` instead.
	// +optional
	EnvVars map[string]string `json:"envVars,omitempty"`
	// EnvFrom is a list of sources to populate environment variables in the container.
	// +optional
	EnvFrom []apiv1.EnvFromSource `json:"envFrom,omitempty"`
	// Labels are the Kubernetes labels to be added to the pod.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations are the Kubernetes annotations to be added to the pod.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// VolumeMounts specifies the volumes listed in ".spec.volumes" to mount into the main container's filesystem.
	// +optional
	VolumeMounts []apiv1.VolumeMount `json:"volumeMounts,omitempty"`
	// Affinity specifies the affinity/anti-affinity settings for the pod.
	// +optional
	Affinity *apiv1.Affinity `json:"affinity,omitempty"`
	// Tolerations specifies the tolerations listed in ".spec.tolerations" to be applied to the pod.
	// +optional
	Tolerations []apiv1.Toleration `json:"tolerations,omitempty"`
	// PodSecurityContext specifies the PodSecurityContext to apply.
	// +optional
	PodSecurityContext *apiv1.PodSecurityContext `json:"podSecurityContext,omitempty"`
	// SecurityContext specifies the container's SecurityContext to apply.
	// +optional
	SecurityContext *apiv1.SecurityContext `json:"securityContext,omitempty"`
	// SchedulerName specifies the scheduler that will be used for scheduling
	// +optional
	SchedulerName *string `json:"schedulerName,omitempty"`
	// Sidecars is a list of sidecar containers that run along side the main Ballista container.
	// +optional
	Sidecars []apiv1.Container `json:"sidecars,omitempty"`
	// InitContainers is a list of init-containers that run to completion before the main Ballista container.
	// +optional
	InitContainers []apiv1.Container `json:"initContainers,omitempty"`
	// HostNetwork indicates whether to request host networking for the pod or not.
	// +optional
	HostNetwork *bool `json:"hostNetwork,omitempty"`
	// NodeSelector is the Kubernetes node selector to be added to the scheduler and executor pods.
	// This field is mutually exclusive with nodeSelector at BallistaApplication level (which will be deprecated).
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// DnsConfig dns settings for the pod, following the Kubernetes specifications.
	// +optional
	DNSConfig *apiv1.PodDNSConfig `json:"dnsConfig,omitempty"`
	// Termination grace period seconds for the pod
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
	// ServiceAccount is the name of the custom Kubernetes service account used by the pod.
	// +optional
	ServiceAccount *string `json:"serviceAccount,omitempty"`
	// HostAliases settings for the pod, following the Kubernetes specifications.
	// +optional
	HostAliases []apiv1.HostAlias `json:"hostAliases,omitempty"`
	// ShareProcessNamespace settings for the pod, following the Kubernetes specifications.
	// +optional
	ShareProcessNamespace *bool `json:"shareProcessNamespace,omitempty"`
}

// Port represents the port definition in the pods objects.
type Port struct {
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	ContainerPort int32  `json:"containerPort"`
}

// BallistaApplicationStatus defines the observed state of BallistaApplication
type BallistaApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BallistaApplication is the Schema for the ballistaapplications API
// BallistaApplication represents a Ballista application running on and using Kubernetes as a cluster manager.
type BallistaApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BallistaApplicationSpec   `json:"spec,omitempty"`
	Status BallistaApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BallistaApplicationList contains a list of BallistaApplication
type BallistaApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BallistaApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BallistaApplication{}, &BallistaApplicationList{})
}
