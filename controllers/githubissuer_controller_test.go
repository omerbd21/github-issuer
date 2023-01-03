package controllers

import (
	"context"
	"fmt"
	"net/http"

	githubv1 "github.com/github-issuer/api/v1"
	"github.com/github-issuer/pkg/github_utils"
	"github.com/google/go-github/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("GithubIssuer controller", func() {
	Context("GithubIssuer controller test", func() {

		const GithubIssuerName = "test-githubissuer"

		ctx := context.Background()
		namespace := &corev1.Namespace{}
		testCounter := 0
		typeNamespaceName := types.NamespacedName{Name: GithubIssuerName, Namespace: GithubIssuerName}

		BeforeEach(func() {
			testCounter++
			// Each test case gets its own namespace
			namespace = &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: GithubIssuerName + fmt.Sprint(testCounter),
				},
			}
			err := k8sClient.Create(ctx, namespace)
			Expect(err).To(Not(HaveOccurred()))
			typeNamespaceName.Namespace = GithubIssuerName + fmt.Sprint(testCounter)
		})

		AfterEach(func() {
			/*err := k8sClient.Delete(ctx, namespace)
			Expect(err).To(Not(HaveOccurred()))*/
		})

		It("should successfully get an issue", func() {
			By("Creating the custom resource for the Kind GithubIssuer")
			githubIssuer := &githubv1.GithubIssuer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      typeNamespaceName.Name,
					Namespace: typeNamespaceName.Namespace,
				},
				Spec: githubv1.GithubIssuerSpec{
					Repo:        "https://github.com/test-user/test-repo",
					Title:       "test-title",
					Description: "test-body",
				},
			}
			err := k8sClient.Create(ctx, githubIssuer)
			Expect(err).Should(BeNil())
			By("Checking the issue exists")
			mockedHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatch(
					mock.GetReposIssuesByOwnerByRepo,
					[]github.Issue{
						{
							Title:  github.String(ISSUE),
							Body:   github.String(DESCRIPTION),
							Number: github.Int(NUMBER),
							Repository: &github.Repository{
								Name: github.String(REPO),
								Owner: &github.User{
									Name: github.String(USER),
								},
							},
						},
					},
				),
				mock.WithRequestMatch(
					mock.PatchReposIssuesByOwnerByRepoByIssueNumber,
					github.Issue{},
				),
			)
			nclient := github.NewClient(mockedHTTPClient)
			issue, err := github_utils.FetchIssue(githubIssuer.Spec.Repo, githubIssuer.Spec.Title, ctx, nclient)
			Expect(issue != nil && err == nil).Should(BeTrue())

		})
		It("should successfully delete an issue", func() {
			By("Creating the custom resource for the Kind GithubIssuer")
			githubIssuer := &githubv1.GithubIssuer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      typeNamespaceName.Name,
					Namespace: typeNamespaceName.Namespace,
				},
				Spec: githubv1.GithubIssuerSpec{
					Repo:        "https://github.com/test-user/test-repo",
					Title:       "test-title",
					Description: "test-body",
				},
			}
			err := k8sClient.Create(ctx, githubIssuer)
			Expect(err).Should(BeNil())
			By("Deleting the custom resource for the Kind GithubIssuer")
			err = k8sClient.Delete(ctx, githubIssuer)
			Expect(err).Should(BeNil())
			By("Checking the issue doesn't exist")
			mockedHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatchHandler(
					mock.GetReposIssuesByOwnerByRepo,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						mock.WriteError(
							w,
							http.StatusInternalServerError,
							"The issue can't be fetched",
						)
					}),
				),
			)
			nclient := github.NewClient(mockedHTTPClient)
			_, err = github_utils.FetchIssue(githubIssuer.Spec.Repo, githubIssuer.Spec.Title, ctx, nclient)
			Expect(err != nil).Should(BeTrue())

		})
		It("should successfully Create an issue", func() {
			By("Creating the custom resource for the Kind GithubIssuer")
			githubIssuer := &githubv1.GithubIssuer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      typeNamespaceName.Name,
					Namespace: typeNamespaceName.Namespace,
				},
				Spec: githubv1.GithubIssuerSpec{
					Repo:        "https://github.com/test-user/test-repo",
					Title:       "test-title",
					Description: "test-body",
				},
			}
			err := k8sClient.Create(ctx, githubIssuer)
			Expect(err).Should(BeNil())
			By("Checking the issue exists")
			mockedHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatch(
					mock.GetReposIssuesByOwnerByRepo,
					[]github.Issue{
						{
							Title:  github.String(ISSUE),
							Body:   github.String(DESCRIPTION),
							Number: github.Int(NUMBER),
							Repository: &github.Repository{
								Name: github.String(REPO),
								Owner: &github.User{
									Name: github.String(USER),
								},
							},
						},
					},
				),
				mock.WithRequestMatch(
					mock.PatchReposIssuesByOwnerByRepoByIssueNumber,
					github.Issue{},
				),
			)
			nclient := github.NewClient(mockedHTTPClient)
			issue, err := github_utils.FetchIssue(githubIssuer.Spec.Repo, githubIssuer.Spec.Title, ctx, nclient)
			Expect(issue != nil && err == nil).Should(BeTrue())

		})
		It("should successfully Update an issue", func() {
			By("Creating the custom resource for the Kind GithubIssuer")
			githubIssuer := &githubv1.GithubIssuer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      typeNamespaceName.Name,
					Namespace: typeNamespaceName.Namespace,
				},
				Spec: githubv1.GithubIssuerSpec{
					Repo:        "https://github.com/test-user/test-repo",
					Title:       "test-title",
					Description: "test-body",
				},
			}
			err := k8sClient.Create(ctx, githubIssuer)
			Expect(err).Should(BeNil())
			By("Updating the custom resource for the Kind GithubIssuer")
			githubIssuer = &githubv1.GithubIssuer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      typeNamespaceName.Name + "2",
					Namespace: typeNamespaceName.Namespace,
				},
				Spec: githubv1.GithubIssuerSpec{
					Repo:        "https://github.com/test-user/test-repo",
					Title:       "test-title",
					Description: "test-body2",
				},
			}
			err = k8sClient.Create(ctx, githubIssuer)
			Expect(err).Should(BeNil())
			By("Checking the issue exists and was changed")
			mockedHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatch(
					mock.GetReposIssuesByOwnerByRepo,
					[]github.Issue{
						{
							Title:  github.String(ISSUE),
							Body:   github.String(DESCRIPTION + "2"),
							Number: github.Int(NUMBER),
							Repository: &github.Repository{
								Name: github.String(REPO),
								Owner: &github.User{
									Name: github.String(USER),
								},
							},
						},
					},
				),
				mock.WithRequestMatch(
					mock.PatchReposIssuesByOwnerByRepoByIssueNumber,
					github.Issue{},
				),
			)
			nclient := github.NewClient(mockedHTTPClient)
			issue, _ := github_utils.FetchIssue(githubIssuer.Spec.Repo, githubIssuer.Spec.Title, ctx, nclient)
			Expect(*issue.Body == "test-body2").Should(BeTrue())

		})

	})
})
