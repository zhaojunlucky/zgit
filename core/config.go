// Package core
package core

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/cfg"
	"gopkg.in/yaml.v3"
)

var configFile = "config.yaml"
var config *Config

// Config represents the zgit configuration structure
type Config struct {
	Global GlobalConfig `yaml:"global"`
	Repos  []RepoConfig `yaml:"repos"`
}

// GlobalConfig represents global configuration settings
type GlobalConfig struct {
	Branches []string     `yaml:"branches"`
	Commit   CommitConfig `yaml:"commit"`
}

// CommitConfig represents commit message configuration
type CommitConfig struct {
	Message string `yaml:"message"`
}

// RepoConfig represents repository-specific configuration
type RepoConfig struct {
	Name     string   `yaml:"name"`
	Branches []string `yaml:"branches"`
}

func (c *Config) MatchBranch(repoName, branch string) (string, error) {
	for _, repo := range c.Repos {
		if repo.Name != repoName {
			continue
		}
		for _, branchPattern := range repo.Branches {
			ticket, err := c.matchBranch(repo.Name, branchPattern, branch)
			if err != nil {
				return "", err
			} else if len(ticket) > 0 {
				return ticket, nil
			}
		}
	}

	for _, branchPattern := range c.Global.Branches {
		ticket, err := c.matchBranch("global", branchPattern, branch)
		if err != nil {
			return "", err
		} else if len(ticket) > 0 {
			return ticket, nil
		}
	}
	return "", errors.New("branch not found or ticket group not matched")
}

func (c *Config) matchBranch(repoName, branchPattern, branch string) (string, error) {
	if !strings.Contains(branchPattern, "?P<ticket>") {
		log.Errorf("branch pattern %s of %s must contain ?P<ticket>", branchPattern, repoName)
		return "", errors.New("branch pattern must contain ?P<ticket>")
	}
	reg := regexp.MustCompile(branchPattern)
	match := reg.FindStringSubmatch(branch)

	if match == nil {
		return "", nil
	}

	// Check if regex has a named group "ticket"
	groupNames := reg.SubexpNames()
	for i, name := range groupNames {
		if name == "ticket" && i < len(match) {
			return match[i], nil
		}
	}
	return "", nil
}

// LoadConfig loads the zgit configuration from the specified file path
func LoadConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	cfgFilePaths := cfg.GetCfgPath("zgit", configFile)
	for _, cfgFilePath := range cfgFilePaths {
		log.Infof("try loading %s", cfgFilePath)
		data, err := os.ReadFile(cfgFilePath)
		if err != nil {
			log.Warnf("failed to read %s: %v", cfgFilePath, err)
			continue
		}

		config = &Config{}
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, err
		}

		err = config.validate()
		if err != nil {
			return nil, err
		}
		log.Infof("used config file: %s", cfgFilePath)
		return config, nil
	}
	return nil, errors.New("config file not found")
}

// validate checks if the config is valid
func (c *Config) validate() error {
	// Check if commit message template contains {{.Ticket}} and {{.Message}}
	// Allow optional spaces: {{ .Ticket }} or {{.Ticket}}
	messageTemplate := c.Global.Commit.Message
	ticketRegex := regexp.MustCompile(`\{\{\s*\.Ticket\s*\}\}`)
	messageRegex := regexp.MustCompile(`\{\{\s*\.Message\s*\}\}`)

	if !ticketRegex.MatchString(messageTemplate) {
		return errors.New("commit message template must contain {{.Ticket}} or {{ .Ticket }}")
	}
	if !messageRegex.MatchString(messageTemplate) {
		return errors.New("commit message template must contain {{.Message}} or {{ .Message }}")
	}

	// Check if at least one global branch pattern exists
	if len(c.Global.Branches) == 0 {
		return errors.New("at least one global branch pattern must be defined")
	}

	// If repos are defined, check that each repo has at least one branch pattern
	if len(c.Repos) > 0 {
		for _, repo := range c.Repos {
			if len(repo.Branches) == 0 {
				return fmt.Errorf("repository '%s' must have at least one branch pattern", repo.Name)
			}
		}
	}

	return nil
}
