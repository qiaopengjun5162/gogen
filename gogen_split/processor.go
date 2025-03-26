package main

import (
	"fmt"
	"os"
	"os/exec"
)

type TemplateProcessor interface {
	Process(src, dest string, config Config) error
}

type GitTemplateProcessor struct {
	Branch string
}

func (p *GitTemplateProcessor) Process(src, dest string, config Config) error {
	if src != config.TemplateSrc {
		return fmt.Errorf("source path mismatch")
	}
	logProgress(fmt.Sprintf("Cloning Git repository from '%s'...", src))
	args := []string{"clone", src, dest}
	if p.Branch != "" {
		args = append(args, "--branch", p.Branch)
		logProgress(fmt.Sprintf("Using branch '%s'...", p.Branch))
	}
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type LocalTemplateProcessor struct{}

func (p *LocalTemplateProcessor) Process(src, dest string, config Config) error {
	if src != config.TemplateSrc {
		return fmt.Errorf("source path mismatch")
	}
	logProgress(fmt.Sprintf("Copying local template from '%s'...", src))
	totalFiles, err := countFiles(src)
	if err != nil {
		return err
	}
	return copyDir(src, dest, totalFiles)
}
