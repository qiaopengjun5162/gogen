package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type TemplateProcessor interface {
	Process(src, dest string, config Config, output io.Writer) error
}

type GitTemplateProcessor struct {
	Branch string
}

func (p *GitTemplateProcessor) Process(src, dest string, config Config, output io.Writer) error {
	if src != config.TemplateSrc {
		return fmt.Errorf("source path '%s' does not match config template source '%s'", src, config.TemplateSrc)
	}

	logProgress(output, fmt.Sprintf("Cloning Git repository from '%s'...", src))
	args := []string{"clone"}
	if p.Branch != "" {
		args = append(args, "--branch", p.Branch)
		logProgress(output, fmt.Sprintf("Using branch '%s'...", p.Branch))
	}
	args = append(args, src, dest)

	cmd := exec.Command("git", args...)
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type LocalTemplateProcessor struct{}

func (p *LocalTemplateProcessor) Process(src, dest string, config Config, output io.Writer) error {
	if src != config.TemplateSrc {
		return fmt.Errorf("source path '%s' does not match config template source '%s'", src, config.TemplateSrc)
	}

	logProgress(output, fmt.Sprintf("Copying local template from '%s'...", src))
	return copyDir(src, dest, output)
}
