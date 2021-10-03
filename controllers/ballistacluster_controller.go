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
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	k8sapiv1 "k8s.io/api/core/v1"
	v1 "github.com/coderplay/ballista-operator/api/v1"
)

// BallistaClusterReconciler reconciles a BallistaCluster object
type BallistaClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ballista.minzhou.info,resources=ballistaclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ballista.minzhou.info,resources=ballistaclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ballista.minzhou.info,resources=ballistaclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *BallistaClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	log := log.FromContext(ctx)

	var cluster = &v1.BallistaCluster{}
	if err := r.Get(ctx, req.NamespacedName, cluster); err != nil {
		log.Error(err, "unable to fetch BallistaCluster")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	clusterCopy := cluster.DeepCopy()

	switch clusterCopy.Status.ClusterState.State {
	case v1.NewState:
		clusterCopy = r.startBallistaCluster(ctx, clusterCopy)
	case v1.Pending, v1.RunningState, v1.UnknownState:
		if err := r.getAndUpdateClusterState(ctx, clusterCopy); err != nil {
			return ctrl.Result{}, err
		}
	}



	return ctrl.Result{}, err
}

var (
	podOwnerKey         = ".metadata.controller"
	podBallistaRoleKey  = "ballista-role"
	apiGVStr            = v1.GroupVersion.String()
)


// SetupWithManager sets up the controller with the Manager.
func (r *BallistaClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// indexed by pod owner
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &k8sapiv1.Pod{}, podOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner...
		pod := rawObj.(*k8sapiv1.Pod)
		owner := metav1.GetControllerOf(pod)
		if owner == nil {
			return nil
		}
		// ...make sure it's a ballista cluster...
		if owner.APIVersion != apiGVStr || owner.Kind != "BallistaCluster" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	// indexed by pod role
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &k8sapiv1.Pod{}, podOwnerKey, func(rawObj client.Object) []string {
		pod := rawObj.(*k8sapiv1.Pod)
		if podRole, ok := pod.Labels[podBallistaRoleKey]; ok {
			return []string{podRole}
		}
		return nil
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.BallistaCluster{}).
		Owns(&k8sapiv1.Pod{}).
		Complete(r)
}

func (c *BallistaClusterReconciler) validateBallistaCluster(cluster *v1.BallistaCluster) error {
	return nil
}

func (r *BallistaClusterReconciler) startBallistaCluster(ctx context.Context, cluster *v1.BallistaCluster)  *v1.BallistaCluster {
	log := log.FromContext(ctx)
	clusterID := uuid.New().String()

	schedulerName := fmt.Sprint("%s-%s", clusterID, "scheduler")
	schedulerPod := &k8sapiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        schedulerName,
			Namespace:   cluster.Namespace,
		},
		Spec: *cluster.Spec.Scheduler.PodSpec.DeepCopy(),
	}

	// ...and create it on the cluster
	if err := r.Create(ctx, schedulerPod); err != nil {
		log.Error(err, "unable to create scheduler pod for Ballista Cluster", "scheduler", schedulerPod)
		// return ctrl.Result{}, err
	}


	return cluster
}

func (r *BallistaClusterReconciler)  handleBallistaClusterDeletion(ctx context.Context, cluster *v1.BallistaCluster)  {
	log := log.FromContext(ctx)
	// BallistaCluster deletion requested, lets delete scheduler pod and executor pods
	if err := r.deleteBallistaResources(ctx, cluster); err != nil {
		glog.Errorf("failed to delete resources associated with deleted SparkApplication %s/%s: %v", app.Namespace, app.Name, err)
	}
}

func (r *BallistaClusterReconciler)  deleteBallistaResources(ctx context.Context, cluster *v1.BallistaCluster) error {
	log := log.FromContext(ctx)

	var childPods = &k8sapiv1.PodList{}
	if err := r.List(ctx, childPods, client.InNamespace(req.Namespace), client.MatchingFields{podOwnerKey: ""}); err != nil {
		log.Error(err, "unable to list child Jobs")
		return ctrl.Result{}, err
	}


	r.Delete()
}

func (r *BallistaClusterReconciler) getAndUpdateClusterState(ctx context.Context, cluster *v1.BallistaCluster) error {
	if err := r.getAndUpdateSchedulerState(ctx, cluster); err != nil {
		return err
	}
	if err := r.getAndUpdateExecutorState(ctx, cluster); err != nil {
		return err
	}
	return nil
}


func (r *BallistaClusterReconciler) getAndUpdateSchedulerState(ctx context.Context, cluster *v1.BallistaCluster) error {
	var pod = &k8sapiv1.PodSpec{}
	r.List(ctx, pod )
	return nil
}

func (r *BallistaClusterReconciler) getAndUpdateExecutorState(ctx context.Context, cluster *v1.BallistaCluster) error {
	return nil
}

func (r *BallistaClusterReconciler) schedulerStateToClusterState()  {

}






