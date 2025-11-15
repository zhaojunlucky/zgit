/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"zhaojunlucky/zgit/cmd"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Configure logrus with timestamp
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	cmd.Execute()
}
