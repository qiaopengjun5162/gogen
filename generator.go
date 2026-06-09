package main

import (
	"fmt"
	"io"
	"os"
)

func generateProject(config Config, output io.Writer) error {
	if _, err := os.Stat(config.ProjectName); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", config.ProjectName)
	}

	var processor TemplateProcessor
	if config.IsLocal {
		processor = &LocalTemplateProcessor{}
	} else {
		processor = &GitTemplateProcessor{Branch: config.Branch}
	}

	if err := processor.Process(config.TemplateSrc, config.ProjectName, config, output); err != nil {
		return cleanupGeneratedProject(config.ProjectName, err)
	}

	if err := replaceTemplateVariables(config.ProjectName, config); err != nil {
		return cleanupGeneratedProject(config.ProjectName, fmt.Errorf("error replacing template variables: %w", err))
	}

	return nil
}

func cleanupGeneratedProject(path string, cause error) error {
	if rmErr := os.RemoveAll(path); rmErr != nil {
		return fmt.Errorf("%v; cleanup failed: %w", cause, rmErr)
	}
	return cause
}
