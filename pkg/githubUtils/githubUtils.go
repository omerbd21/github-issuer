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

func FetchIssue(user string, repo string, issueTitle string, ctx context.Context, client *github.Client) (*github.Issue, error) {
	opts := github.IssueListByRepoOptions{}
	issues, _, err := client.Issues.ListByRepo(ctx, user, repo, &opts)
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

func CreateIssue(user string, repo string, issueTitle string, description string, ctx context.Context, client *github.Client) error {
	req := github.IssueRequest{
		Title: &issueTitle,
		Body:  &description,
	}
	_, _, err := client.Issues.Create(ctx, user, repo, &req)
	return err
}

func UpdateIssue(user string, repo string, issueTitle string, description string, ctx context.Context, client *github.Client) error {
	issue, err := FetchIssue(user, repo, issueTitle, ctx, client)
	if err != nil {
		return err
	}
	req := github.IssueRequest{
		Title: &issueTitle,
		Body:  &description,
	}
	if description != *issue.Body {
		_, _, err = client.Issues.Edit(ctx, user, repo, *issue.Number, &req)
		return err
	}
	return err
}
