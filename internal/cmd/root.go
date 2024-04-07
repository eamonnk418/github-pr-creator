package cmd

import (
	"context"
	githubclient "github-pr-creator/internal/client/github"
	"os"

	"github.com/spf13/cobra"
)

const (
	owner      = "eamonnk418"
	repo       = "github-metrics"
	branchName = "dependabot-config"
)

func NewRootCmd() *cobra.Command {
	opts := struct {
		filepaths []string
		author    string
		email     string
		repo      string
	}{}

	rootCmd := &cobra.Command{
		Use:   "github-pr-creator",
		Short: "github-pr-creator",
		RunE: func(cmd *cobra.Command, args []string) error {
			githubClient := githubclient.New(os.Getenv("GITHUB_TOKEN"))

			defaultBranch, featureBranch, err := githubClient.CheckoutBranch(context.TODO(), owner, repo, branchName)
			if err != nil {
				return err
			}

			tree, err := githubClient.GitAdd(context.TODO(), owner, repo, featureBranch, opts.filepaths)
			if err != nil {
				return err
			}

			if err := githubClient.Commit(context.TODO(), owner, repo, opts.author, opts.email, featureBranch, tree); err != nil {
				return err
			}

			if err := githubClient.Push(context.TODO(), owner, opts.repo, "feat(dependabot): add config file to repo", featureBranch.GetRef(), branchName, defaultBranch.GetRef(), "This change adds the dependabot config to the repo"); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.AddCommand(NewGenerateCmd())

	rootCmd.Flags().StringSliceVarP(&opts.filepaths, "filepaths", "f", []string{".github/dependabot.yml"}, "The filepaths to each file you want to add to the git tree.")
	rootCmd.Flags().StringVarP(&opts.author, "author", "a", "github-actions[bot]", "Author for the commit")
	rootCmd.Flags().StringVarP(&opts.email, "email", "e", "160419984+github-actions[bot]@users.noreply.github.com", "Email for the Git Author of the commit")
	rootCmd.Flags().StringVarP(&opts.repo, "repo", "r", "", "GitHub Repository to create pr in")
	
	
	return rootCmd
}

func Execute() {
	cobra.CheckErr(NewRootCmd().Execute())
}
