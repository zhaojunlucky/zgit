/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var repoDir string



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zgit",
	Short: "Enhanced git workflow tool with automatic ticket tracking",
	Long: `zgit is a Git workflow enhancement tool that automates common tasks.

Features:
  - Automatic ticket extraction from branch names
  - Template-based commit message formatting
  - Force-pull for syncing after force-pushed branches
  - Repository-specific and global configuration support

Configuration:
  zgit looks for config.yaml in the current directory or ~/.zgit/config.yaml
  
  Example config:
    global:
      branches:
        - usr/name/(?P<ticket>JIRA-\d+)
      commit:
        message: "[{{.Ticket}}] {{.Message}}"

Commands:
  commit      - Commit with automatic ticket prefix
  force-pull  - Force pull by recreating local branch from origin`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Change to repo directory if specified
		if repoDir != "" {
			if err := os.Chdir(repoDir); err != nil {
				log.Fatalf("failed to change to directory %s: %v", repoDir, err)
			}
			log.Infof("changed to directory: %s", repoDir)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global persistent flags available to all subcommands
	rootCmd.PersistentFlags().StringVarP(&repoDir, "repo-dir", "C", "", "Git repository directory (default is current directory)")
}


