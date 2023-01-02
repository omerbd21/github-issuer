package github_utils

import (
	"context"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Github Utils", func() {
	Context("get the user and the repo from url", func() {
		It("Should return the correct user and repo", func() {
			url := "https://github.com/test-user/test-repo"
			repoConfig := divideUserAndRepo(url)
			Expect(repoConfig["user"] == "test-user" && repoConfig["repo"] == "test-repo").Should(BeTrue())
		})
	})
	Context("crud methods for github_utils", func() {
		It("Should fetch the issue", func() {
			c := setupFakeClient("GET")
			ctx := context.Background()
			_, err := FetchIssue("https://github.com/test-user/test-repo", "test-title", ctx, c)
			Expect(err).Should(BeNil())
		})
		It("Should create the issue", func() {
			c := setupFakeClient("POST")
			ctx := context.Background()
			err := CreateIssue("https://github.com/test-user/test-repo", "test-title", "test-body", ctx, c)
			Expect(err).Should(BeNil())
		})
		It("Should delete the issue", func() {
			c := setupFakeClient("PATCH")
			ctx := context.Background()
			err := UpdateIssue("https://github.com/test-user/test-repo", "test-title", "test-body2", ctx, c)
			Expect(err).Should(BeNil())
		})
		It("Should delete the issue", func() {
			c := setupFakeClient("PATCH")
			ctx := context.Background()
			err := DeleteIssue("https://github.com/test-user/test-repo", "test-title", ctx, c)
			Expect(err).Should(BeNil())
		})
		It("Should return an error for get", func() {
			c := setupFakeClient("GET_ERROR")
			ctx := context.Background()
			_, err := FetchIssue("https://github.com/no-user/no-repo", "no-issue", ctx, c)
			Expect(err).ShouldNot(BeNil())
		})
		It("Should return an error for create", func() {
			c := setupFakeClient("CREATE_ERROR")
			ctx := context.Background()
			err := CreateIssue("https://github.com/no-user/no-repo", "no-issue", "no-body", ctx, c)
			Expect(err).ShouldNot(BeNil())
		})
		It("Should return an error for update", func() {
			c := setupFakeClient("UPDATE_ERROR")
			ctx := context.Background()
			err := UpdateIssue("https://github.com/no-user/no-repo", "no-issue", "no-body", ctx, c)
			Expect(err).ShouldNot(BeNil())
		})

	})
})

func setupFakeClient(method string) *github.Client {
	mockedHTTPClient := mock.NewMockedHTTPClient()
	if method == "PATCH" {
		mockedHTTPClient = mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposIssuesByOwnerByRepo,
				[]github.Issue{
					{
						Title:  github.String("test-title"),
						Body:   github.String("test-body"),
						Number: github.Int(1),
						Repository: &github.Repository{
							Name: github.String("test-repo"),
							Owner: &github.User{
								Name: github.String("test-user"),
							},
						},
					},
				},
			),
			mock.WithRequestMatch(
				mock.PatchReposIssuesByOwnerByRepoByIssueNumber,
				github.Issue{
					Title:  github.String("test-title"),
					Body:   github.String("test-body"),
					Number: github.Int(1),
					Repository: &github.Repository{
						Name: github.String("test-repo"),
						Owner: &github.User{
							Name: github.String("test-user"),
						},
					},
				},
			),
		)
	} else if method == "GET" {
		mockedHTTPClient = mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposIssuesByOwnerByRepo,
				[]github.Issue{
					{
						Title:  github.String("test-title"),
						Body:   github.String("test-body"),
						Number: github.Int(1),
						Repository: &github.Repository{
							Name: github.String("test-repo"),
							Owner: &github.User{
								Name: github.String("test-user"),
							},
						},
					},
				},
			),
		)
	} else if method == "POST" {
		mockedHTTPClient = mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.PostReposIssuesByOwnerByRepo,
				github.Issue{
					Title: github.String("test-title"),
					Body:  github.String("test-body"),
					Repository: &github.Repository{
						Name: github.String("test-repo"),
						Owner: &github.User{
							Name: github.String("test-user"),
						},
					},
				},
			))
	} else if method == "GET_ERROR" {
		mockedHTTPClient = mock.NewMockedHTTPClient(
			mock.WithRequestMatchHandler(
				mock.GetReposIssuesByOwnerByRepo,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						"The issue can't be created",
					)
				}),
			),
		)
	} else if method == "CREATE_ERROR" {
		mockedHTTPClient = mock.NewMockedHTTPClient(
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesByOwnerByRepo,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						"The issue can't be created",
					)
				}),
			),
		)
	} else if method == "UPDATE_ERROR" {
		mockedHTTPClient = mock.NewMockedHTTPClient(
			mock.WithRequestMatchHandler(
				mock.PatchReposIssuesByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						"The issue can't be created",
					)
				}),
			),
		)
	}
	c := github.NewClient(mockedHTTPClient)
	return c
}
