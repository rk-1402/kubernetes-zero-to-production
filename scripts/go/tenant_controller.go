// Capstone 9 — Custom Operator with Kubebuilder
// Tenant reconciler: when a Tenant CR is applied, provision a Namespace +
// ResourceQuota + default-deny NetworkPolicy + admin RoleBinding.
// File: internal/controller/tenant_controller.go
package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	platformv1 "github.com/acme/tenant-op/api/v1"
)

type TenantReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *TenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var t platformv1.Tenant
	if err := r.Get(ctx, req.NamespacedName, &t); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	nsName := "team-" + t.Spec.Team

	// 1. Namespace
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: nsName, Labels: map[string]string{"tenant": t.Name}}}
	_ = ctrl.SetControllerReference(&t, ns, r.Scheme)
	if err := r.Client.Patch(ctx, ns, client.Apply, client.ForceOwnership, client.FieldOwner("tenant-op")); err != nil {
		return ctrl.Result{}, err
	}

	// 2. ResourceQuota
	rq := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{Name: "default", Namespace: nsName},
		Spec: corev1.ResourceQuotaSpec{Hard: corev1.ResourceList{
			corev1.ResourceRequestsCPU:    resource.MustParse(t.Spec.CPU),
			corev1.ResourceRequestsMemory: resource.MustParse(t.Spec.Memory),
		}},
	}
	_ = ctrl.SetControllerReference(&t, rq, r.Scheme)
	if err := r.Client.Patch(ctx, rq, client.Apply, client.ForceOwnership, client.FieldOwner("tenant-op")); err != nil {
		return ctrl.Result{}, err
	}

	// 3. NetworkPolicy default-deny
	np := &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "default-deny", Namespace: nsName},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []netv1.PolicyType{netv1.PolicyTypeIngress, netv1.PolicyTypeEgress},
		},
	}
	_ = ctrl.SetControllerReference(&t, np, r.Scheme)
	_ = r.Client.Patch(ctx, np, client.Apply, client.ForceOwnership, client.FieldOwner("tenant-op"))

	// 4. RoleBinding for admins
	if len(t.Spec.Admins) > 0 {
		rb := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{Name: "team-admin", Namespace: nsName},
			RoleRef:    rbacv1.RoleRef{Kind: "ClusterRole", Name: "admin", APIGroup: "rbac.authorization.k8s.io"},
		}
		for _, u := range t.Spec.Admins {
			rb.Subjects = append(rb.Subjects, rbacv1.Subject{Kind: "User", Name: u, APIGroup: "rbac.authorization.k8s.io"})
		}
		_ = ctrl.SetControllerReference(&t, rb, r.Scheme)
		_ = r.Client.Patch(ctx, rb, client.Apply, client.ForceOwnership, client.FieldOwner("tenant-op"))
	}

	t.Status.Namespace = nsName
	t.Status.Ready = true
	return ctrl.Result{}, r.Status().Update(ctx, &t)
}

// TenantSpec / TenantStatus (api/v1/tenant_types.go)
//
// type TenantSpec struct {
//     Team   string   `json:"team"`
//     CPU    string   `json:"cpu"`
//     Memory string   `json:"memory"`
//     Admins []string `json:"admins"`
// }
// type TenantStatus struct {
//     Namespace string `json:"namespace,omitempty"`
//     Ready     bool   `json:"ready,omitempty"`
// }
