/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package core

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// GetCurrentBranch returns the current git branch name
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// RunGitCommand executes a git command with the given arguments
func RunGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetRepoFullName checks if the current directory is a git repository
// and returns the full repository name (e.g., "owner/repo")
func GetRepoFullName() (string, error) {
	// Check if it's a git repository
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return "", errors.New("not a git repository")
	}

	// Get the remote origin URL
	cmd = exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.New("failed to get remote origin URL")
	}

	url := strings.TrimSpace(string(output))
	
	// Parse the URL to extract owner/repo
	// Handle both SSH (git@github.com:owner/repo.git) and HTTPS (https://github.com/owner/repo.git)
	var repoFullName string
	
	if strings.HasPrefix(url, "git@") {
		// SSH format: git@github.com:owner/repo.git
		parts := strings.Split(url, ":")
		if len(parts) >= 2 {
			repoFullName = strings.TrimSuffix(parts[1], ".git")
		}
	} else if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		// HTTPS format: https://github.com/owner/repo.git
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			owner := parts[len(parts)-2]
			repo := strings.TrimSuffix(parts[len(parts)-1], ".git")
			repoFullName = owner + "/" + repo
		}
	}

	if repoFullName == "" {
		return "", errors.New("failed to parse repository name from URL")
	}

	return repoFullName, nil
}
