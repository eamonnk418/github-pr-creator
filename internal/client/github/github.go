package githubclient

import (
	"context"
	"fmt"
	"github-pr-creator/internal/config"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v61/github"
)

type Client struct {
	Client *github.Client
}

func New(token string) *Client {
	httpClient := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   10 * time.Second,
	}

	client := github.NewClient(httpClient).WithAuthToken(token)

	return &Client{Client: client}
}

func (c Client) FindPackageJSON(ctx context.Context, owner, repo string) error {
	fileContent, directoryContent, resp, err := c.Client.Repositories.GetContents(ctx, owner, repo, "package.json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("Found file %s\n", fileContent.GetName())

	for _, content := range directoryContent {
		fmt.Printf("Found file %s\n", content.GetName())
	}

	return nil
}

func (c Client) GenerateDependabotConfig(ctx context.Context, owner, repo string) error {
	repository, _, err := c.Client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return err
	}

	// Get the primary language of the repository
	lang := repository.GetLanguage()

	// Get the directories for each ecosystem
	if err := c.FindPackageJSON(ctx, owner, repo); err != nil {
		return err
	}

	dependabot := &config.Dependabot{
		Version: 2,
		Updates: []*config.Update{
			{
				PackageEcosystem: lang,
				Directory:        "/",
				Schedule: &config.Schedule{
					Interval: "weekly",
				},
			},
		},
	}

	fmt.Printf("Creating Dependabot Config for %s which uses %s.\n", owner+"/"+repo, lang)
	fmt.Printf("Dependabot Config: %v\n", dependabot)

	return nil
}

func (c Client) CheckoutBranch(ctx context.Context, owner, repo, branchName string) (*github.Reference, *github.Reference, error) {
	repository, _, err := c.Client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, nil, err
	}

	defaultBranch := repository.GetDefaultBranch()

	baseRef, _, err := c.Client.Git.GetRef(ctx, owner, repo, "refs/heads/"+defaultBranch)
	if err != nil {
		return nil, nil, err
	}

	newRef, _, err := c.Client.Git.CreateRef(ctx, owner, repo, &github.Reference{
		Ref: github.String("refs/heads/" + branchName),
		Object: &github.GitObject{
			SHA: github.String(baseRef.GetObject().GetSHA()),
		},
	})

	return baseRef, newRef, nil
}

func (c Client) GitAdd(ctx context.Context, owner, repo string, ref *github.Reference, files []string) (*github.Tree, error) {
	var entries []*github.TreeEntry

	for _, file := range files {
		fileContents, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		entries = append(entries, &github.TreeEntry{
			Path:    github.String(file),
			Mode:    github.String("100644"),
			Type:    github.String("blob"),
			Content: github.String(string(fileContents)),
		})
	}

	tree, _, err := c.Client.Git.CreateTree(ctx, owner, repo, ref.GetObject().GetSHA(), entries)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func (c Client) Commit(ctx context.Context, owner, repo, authorName, authorEmail string, ref *github.Reference, tree *github.Tree) error {
	parent, _, err := c.Client.Repositories.GetCommit(ctx, owner, repo, ref.GetObject().GetSHA(), nil)
	if err != nil {
		return err
	}

	parent.Commit.SHA = parent.SHA

	date := time.Now()
	author := &github.CommitAuthor{Date: &github.Timestamp{Time: date}, Name: github.String(authorName), Email: github.String(authorEmail)}
	commit := &github.Commit{Author: author, Message: github.String("feat(dependabot): add config to repo"), Tree: tree, Parents: []*github.Commit{parent.Commit}}

	newCommit, _, err := c.Client.Git.CreateCommit(ctx, owner, repo, commit, nil)
	if err != nil {
		return err
	}

	ref.Object.SHA = newCommit.SHA

	if _, _, err := c.Client.Git.UpdateRef(ctx, owner, repo, ref, false); err != nil {
		return err
	}

	return nil
}

func (c Client) Push(ctx context.Context, owner, repo, title, head, headRepo, base, body string) error {
	pullRequest, _, err := c.Client.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title:               github.String(title),
		Head:                github.String(head),
		HeadRepo:            github.String(headRepo),
		Base:                github.String(base),
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Created pull request %s\n", pullRequest.GetHTMLURL())

	return nil
}
