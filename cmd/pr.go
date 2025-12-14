/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"zhaojunlucky/zgit/core"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var prRemoteName string
var prBaseBranch string

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Open GitHub pull request creation page",
	Long: `Open the GitHub pull request creation page in your default web browser.

This command opens the compare URL to create a PR from the current branch
to the base branch (default branch by default).

Examples:
  zgit pr                    # Compare current branch with default branch
  zgit pr -b main            # Compare current branch with 'main' branch
  zgit pr -r upstream        # Use 'upstream' remote
  zgit pr -b develop -r fork # Compare with 'develop' on 'fork' remote`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get current branch
		currentBranch, err := core.GetCurrentBranch()
		if err != nil {
			log.Fatalf("failed to get current branch: %v", err)
		}
		log.Infof("current branch: %s", currentBranch)

		// Get remote URL
		remoteURL, err := getRemoteURL(prRemoteName)
		if err != nil {
			log.Fatalf("failed to get remote URL: %v", err)
		}

		// Parse to web URL
		webURL, err := parseGitURLToWeb(remoteURL)
		if err != nil {
			log.Fatalf("failed to parse git URL: %v", err)
		}

		// Get base branch (default branch if not specified)
		baseBranch := prBaseBranch
		if baseBranch == "" {
			baseBranch, err = getDefaultBranch(prRemoteName)
			if err != nil {
				log.Fatalf("failed to get default branch: %v", err)
			}
		}
		log.Infof("base branch: %s", baseBranch)

		// Build PR URL: https://github.com/owner/repo/compare/base...head
		prURL := fmt.Sprintf("%s/compare/%s...%s", webURL, baseBranch, currentBranch)
		log.Infof("opening %s", prURL)

		if err := openBrowser(prURL); err != nil {
			log.Fatalf("failed to open browser: %v", err)
		}
	},
}

// getDefaultBranch gets the default branch for the remote
func getDefaultBranch(remote string) (string, error) {
	// Try to get the default branch from remote HEAD
	cmd := exec.Command("git", "symbolic-ref", fmt.Sprintf("refs/remotes/%s/HEAD", remote))
	output, err := cmd.Output()
	if err == nil {
		// refs/remotes/origin/HEAD -> refs/remotes/origin/main
		ref := strings.TrimSpace(string(output))
		// Extract branch name from refs/remotes/origin/main
		parts := strings.Split(ref, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}

	// Fallback: try common default branch names
	commonDefaults := []string{"main", "master"}
	for _, branch := range commonDefaults {
		cmd := exec.Command("git", "rev-parse", "--verify", fmt.Sprintf("refs/remotes/%s/%s", remote, branch))
		if err := cmd.Run(); err == nil {
			return branch, nil
		}
	}

	return "", fmt.Errorf("could not determine default branch for remote '%s'", remote)
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().StringVarP(&prRemoteName, "remote", "r", "origin", "Remote name (default: origin)")
	prCmd.Flags().StringVarP(&prBaseBranch, "base", "b", "", "Base branch for comparison (default: remote's default branch)")
}
