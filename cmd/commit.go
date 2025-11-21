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
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		
		// If -m flag is not provided, pass all args directly to git commit
		if !cmd.Flags().Changed("message") {
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
		gitArgs = append(gitArgs, args...)
		if err := core.RunGitCommand(gitArgs...); err != nil {
			log.Fatalf("failed to commit: %v", err)
		}
		log.Info("commit successful")
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("message", "m", "", "Commit message (optional, if not provided, git commit will be called directly)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
