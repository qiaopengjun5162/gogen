package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Config holds the user-provided variables.
type Config struct {
	ProjectName string
}

func main() {
	// 解析命令行参数
	git := flag.String("git", "", "Git repository URL for the template")
	flag.Parse()

	// 检查必需参数
	if *git == "" {
		fmt.Println("Error: Please provide a Git repository URL for the template using the --git flag.")
		fmt.Println("Usage: gogen --git=<template-repo-url>")
		os.Exit(1)
	}

	// 获取用户输入
	config := getUserInput()

	// 验证用户输入
	if err := validateInput(*git, config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// 克隆模版并生成项目
	err := generateProject(*git, config)
	if err != nil {
		fmt.Printf("Error generating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Project '%s' generated successfully!\n", config.ProjectName)
}

// getUserInput prompts the user for input and returns a Config struct.
func getUserInput() Config {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter project name: ")
	scanner.Scan()
	projectName := strings.TrimSpace(scanner.Text())

	return Config{
		ProjectName: projectName,
	}
}

// validateInput validates the git URL and project name to prevent command injection.
func validateInput(git string, config Config) error {
	// 验证 Git URL 格式
	validGitURL := regexp.MustCompile(`^(https://|git@).+$`)
	if !validGitURL.MatchString(git) {
		return fmt.Errorf("invalid Git URL format")
	}

	// 验证项目名称
	validProjectName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validProjectName.MatchString(config.ProjectName) {
		return fmt.Errorf("invalid project name")
	}

	return nil
}

// generateProject clones the template repository and generates the project.
func generateProject(git string, config Config) error {
	// 克隆模版到指定目录
	cmd := exec.Command("git", "clone", git, config.ProjectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error cloning template: %v", err)
	}
	return nil
}
