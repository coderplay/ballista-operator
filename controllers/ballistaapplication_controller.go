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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ballistaminzhouinfov1 "github.com/coderplay/ballista-operator/api/v1"
)

// BallistaApplicationReconciler reconciles a BallistaApplication object
type BallistaApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ballista.minzhou.info,resources=ballistaapplications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ballista.minzhou.info,resources=ballistaapplications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ballista.minzhou.info,resources=ballistaapplications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *BallistaApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)



	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BallistaApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ballistaminzhouinfov1.BallistaApplication{}).
		Complete(r)
}
