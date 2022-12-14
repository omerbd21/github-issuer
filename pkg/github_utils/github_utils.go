package github_utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func divideUserAndRepo(repo string) map[string]string {
	split := strings.Split(repo, "/")

	return map[string]string{
		"user": split[len(split)-2],
		"repo": split[len(split)-1],
	}
}

func CreateClient(ctx context.Context, token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}

func FetchIssue(repo string, issueTitle string, ctx context.Context, client *github.Client) (*github.Issue, error) {
	githubAuth := divideUserAndRepo(repo)
	opts := github.IssueListByRepoOptions{}
	issues, _, err := client.Issues.ListByRepo(ctx, githubAuth["user"], githubAuth["repo"], &opts)
	if err != nil {
		return &github.Issue{}, err
	}
	for _, issue := range issues {
		if *issue.Title == issueTitle {
			return issue, nil
		}
	}
	return &github.Issue{}, fmt.Errorf("%v", "The issue wasn't found")
}

func CreateIssue(repo string, issueTitle string, description string, ctx context.Context, client *github.Client) error {
	githubAuth := divideUserAndRepo(repo)
	req := github.IssueRequest{
		Title: &issueTitle,
		Body:  &description,
	}
	_, _, err := client.Issues.Create(ctx, githubAuth["user"], githubAuth["repo"], &req)
	return err
}

func UpdateIssue(repo string, issueTitle string, description string, ctx context.Context, client *github.Client) error {
	githubAuth := divideUserAndRepo(repo)
	issue, err := FetchIssue(repo, issueTitle, ctx, client)
	if err != nil {
		return err
	}
	req := github.IssueRequest{
		Title: issue.Title,
		Body:  &description,
	}
	if description != *issue.Body {
		_, _, err = client.Issues.Edit(ctx, githubAuth["user"], githubAuth["repo"], *issue.Number, &req)
		return err
	}
	return err
}

func DeleteIssue(repo string, issueTitle string, ctx context.Context, client *github.Client) error {
	githubAuth := divideUserAndRepo(repo)
	state := "closed"
	issue, err := FetchIssue(repo, issueTitle, ctx, client)
	if err != nil {
		return err
	}
	req := github.IssueRequest{
		Title: issue.Title,
		Body:  issue.Body,
		State: &state,
	}
	_, _, err = client.Issues.Edit(ctx, githubAuth["user"], githubAuth["repo"], *issue.Number, &req)
	return err
}
