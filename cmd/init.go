/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const configURL = "https://raw.githubusercontent.com/zhaojunlucky/zgit/main/config.yaml"

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize zgit configuration",
	Long: `Download the default config.yaml from GitHub and save it to ~/.config/zgit/config.yaml.
	
If the configuration file already exists, you will be prompted to confirm whether to override it.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to get home directory: %v", err)
		}

		// Create config directory path
		configDir := filepath.Join(homeDir, ".config", "zgit")
		configPath := filepath.Join(configDir, "config.yaml")

		// Check if config file already exists
		if _, err := os.Stat(configPath); err == nil {
			// File exists, ask user for confirmation
			fmt.Printf("Config file already exists at %s\n", configPath)
			fmt.Print("Do you want to override it? (y/N): ")
			
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("failed to read user input: %v", err)
			}
			
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				log.Info("Init cancelled")
				return
			}
		}

		// Create config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Fatalf("failed to create config directory: %v", err)
		}

		// Download config file from GitHub
		log.Infof("Downloading config from %s", configURL)
		resp, err := http.Get(configURL)
		if err != nil {
			log.Fatalf("failed to download config: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("failed to download config: HTTP %d", resp.StatusCode)
		}

		// Create the config file
		file, err := os.Create(configPath)
		if err != nil {
			log.Fatalf("failed to create config file: %v", err)
		}
		defer file.Close()

		// Write the downloaded content to file
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Fatalf("failed to write config file: %v", err)
		}

		log.Infof("Config file successfully created at %s", configPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
