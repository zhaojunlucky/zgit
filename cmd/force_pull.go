// Package cmd /*
package cmd

import (
	"fmt"
	"zhaojunlucky/zgit/core"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// forcePullCmd represents the forcePull command
var forcePullCmd = &cobra.Command{
	Use:   "force-pull",
	Short: "Force pull by deleting local branch and checking out from origin",
	Long: `Force pull is useful when the remote branch has been force-pushed.
It deletes the current local branch, fetches from origin, and checks out 
the branch again from the remote, effectively syncing with the force-pushed changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceBranch, _ := cmd.Flags().GetString("branch")

		// Get current branch name
		currentBranch, err := core.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}

		logrus.Infof("Current branch: %s", currentBranch)

		// Checkout to source branch
		if err := core.RunGitCommand("checkout", sourceBranch); err != nil {
			return fmt.Errorf("failed to checkout to %s: %w", sourceBranch, err)
		}

		// Delete the local branch
		if err := core.RunGitCommand("branch", "-D", currentBranch); err != nil {
			return fmt.Errorf("failed to delete branch %s: %w", currentBranch, err)
		}
		logrus.Infof("Deleted local branch: %s", currentBranch)

		// Fetch from origin
		if err := core.RunGitCommand("fetch", "origin", currentBranch); err != nil {
			return fmt.Errorf("failed to fetch from origin: %w", err)
		}
		logrus.Info("Fetched from origin")

		// Checkout the branch from origin
		if err := core.RunGitCommand("checkout", "-b", currentBranch, fmt.Sprintf("origin/%s", currentBranch)); err != nil {
			return fmt.Errorf("failed to checkout branch %s from origin: %w", currentBranch, err)
		}
		logrus.Infof("Successfully force-pulled branch: %s", currentBranch)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(forcePullCmd)
	forcePullCmd.Flags().StringP("branch", "b", "main", "The source branch to checkout before deleting current branch")
}
