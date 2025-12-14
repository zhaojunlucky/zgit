/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"zhaojunlucky/zgit/core"

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
  - Pass-through for any other git commands

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
  force-pull  - Force pull by recreating local branch from origin
  init        - Initialize zgit configuration
  version     - Show version information
  
  Any other command will be passed directly to git`,
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true, // Allow unknown flags to pass through to git
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Change to repo directory if specified
		if repoDir != "" {
			if err := os.Chdir(repoDir); err != nil {
				log.Fatalf("failed to change to directory %s: %v", repoDir, err)
			}
			log.Infof("changed to directory: %s", repoDir)
		}
	},
	// Handle unknown subcommands by passing them to git
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no args, show help
		if len(args) == 0 {
			return cmd.Help()
		}
		
		// Check if it's a known subcommand
		for _, c := range cmd.Commands() {
			if c.Name() == args[0] {
				// Let cobra handle it
				return nil
			}
		}
		
		// Unknown command - pass to git
		log.Infof("passing command to git: %v", args)
		gitArgs := args
		if err := core.RunGitCommand(gitArgs...); err != nil {
			log.Fatalf("git command failed: %v", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Check if the first argument is an unknown command
	// If so, pass everything directly to git to avoid flag parsing issues
	if len(os.Args) > 1 {
		// Handle -C or --repo-dir flag
		argIdx := 1
		if os.Args[argIdx] == "-C" || os.Args[argIdx] == "--repo-dir" {
			if len(os.Args) > argIdx+1 {
				repoDir = os.Args[argIdx+1]
				argIdx += 2
			}
		}
		
		if argIdx < len(os.Args) {
			subcommand := os.Args[argIdx]
			
			// Skip if it's a known flag or known command
			if subcommand != "-h" && subcommand != "--help" && !isKnownCommand(subcommand) {
				// Change to repo directory if specified
				if repoDir != "" {
					if err := os.Chdir(repoDir); err != nil {
						log.Fatalf("failed to change to directory %s: %v", repoDir, err)
					}
					log.Infof("changed to directory: %s", repoDir)
				}
				
				// Unknown command - pass everything to git
				gitArgs := os.Args[argIdx:]
				log.Infof("passing command to git: %v", gitArgs)
				if err := core.RunGitCommand(gitArgs...); err != nil {
					log.Fatalf("git command failed: %v", err)
				}
				return
			}
		}
	}
	
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// isKnownCommand checks if a command is a known zgit subcommand
func isKnownCommand(cmd string) bool {
	knownCommands := []string{"commit", "force-pull", "init", "version", "completion", "help", "open", "pr"}
	for _, known := range knownCommands {
		if cmd == known {
			return true
		}
	}
	return false
}

func init() {
	// Global persistent flags available to all subcommands
	rootCmd.PersistentFlags().StringVarP(&repoDir, "repo-dir", "C", "", "Git repository directory (default is current directory)")
}


