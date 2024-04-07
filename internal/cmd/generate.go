package cmd

import (
	"context"
	githubclient "github-pr-creator/internal/client/github"
	"os"

	"github.com/spf13/cobra"
)

func NewGenerateCmd() *cobra.Command {
	opts := struct {
		repo string
	}{}

	generateCmd := &cobra.Command{
		Use: "generate",
		RunE: func(cmd *cobra.Command, args []string) error {
			githubClient := githubclient.New(os.Getenv("GITHUB_TOKEN"))

			if err := githubClient.FindPackageJSON(context.TODO(), owner, repo); err != nil {
				return err
			}

			return nil
		},
	}

	generateCmd.Flags().StringVarP(&opts.repo, "repo", "r", "", "GitHub Repository to create pr in")

	return generateCmd
}
