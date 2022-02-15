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

// MutationArgs is a named list of arguments used in a mutation
type MutationArgs map[string]DataTypeField

// FieldSpec defines the desired state of Field
type MutationSpec struct {
	// Type to return upon completion of the operation
	Type string `json:"type"`

	// Description of the mutation
	Description string `json:"description"`

	// Args received the the mutation
	Args MutationArgs `json:"args"`
}

// DataTypeBaseMutation is a mutation that is automatically supplied by the data
// store adapter implementation
type BaseMutation string

const (
	// DataTypeBaseMutationCreate is a base query that fetches a single record based
	// on ID
	BaseMutationCreate BaseMutation = "CREATE"

	// DataTypeBaseMutationList is a base query that fetches a list of records of
	// the specified type
	BaseMutationUpdate BaseMutation = "UPDATE"

	// DataTypeBaseMutationDelete is a base query that fetches a list of
	// records of the specified type
	BaseMutationDelete BaseMutation = "DELETE"
)

// FieldStatus defines the observed state of Field
type MutationStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Field is the Schema for the mutations API
// +k8s:openapi-gen=true
type Mutation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MutationSpec   `json:"spec,omitempty"`
	Status MutationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FieldList contains a list of Field
type MutationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mutation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mutation{}, &MutationList{})
}
