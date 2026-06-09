package main

import "regexp"

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
)

type Config struct {
	ProjectName string
	TemplateSrc string
	IsLocal     bool
	Branch      string
	Vars        map[string]string
	Yes         bool
	ShowVersion bool
}

var (
	validGitURL      = regexp.MustCompile(`^(https://|git@).+$`)
	validProjectName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	validVariableKey = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	repoNameExtract  = regexp.MustCompile(`([^/?]+?)(?:\.git)?(?:\?.*)?$`)
	nameSanitizer    = regexp.MustCompile(`[^a-zA-Z0-9_-]`)
)
