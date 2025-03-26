package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func parseFlags(git, local, branch string) (Config, error) {
	if git != "" && local != "" {
		return Config{}, fmt.Errorf("cannot specify both --git and --local")
	}
	if git == "" && local == "" {
		return Config{}, fmt.Errorf("please provide a template source using --git or --local")
	}
	return Config{
		TemplateSrc: git,
		IsLocal:     local != "",
		Branch:      branch,
	}, nil
}

func extractRepoName(gitURL string) string {
	matches := RepoNameExtract.FindStringSubmatch(gitURL)
	if len(matches) > 1 {
		return matches[1]
	}
	return "new_project"
}

func getUserInput(config Config) Config {
	scanner := bufio.NewScanner(os.Stdin)
	defaultName := config.ProjectName
	if defaultName == "" {
		if !config.IsLocal {
			defaultName = extractRepoName(config.TemplateSrc)
		} else {
			base := filepath.Base(config.TemplateSrc)
			defaultName = strings.ReplaceAll(base, "[^a-zA-Z0-9_-]", "_")
		}
	}

	logInput(fmt.Sprintf("Enter project name (default: %s): ", defaultName))
	scanner.Scan()
	if name := strings.TrimSpace(scanner.Text()); name != "" {
		config.ProjectName = name
	} else {
		config.ProjectName = defaultName
	}

	logInput(fmt.Sprintf("Generate project '%s' from %s? (Y/n): ", config.ProjectName, config.TemplateSrc))
	scanner.Scan()
	if strings.ToLower(strings.TrimSpace(scanner.Text())) == "n" {
		logInfo("Operation canceled.")
		os.Exit(0)
	}
	return config
}

func validateInput(config Config) error {
	if config.IsLocal {
		if _, err := os.Stat(config.TemplateSrc); os.IsNotExist(err) {
			return fmt.Errorf("local template path '%s' does not exist", config.TemplateSrc)
		}
	} else {
		if !ValidGitURL.MatchString(config.TemplateSrc) {
			return fmt.Errorf("invalid Git URL format: %s", config.TemplateSrc)
		}
	}
	if !ValidProjectName.MatchString(config.ProjectName) {
		return fmt.Errorf("invalid project name '%s': only letters, numbers, '_', and '-' are allowed", config.ProjectName)
	}
	return nil
}
