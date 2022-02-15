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

package datatype

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	manutarav1 "github.com/spiegela/manutara/pkg/apis/manutara/v1beta1"
)

var log = logf.Log.WithName("datatype-controller")

// Add creates a new DataType Controller and adds it to the Manager with default
// RBAC. The Manager will set fields on the Controller and Start it when the
// Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDataType{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("datatype-controller", mgr,
		controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	// Watch for changes to DataType
	err = c.Watch(&source.Kind{Type: &manutarav1.DataType{}},
		&handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}

var _ reconcile.Reconciler = &ReconcileDataType{}

// ReconcileDataType reconciles a DataType object
type ReconcileDataType struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a DataType object and makes
// changes based on the state read and what is in the DataType.Spec
// +kubebuilder:rbac:groups=manutara.spiegela.github.io,resources=datatypes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=manutara.spiegela.github.io,resources=datastores,verbs=get;list
// +kubebuilder:rbac:groups=manutara.spiegela.github.io,resources=datatypes/status,verbs=get;update;patch
func (r *ReconcileDataType) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the DataType instance
	instance := &manutarav1.DataType{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically
			// garbage collected. For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	err = r.createMutations(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileDataType) createMutations(instance *manutarav1.DataType) error {
	for _, mutType := range instance.Spec.BaseMutationsEnabled {
		err := r.Client.Create(context.TODO(), newMutation(mutType, instance))
		if errors.IsAlreadyExists(err) {
			err := r.Client.Update(context.TODO(), newMutation(mutType, instance))
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func mutationName(verb manutarav1.BaseMutation, typeName string) string {
	return strings.ToLower(string(verb)) + strings.Title(typeName)
}

func newMutation(mutType manutarav1.BaseMutation, instance *manutarav1.DataType) *manutarav1.Mutation {
	typeName := instance.GetName()
	return &manutarav1.Mutation{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "manutara.spiegela.github.io",
			Kind:       "Field",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      mutationName(mutType, typeName),
			Namespace: instance.GetNamespace(),
			Labels:    instance.GetLabels(),
		},
		Spec: manutarav1.MutationSpec{
			Type:        typeName,
			Description: mutationDescription(mutType, typeName),
			Args:        mutationArgs(mutType, instance),
		},
	}
}

func mutationArgs(mutType manutarav1.BaseMutation, instance *manutarav1.DataType) manutarav1.MutationArgs {
	mutArgs := manutarav1.MutationArgs{}
	switch mutType {
	case manutarav1.BaseMutationUpdate, manutarav1.BaseMutationDelete:
		mutArgs[manutarav1.IDFieldName] = manutarav1.DataTypeField{
			AllowNull:   false,
			Description: fmt.Sprintf("ID of the %s element", instance.GetName()),
			Type:        manutarav1.DataTypeFieldID,
			IsList:      false,
		}
	}
	switch mutType {
	case manutarav1.BaseMutationCreate, manutarav1.BaseMutationUpdate:
		mutArgs[string(instance.GetName())] = manutarav1.DataTypeField{
			AllowNull:       false,
			Description:     instance.Spec.Description,
			UserDefinedType: instance.GetName(),
			IsList:          false,
		}
	}
	return mutArgs
}

func mutationDescription(mutType manutarav1.BaseMutation, typeName string) string {
	switch mutType {
	case manutarav1.BaseMutationCreate:
		return fmt.Sprintf("Create a new %s", typeName)
	case manutarav1.BaseMutationUpdate:
		return fmt.Sprintf("Update an existing %s", typeName)
	case manutarav1.BaseMutationDelete:
		return fmt.Sprintf("Delete a %s", typeName)
	default:
		return ""
	}
}
