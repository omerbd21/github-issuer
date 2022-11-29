package githubUtils

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func CreateClient(ctx context.Context) (*github.Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}

func FetchIssue(repo string, issueTitle string, ctx context.Context, client *github.Client) (*github.Issue, error) {
	opts := github.IssueListByRepoOptions{}
	issues, _, err := client.Issues.ListByRepo(ctx, os.Getenv("GITHUB_USERNAME"), repo, &opts)
	if err != nil {
		return nil, err
	}
	for _, issue := range issues {
		if *issue.Title == issueTitle {
			return issue, nil
		}
	}
	return nil, fmt.Errorf("%v", "The issue wasn't found")
}

func CreateIssue(repo string, issueTitle string, description string, ctx context.Context, client *github.Client) error {
	req := github.IssueRequest{
		Title: &issueTitle,
		Body:  &description,
	}
	_, _, err := client.Issues.Create(ctx, os.Getenv("GITHUB_USERNAME"), repo, &req)
	return err
}

func UpdateIssue(repo string, issueTitle string, description string, ctx context.Context, client *github.Client) error {
	issue, err := FetchIssue(repo, issueTitle, ctx, client)
	if err != nil {
		return err
	}
	req := github.IssueRequest{
		Title: &issueTitle,
		Body:  &description,
	}
	if description != *issue.Body {
		_, _, err = client.Issues.Edit(ctx, os.Getenv("GITHUB_USERNAME"), repo, *issue.Number, &req)
		return err
	}
	return err
}
