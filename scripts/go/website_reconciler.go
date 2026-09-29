// Topic 58 — Operators and the Operator Pattern
// Minimal Kubebuilder reconciler: keep a Deployment matching a Website CR,
// then write status back. File: internal/controller/website_controller.go
package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	webv1 "github.com/me/website-op/api/v1"
)

type WebsiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *WebsiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var site webv1.Website
	if err := r.Get(ctx, req.NamespacedName, &site); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Desired Deployment
	dep := buildDeployment(&site)
	if err := ctrl.SetControllerReference(&site, dep, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Create or update
	var existing appsv1.Deployment
	if err := r.Get(ctx, client.ObjectKeyFromObject(dep), &existing); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, r.Create(ctx, dep)
		}
		return ctrl.Result{}, err
	}
	existing.Spec = dep.Spec
	if err := r.Update(ctx, &existing); err != nil {
		return ctrl.Result{}, err
	}

	// Update status
	site.Status.URL = fmt.Sprintf("http://%s.svc", site.Name)
	return ctrl.Result{RequeueAfter: 30 * time.Second}, r.Status().Update(ctx, &site)
}

// buildDeployment converts a Website CR into the desired Deployment (implement
// the pod template from site.Spec.GitRepo / site.Spec.Replicas).
func buildDeployment(site *webv1.Website) *appsv1.Deployment {
	// ... construct and return the *appsv1.Deployment ...
	return &appsv1.Deployment{}
}
