package main

import (
	"flag"
	"fmt"
)

var GitCommit string

func main() {
	if GitCommit != "" {
		fmt.Printf("Built from Git commit: %s\n", GitCommit)
	}

	// 解析命令行参数
	git := flag.String("git", "", "Git repository URL for the template")
	local := flag.String("local", "", "Local path to the template")
	branch := flag.String("branch", "", "Git branch to clone (optional)")
	flag.Parse()

	config, err := parseFlags(*git, *local, *branch)
	if err != nil {
		handleError(err, "Please provide a template source using --git or --local flag.")
	}

	// 获取用户输入并验证
	config = getUserInput(config)
	if err := validateInput(config); err != nil {
		handleError(err, "")
	}

	// 生成项目
	logInfo(fmt.Sprintf("Generating project '%s'...", config.ProjectName))
	if err := generateProject(config); err != nil {
		handleError(err, "Failed to generate project")
	}
	logSuccess(fmt.Sprintf("Project '%s' generated successfully!", config.ProjectName))
}
