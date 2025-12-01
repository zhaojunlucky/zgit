/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"os"
	"text/template"
	"zhaojunlucky/zgit/core"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit [flags]...",
	Short: "Commit changes with automatic ticket prefix from branch name",
	Long: `Commit changes with automatic ticket extraction from branch name.

The commit command extracts the ticket number from your current branch name
using configured patterns and automatically formats the commit message using
the template defined in the config file.

Example:
  If you're on branch "usr/john/JIRA-1234" and run:
    zgit commit -m "fix bug"
  
  It will execute:
    git commit -m "[JIRA-1234] fix bug"

  You can also pass other git flags:
    zgit commit --amend
    zgit commit -m "fix bug" --no-verify`,
	DisableFlagParsing: true,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse args manually to find -m flag
		var message string
		var messageIndex int = -1
		var otherArgs []string
		
		for i := 0; i < len(args); i++ {
			if args[i] == "-m" || args[i] == "--message" {
				if i+1 < len(args) {
					message = args[i+1]
					messageIndex = i
					i++ // Skip the message value
				}
			} else if messageIndex == -1 || i < messageIndex || i > messageIndex+1 {
				otherArgs = append(otherArgs, args[i])
			}
		}
		
		// If -m flag is not provided, pass all args directly to git commit
		if messageIndex == -1 {
			log.Info("no -m flag provided, calling git commit directly with args")
			gitArgs := append([]string{"commit"}, args...)
			if err := core.RunGitCommand(gitArgs...); err != nil {
				log.Fatalf("failed to commit: %v", err)
			}
			log.Info("commit successful")
			return
		}
		
		log.Infof("commit called with message: %s", message)
		
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("current working directory: %s", pwd)

		// Check if current directory is a git repo and get repo full name
		repoFullName, err := core.GetRepoFullName()
		if err != nil {
			log.Fatalf("failed to get repository name: %v", err)
		}
		log.Infof("repository: %s", repoFullName)

		// Match branch
		branch, err := core.GetCurrentBranch()
		if err != nil {
			log.Fatalf("failed to get current branch: %v", err)
		}
		log.Infof("branch: %s", branch)

		// Match config
		config, err := core.LoadConfig()
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
		log.Infof("config: %v", config)

		// Match branch
		ticket, err := config.MatchBranch(repoFullName, branch)
		if err != nil {
			log.Fatalf("failed to match branch: %v", err)
		}
		log.Infof("found ticket: %s from branch %s", ticket, branch)

		// Render commit message template
		tmpl, err := template.New("commit").Parse(config.Global.Commit.Message)
		if err != nil {
			log.Fatalf("failed to parse commit message template: %v", err)
		}

		var buf bytes.Buffer
		data := map[string]string{
			"Ticket":  ticket,
			"Message": message,
		}
		if err := tmpl.Execute(&buf, data); err != nil {
			log.Fatalf("failed to render commit message template: %v", err)
		}
		commitMessage := buf.String()
		log.Infof("rendered commit message: %s", commitMessage)

		// Execute git commit with the formatted message and any additional args
		gitArgs := []string{"commit", "-m", commitMessage}
		gitArgs = append(gitArgs, otherArgs...)
		if err := core.RunGitCommand(gitArgs...); err != nil {
			log.Fatalf("failed to commit: %v", err)
		}
		log.Info("commit successful")
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
