/*
Copyright 2022.

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
	"os"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	githubv1 "github.com/github-issuer/api/v1"
	"github.com/github-issuer/pkg/githubUtils"
	"github.com/go-logr/logr"
	"github.com/google/go-github/github"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// GithubIssuerReconciler reconciles a GithubIssuer object
type GithubIssuerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Secret string
}

const FinalizerName = "github.benda.io/finalizer"

func (r *GithubIssuerReconciler) updateConditions(ctx context.Context, githubIssuer *githubv1.GithubIssuer, conditionType string, reason string, msg string, status metav1.ConditionStatus) error {
	condition := metav1.Condition{Type: conditionType, Status: status, Reason: reason, Message: msg, LastTransitionTime: metav1.Time{Time: time.Now()}}
	githubIssuer.Status.Conditions = append(githubIssuer.Status.Conditions, condition)
	return r.Client.Status().Update(ctx, githubIssuer)

}

//+kubebuilder:rbac:groups=github.benda.io,resources=githubissuers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=github.benda.io,resources=githubissuers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=github.benda.io,resources=githubissuers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GithubIssuer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *GithubIssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)

	var githubIssuer githubv1.GithubIssuer
	if err := r.Get(ctx, req.NamespacedName, &githubIssuer); err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Unable to fetch GithubIssuer", "githubIssuer", req.NamespacedName.String())
		return ctrl.Result{Requeue: true}, client.IgnoreNotFound(err)
	}
	githubClient, err := githubUtils.CreateClient(ctx, os.Getenv("GITHUB_PASSWORD"))
	if err != nil {
		log.Error(err, "Unable to create GitHub Client", "githubIssuer", req.NamespacedName.String())
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if githubIssuer.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&githubIssuer, FinalizerName) {
			if err := r.addFinalizer(ctx, log, &githubIssuer); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(&githubIssuer, FinalizerName) {
			if res, err := r.deleteIssue(ctx, log, &githubIssuer, githubClient); err != nil {
				return res, err
			}
			return ctrl.Result{}, nil
		}
	}
	issue, err := githubUtils.FetchIssue(githubIssuer.Spec.Repo, githubIssuer.Spec.Title, ctx, githubClient)
	if err != nil && strings.Contains(err.Error(), "The issue wasn't found") {
		err = githubUtils.CreateIssue(githubIssuer.Spec.Repo, githubIssuer.Spec.Title, githubIssuer.Spec.Description, ctx, githubClient)
		if err != nil {
			err = r.updateConditions(ctx, &githubIssuer, "IssueCreated", "IssueCreated", "Issue was created", metav1.ConditionTrue)
			if err != nil {
				log.Error(err, "Unable to update githubIssuer status", "githubIssuer", req.NamespacedName.String(), "issue", issue)
			}
		} else {
			err = r.updateConditions(ctx, &githubIssuer, "IssueNotCreated", "IssueNotCreated", "Issue was not created", metav1.ConditionFalse)
			if err != nil {
				log.Error(err, "Unable to update githubIssuer status", "githubIssuer", req.NamespacedName.String(), "issue", issue)
			}
		}
	} else if err != nil {
		log.Error(err, "Unable to fetch the specific issue in repo", "githubIssuer", req.NamespacedName.String(), "repo", githubIssuer.Spec.Repo, "issue", issue)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} else {
		if err := githubUtils.UpdateIssue(githubIssuer.Spec.Repo, githubIssuer.Spec.Title, githubIssuer.Spec.Description, ctx, githubClient); err != nil {
			err = r.updateConditions(ctx, &githubIssuer, "IssueUpdated", "IssueUpdated", "Issue was updated, issue status: "+*issue.State, metav1.ConditionTrue)
			if err != nil {
				log.Error(err, "Unable to update githubIssuer status", "githubIssuer", req.NamespacedName.String(), "issue", issue)
				err = r.updateConditions(ctx, &githubIssuer, "IssueNotUpdated", "IssueNotUpdated", "Issue was not updated", metav1.ConditionTrue)
				if err != nil {
					log.Info(err.Error())
				}
				log.Error(err, "Unable to update the issue", "githubIssuer", req.NamespacedName.String(), "repo", githubIssuer.Spec.Repo, "issue", issue)
			}
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	}

	return ctrl.Result{}, nil
}

func (r *GithubIssuerReconciler) deleteIssue(ctx context.Context, log logr.Logger, githubIssuer *githubv1.GithubIssuer, githubClient *github.Client) (ctrl.Result, error) {
	repo := githubIssuer.Spec.Repo
	issue := githubIssuer.Spec.Title
	if err := githubUtils.DeleteIssue(repo, issue, ctx, githubClient); err != nil {
		log.Error(err, "unable to delete issue from github", "githubIssuer", githubIssuer.Name, "issue", issue)
		return ctrl.Result{Requeue: true}, err
	}
	controllerutil.RemoveFinalizer(githubIssuer, FinalizerName)
	if err := r.Update(ctx, githubIssuer); err != nil {
		log.Error(err, "unable to remove finalizer from githubissuer", "githubIssuer", githubIssuer.Name)
		return ctrl.Result{Requeue: true}, err
	}
	log.Info("issue was deleted", "githubIssuer", githubIssuer.Name)
	return ctrl.Result{}, nil
}

func (r *GithubIssuerReconciler) addFinalizer(ctx context.Context, log logr.Logger, githubIssuer *githubv1.GithubIssuer) error {
	controllerutil.AddFinalizer(githubIssuer, FinalizerName)
	if err := r.Update(ctx, githubIssuer); err != nil {
		log.Error(err, "unable to add finalizer to githubIssuer", "githubIssuer", githubIssuer.Name)
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GithubIssuerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&githubv1.GithubIssuer{}).
		Complete(r)
}
