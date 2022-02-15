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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DataTypeFieldType string

const DataStoreSelectorLabel = "datastore.manutara.spiegela.github.io/name"

const (
	// DataTypeFieldID is a GraphQL ID data type
	DataTypeFieldID DataTypeFieldType = "ID"

	// DataTypeFieldString is a GraphQL String data type
	DataTypeFieldString DataTypeFieldType = "String"

	// DataTypeFieldInt is a GraphQL Int data type
	DataTypeFieldInt DataTypeFieldType = "Int"

	// DataTypeFieldFloat is a GraphQL Float data type
	DataTypeFieldFloat DataTypeFieldType = "Float"

	// DataTypeFieldBoolean is a GraphQL Boolean data type
	DataTypeFieldBoolean DataTypeFieldType = "Boolean"

	// DataTypeFieldDate is a GraphQL Date data type
	DataTypeFieldDate DataTypeFieldType = "Date"
)

// DataTypeBaseQuery is query that is automatically supplied by the data store
// adapter implementation
type DataTypeBaseQuery string

const (
	// DataTypeBaseQueryGet is a base query that fetches a single record based
	// on ID
	DataTypeBaseQueryGet DataTypeBaseQuery = "GET"

	// DataTypeBaseQueryList is a base query that fetches a list of records of
	// the specified type
	DataTypeBaseQueryList DataTypeBaseQuery = "LIST"

	// DataTypeBaseQueryConnection is a base query that fetches a list of
	// records of the specified type
	DataTypeBaseQueryConnection DataTypeBaseQuery = "CONNECTION"
)

// Name of the ID field
const IDFieldName = "ID"

// DataTypeFieldUnionType is a list of types that are valid as values for a
// field
type DataTypeFieldUnionType []string

// DataTypeFieldEnumType is a list of values that are valid for a field
type DataTypeFieldEnumType []string

// DataTypeField defines a GraphQL field for the data type
type DataTypeField struct {
	// Description is the user presented short description of the field
	Description string `json:"description"`

	// Type is the data type of the field. If the field references a user
	// defined type, such as an interface or a GraphQL Data Type, this should be
	// empty
	Type DataTypeFieldType `json:"type"`

	// UserDefinedType is a type that is defined within the GraphQL schema and
	// should either be an DataType name or an interface
	UserDefinedType string `json:"userDefinedType"`

	// UnionType is a list of types that are valid for for this field
	UnionType DataTypeFieldUnionType `json:"union"`

	// EnumType is a list of valid values that are valid for for this field
	EnumType DataTypeFieldEnumType `json:"enum"`

	// IsList specifies if this type is a list of objects
	IsList bool `json:"isList"`

	// AllowEmpty defines if an empty list is allowed in the case that this is
	// a list field
	AllowEmpty bool `json:"allowEmpty"`

	// AllowNull defines if a null value is allowed for the field
	AllowNull bool `json:"allowNull"`
}

// DataTypeFields is a named list of fields used in a query
type DataTypeFields map[string]DataTypeField

// DataTypeSpec defines the desired state of DataType
type DataTypeSpec struct {
	// ServiceName is the name of the service which will present type's GraphQL
	// queries and mutations
	ServiceName string `json:"serviceName"`

	// DataStores is a set label selector used to select data stores for
	DataStores []string `json:"dataStores"`

	// Description is the user presented short description of the data type
	Description string `json:"description"`

	// Fields is the list of fields making up the data type
	Fields DataTypeFields `json:"fields"`

	// BaseQueriesEnabled is a list of base queries enabled for this type
	BaseQueriesEnabled []DataTypeBaseQuery `json:"baseQueriesEnabled"`

	// BaseMutationsEnabled is a list of base mutations enabled for this type
	BaseMutationsEnabled []BaseMutation `json:"baseMutationsEnabled"`
}

// DataTypeStatus defines the observed state of DataType
type DataTypeStatus struct {
	// CurrentRevision is current schema version that has been applied
	CurrentRevision string `json:"observedRevision"`

	// UpdateRevision is current schema version that has been applied
	UpdateRevision string `json:"updateRevision"`

	Conditions []DataTypeCondition `json:"conditions"`
}

type DataTypeConditionType string

const (
	// DataTypeMigrating describes a condition in which a migration is currently
	// progressing on a data type
	DataTypeMigrating DataTypeConditionType = "Migrating"

	// DataTypeMigrationComplete describes a condition in which a migration
	// has successfully completed
	DataTypeMigrationComplete DataTypeConditionType = "MigrationComplete"

	// DataTypeMigrationFailed describes a condition in which a migration
	// has exited with a failure
	DataTypeMigrationFailed DataTypeConditionType = "MigrationFailed"
)

// DataTypeCondition describes the state of a datatype at a certain point.
type DataTypeCondition struct {
	// Type of statefulset condition.
	Type DataTypeConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`

	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DataType is the Schema for the datatypes API
// +k8s:openapi-gen=true
type DataType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataTypeSpec   `json:"spec,omitempty"`
	Status DataTypeStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DataTypeList contains a list of DataType
type DataTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataType `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataType{}, &DataTypeList{})
}
