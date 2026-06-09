package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func promptForInput(config Config, input io.Reader, output io.Writer) (Config, bool, error) {
	scanner := bufio.NewScanner(input)

	if config.ProjectName == "" {
		defaultName := defaultProjectName(config)
		if config.Yes {
			config.ProjectName = defaultName
			return config, true, nil
		}

		logInput(output, fmt.Sprintf("Enter project name (default: %s): ", defaultName))
		if !scanner.Scan() {
			return Config{}, false, scannerError(scanner, "project name")
		}
		name := strings.TrimSpace(scanner.Text())
		if name == "" {
			config.ProjectName = defaultName
		} else {
			config.ProjectName = name
		}
	}
	if config.Yes {
		return config, true, nil
	}

	logInput(output, fmt.Sprintf("Generate project '%s' from %s? (Y/n): ", config.ProjectName, config.TemplateSrc))
	if !scanner.Scan() {
		return Config{}, false, scannerError(scanner, "confirmation")
	}
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))

	return config, response != "n", nil
}

func scannerError(scanner *bufio.Scanner, prompt string) error {
	if err := scanner.Err(); err != nil {
		return err
	}
	return fmt.Errorf("missing %s input: %w", prompt, io.ErrUnexpectedEOF)
}

func defaultProjectName(config Config) string {
	if config.IsLocal {
		base := filepath.Base(filepath.Clean(config.TemplateSrc))
		name := nameSanitizer.ReplaceAllString(base, "_")
		if name != "" && name != "." {
			return name
		}
		return "new_project"
	}

	matches := repoNameExtract.FindStringSubmatch(config.TemplateSrc)
	if len(matches) > 1 && matches[1] != "" {
		return matches[1]
	}
	return "new_project"
}

func validateInput(config Config) error {
	if config.IsLocal {
		info, err := os.Stat(config.TemplateSrc)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("local template path '%s' does not exist", config.TemplateSrc)
			}
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("local template path '%s' is not a directory", config.TemplateSrc)
		}
	} else if !validGitURL.MatchString(config.TemplateSrc) {
		return fmt.Errorf("invalid Git URL format: %s", config.TemplateSrc)
	}

	if !validProjectName.MatchString(config.ProjectName) {
		return fmt.Errorf("invalid project name '%s': only letters, numbers, '_', and '-' are allowed", config.ProjectName)
	}

	return nil
}
