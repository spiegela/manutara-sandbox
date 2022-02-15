/*
Copyright 2019 Aaron Spiegel.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DataStoreBindingSpec defines the desired state of DataStoreBinding
type DataStoreBindingSpec struct {
	//
}

// DataStoreBindingStatus defines the observed state of DataStoreBinding
type DataStoreBindingStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DataStoreBinding is the Schema for the datastorebindings API
// +k8s:openapi-gen=true
type DataStoreBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataStoreBindingSpec   `json:"spec,omitempty"`
	Status DataStoreBindingStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DataStoreBindingList contains a list of DataStoreBinding
type DataStoreBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataStoreBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataStoreBinding{}, &DataStoreBindingList{})
}
