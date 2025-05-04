/*
Copyright 2025.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ArmorProxySpec defines the desired state of ArmorProxy
type ArmorProxySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// DeploymentName is the name of the deployment to be protected
	//+kubebuilder:validation:Required
	ServiceName string `json:"serviceName"`
	// ServiceNamespace is the namespace of the deployment to be protected
	//+kubebuilder:validation:Required
	ServiceNamespace string `json:"serviceNamespace"`

	// IngressHost is the optional hostname for ingress
	// +optional
	IngressHost string `json:"ingressHost,omitempty"`

	// IngressClassName is the optional ingress class name for the ingress resource
	// +optional
	IngressClassName string `json:"ingressClassName,omitempty"`

	// If auth to be enabled
	// +optional
	Oidc bool `json:"oidc,omitempty"`

	//If WAF to be enabled
	Waf bool `json:"waf,omitempty"`
}

// ArmorProxyStatus defines the observed state of ArmorProxy
type ArmorProxyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Conditions represent the latest available observations of the ArmorProxy
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ArmorProxy is the Schema for the armorproxies API
type ArmorProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArmorProxySpec   `json:"spec,omitempty"`
	Status ArmorProxyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ArmorProxyList contains a list of ArmorProxy
type ArmorProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ArmorProxy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ArmorProxy{}, &ArmorProxyList{})
}
