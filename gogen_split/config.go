package main

import "regexp"

// ANSI 颜色代码
const (
	ColorReset = "\033[0m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorCyan  = "\033[36m"
)

// Config holds the user-provided variables.
type Config struct {
	ProjectName string
	TemplateSrc string
	IsLocal     bool
	Branch      string
}

// 正则表达式常量
var (
	ValidGitURL      = regexp.MustCompile(`^(https://|git@).+$`)
	ValidProjectName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	RepoNameExtract  = regexp.MustCompile(`([^/?]+)(?:\.git)?(?:\?.*)?$`)
)
