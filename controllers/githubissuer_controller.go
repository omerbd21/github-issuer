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
	"errors"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	githubv1 "github.com/github-issuer/api/v1"
	"github.com/github-issuer/pkg/githubUtils"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GithubIssuerReconciler reconciles a GithubIssuer object
type GithubIssuerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *GithubIssuerReconciler) updateConditions(ctx context.Context, githubIssuer *githubv1.GithubIssuer, conditionType string, reason string, msg string, status metav1.ConditionStatus) error {
	condition := metav1.Condition{Type: conditionType, Status: status, Reason: reason, Message: msg, LastTransitionTime: metav1.Time{Time: time.Now()}}
	githubIssuer.Status.Conditions = append(githubIssuer.Status.Conditions, condition)
	return r.Client.Status().Update(ctx, githubIssuer)

}

/*func getConditionParameters(action string, success bool) {
	conditionParameters := {

	}
}*/

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
	log := logrus.New()
	log.SetFormatter(&ecslogrus.Formatter{})

	var githubIssuer githubv1.GithubIssuer
	if err := r.Get(ctx, req.NamespacedName, &githubIssuer); err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.WithError(errors.New(err.Error())).WithFields(logrus.Fields{"GitHubIssuer": req.NamespacedName}).Error("unable to fetch GithubIssuer")
		return ctrl.Result{Requeue: true}, client.IgnoreNotFound(err)
	}
	githubClient, err := githubUtils.CreateClient(ctx)
	if err != nil {
		log.WithError(errors.New(err.Error())).WithFields(logrus.Fields{"GitHubIssuer": req.NamespacedName}).Error("unable to create GitHub Client")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	_, err = githubUtils.FetchIssue(githubIssuer.Spec.Repo, githubIssuer.Spec.Title, ctx, githubClient)
	if err != nil && strings.Contains(err.Error(), "The issue wasn't found") {
		err = githubUtils.CreateIssue(githubIssuer.Spec.Repo, githubIssuer.Spec.Title, githubIssuer.Spec.Description, ctx, githubClient)
		if err != nil {
			err = r.updateConditions(ctx, &githubIssuer, "IssueCreated", "IssueCreated", "Issue was created", metav1.ConditionTrue)
			if err != nil {
				log.Info("Github Issuer status update failed.")
			}
		} else {
			err = r.updateConditions(ctx, &githubIssuer, "IssueCreated", "IssueNotCreated", "Issue was not created", metav1.ConditionFalse)
			if err != nil {
				log.Info("Github Issuer status update failed.")
			}
		}
	} else if err != nil {
		log.WithError(errors.New(err.Error())).WithFields(logrus.Fields{"GitHubIssuer": req.NamespacedName, "repo": githubIssuer.Spec.Repo, "issue": githubIssuer.Spec.Title}).Error("unable to fetch the specific Issue in Repo")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} else {
		if err := githubUtils.UpdateIssue(githubIssuer.Spec.Repo, githubIssuer.Spec.Title, githubIssuer.Spec.Description, ctx, githubClient); err != nil {
			err = r.updateConditions(ctx, &githubIssuer, "IssueUpdated", "IssueUpdated", "Issue was updated", metav1.ConditionTrue)
			if err != nil {
				log.Info("Github Issuer status update failed.")
			}
			log.WithError(errors.New(err.Error())).WithFields(logrus.Fields{"GitHubIssuer": req.NamespacedName, "repo": githubIssuer.Spec.Repo, "issue": githubIssuer.Spec.Title}).Error("unable to update the issue")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		err = r.updateConditions(ctx, &githubIssuer, "IssueUpdated", "IssueNotUpdated", "Issue was not updated", metav1.ConditionTrue)
		if err != nil {
			log.Info(err.Error())
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GithubIssuerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&githubv1.GithubIssuer{}).
		Complete(r)
}
