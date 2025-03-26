package main

import (
	"fmt"
	"os"
)

func generateProject(config Config) error {
	if _, err := os.Stat(config.ProjectName); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", config.ProjectName)
	}

	var processor TemplateProcessor
	if config.IsLocal {
		processor = &LocalTemplateProcessor{}
	} else {
		processor = &GitTemplateProcessor{Branch: config.Branch}
	}

	if err := processor.Process(config.TemplateSrc, config.ProjectName, config); err != nil {
		if rmErr := os.RemoveAll(config.ProjectName); rmErr != nil {
			return fmt.Errorf("failed to process template: %v, cleanup failed: %v", err, rmErr)
		}
		return err
	}

	if err := replaceTemplateVariables(config.ProjectName, config); err != nil {
		if rmErr := os.RemoveAll(config.ProjectName); rmErr != nil {
			return fmt.Errorf("failed to replace variables: %v, cleanup failed: %v", err, rmErr)
		}
		return err
	}
	return nil
}
