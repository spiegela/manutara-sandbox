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

type NetworkProtocol string

const (
	// NetworkProtocolTCP
	NetworkProtocolTCP NetworkProtocol = "TCP"

	// NetworkProtocolUDP
	NetworkProtocolUDP NetworkProtocol = "UDP"
)

// DataStoreSpec defines the desired state of DataStore
type DataStoreSpec struct {
	// Name is the human readable of the data store
	Name string `json:"name"`

	// Configuration contains the driver specific configuration details for
	// creating a connection
	Configuration DataStoreConfiguration `json:"configuration"`
}

// DataStoreConfiguration contains the driver specific configuration details for
// creating a connection
type DataStoreConfiguration struct {
	// MySQLConfiguration is the detailed configuration for a MySQL connection
	// for this data store
	MySQL MySQLDataStoreConfiguration `json:"mysql"`
}

// MySQLConfiguration is the detailed configuration for a MySQL connection
type MySQLDataStoreConfiguration struct {
	// Username is the MySQL user used for authentication
	Username string `json:"username"`

	// Password is the MySQL password used for authentication
	Password string `json:"password"`

	// Protocol is the network protocol used for the collection
	Protocol string `json:"protocol"`

	// Port is the network port used for the connection
	Port int32 `json:"port"`

	// Database is the name of the database instance used for the connection
	Database string `json:"database"`

	// Options are the connection options used for the connection
	Options map[string]string `json:"options"`
}

// DataStoreStatus defines the observed state of DataStore
type DataStoreStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DataStore is the Schema for the datastores API
// +k8s:openapi-gen=true
type DataStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataStoreSpec   `json:"spec,omitempty"`
	Status DataStoreStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DataStoreList contains a list of DataStore
type DataStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataStore{}, &DataStoreList{})
}
