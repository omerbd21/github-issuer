package github_utils

import (
	"context"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	REGULAR_URL       = "https://github.com/test-user/test-repo"
	ERROR_URL         = "https://github.com/no-user/no-repo"
	USER              = "test-user"
	REPO              = "test-repo"
	ISSUE             = "test-title"
	ERROR_ISSUE       = "no-title"
	DESCRIPTION       = "test-body"
	ERROR_DESCRIPTION = "no-body"
	NUMBER            = 1
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
			_, err := FetchIssue(REGULAR_URL, ISSUE, ctx, c)
			Expect(err).Should(BeNil())
		})
		It("Should create the issue", func() {
			c := setupFakeClient("POST")
			ctx := context.Background()
			err := CreateIssue(REGULAR_URL, ISSUE, DESCRIPTION, ctx, c)
			Expect(err).Should(BeNil())
		})
		It("Should delete the issue", func() {
			c := setupFakeClient("PATCH")
			ctx := context.Background()
			err := UpdateIssue(REGULAR_URL, ISSUE, DESCRIPTION, ctx, c)
			Expect(err).Should(BeNil())
		})
		It("Should delete the issue", func() {
			c := setupFakeClient("PATCH")
			ctx := context.Background()
			err := DeleteIssue(REGULAR_URL, ISSUE, ctx, c)
			Expect(err).Should(BeNil())
		})
		It("Should return an error for get", func() {
			c := setupFakeClient("GET_ERROR")
			ctx := context.Background()
			_, err := FetchIssue(ERROR_URL, ERROR_ISSUE, ctx, c)
			Expect(err).ShouldNot(BeNil())
		})
		It("Should return an error for create", func() {
			c := setupFakeClient("CREATE_ERROR")
			ctx := context.Background()
			err := CreateIssue(ERROR_URL, ERROR_ISSUE, ERROR_DESCRIPTION, ctx, c)
			Expect(err).ShouldNot(BeNil())
		})
		It("Should return an error for update", func() {
			c := setupFakeClient("UPDATE_ERROR")
			ctx := context.Background()
			err := UpdateIssue(ERROR_URL, ERROR_ISSUE, ERROR_DESCRIPTION, ctx, c)
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
				github.Issue{
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
			),
		)
	} else if method == "GET" {
		mockedHTTPClient = mock.NewMockedHTTPClient(
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
		)
	} else if method == "POST" {
		mockedHTTPClient = mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.PostReposIssuesByOwnerByRepo,
				github.Issue{
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
			))
	} else if method == "GET_ERROR" {
		mockedHTTPClient = mock.NewMockedHTTPClient(
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
						"The issue can't be updated",
					)
				}),
			),
		)
	}
	c := github.NewClient(mockedHTTPClient)
	return c
}
