/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var remoteName string

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the GitHub repository in your browser",
	Long: `Open the GitHub repository URL in your default web browser.

By default, it uses the 'origin' remote. You can specify a different remote
using the -r or --remote flag.

Examples:
  zgit open              # Opens the origin remote's GitHub page
  zgit open -r upstream  # Opens the upstream remote's GitHub page`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := getRemoteURL(remoteName)
		if err != nil {
			log.Fatalf("failed to get remote URL: %v", err)
		}

		webURL, err := parseGitURLToWeb(url)
		if err != nil {
			log.Fatalf("failed to parse git URL: %v", err)
		}

		log.Infof("opening %s", webURL)
		if err := openBrowser(webURL); err != nil {
			log.Fatalf("failed to open browser: %v", err)
		}
	},
}

// getRemoteURL gets the URL for the specified remote
func getRemoteURL(remote string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", remote)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("remote '%s' not found", remote)
	}
	return strings.TrimSpace(string(output)), nil
}

// parseGitURLToWeb converts a git remote URL to a web browser URL
// Supports:
//   - SSH: git@github.com:owner/repo.git
//   - HTTPS: https://github.com/owner/repo.git
//   - SSH with ssh:// prefix: ssh://git@github.com/owner/repo.git
func parseGitURLToWeb(gitURL string) (string, error) {
	gitURL = strings.TrimSpace(gitURL)
	gitURL = strings.TrimSuffix(gitURL, ".git")

	// SSH format: git@github.com:owner/repo
	if strings.HasPrefix(gitURL, "git@") {
		// git@github.com:owner/repo -> https://github.com/owner/repo
		re := regexp.MustCompile(`^git@([^:]+):(.+)$`)
		matches := re.FindStringSubmatch(gitURL)
		if len(matches) == 3 {
			host := matches[1]
			path := matches[2]
			return fmt.Sprintf("https://%s/%s", host, path), nil
		}
	}

	// SSH with ssh:// prefix: ssh://git@github.com/owner/repo
	if strings.HasPrefix(gitURL, "ssh://") {
		// ssh://git@github.com/owner/repo -> https://github.com/owner/repo
		re := regexp.MustCompile(`^ssh://git@([^/]+)/(.+)$`)
		matches := re.FindStringSubmatch(gitURL)
		if len(matches) == 3 {
			host := matches[1]
			path := matches[2]
			return fmt.Sprintf("https://%s/%s", host, path), nil
		}
	}

	// HTTPS format: https://github.com/owner/repo
	if strings.HasPrefix(gitURL, "https://") || strings.HasPrefix(gitURL, "http://") {
		return gitURL, nil
	}

	return "", fmt.Errorf("unsupported git URL format: %s", gitURL)
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().StringVarP(&remoteName, "remote", "r", "origin", "Remote name to open (default: origin)")
}
