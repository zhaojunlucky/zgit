/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"zhaojunlucky/zgit/core"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version number and build date of zgit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ZGit: https://gundamz.net/zgit")
		fmt.Println("Author: https://exia.dev")
		fmt.Printf("ZGit version: %s\n", core.Version)
		if core.BuildDate != "" {
			fmt.Printf("Build date: %s\n", core.BuildDate)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
