// gogen is a project generation tool that creates new projects from Git repositories or local templates.
// It supports cloning from Git URLs with optional branch specification or copying from local paths.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var GitCommit string

// ANSI 颜色代码
const (
	ColorReset = "\033[0m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorCyan  = "\033[36m" // 信息、输入、进度
)

// Config holds the user-provided variables.
type Config struct {
	ProjectName string
	TemplateSrc string // Git URL 或本地路径
	IsLocal     bool   // 是否使用本地模板
	Branch      string // Git 分支（可选）
}

// TemplateProcessor 处理模板的接口
type TemplateProcessor interface {
	Process(src, dest string, config Config) error
}

func main() {
	if GitCommit != "" {
		fmt.Printf("Built from Git commit: %s\n", GitCommit)
	}

	// 解析命令行参数
	git := flag.String("git", "", "Git repository URL for the template")
	local := flag.String("local", "", "Local path to the template")
	branch := flag.String("branch", "", "Git branch to clone (optional)")
	flag.Parse()

	// 检查参数
	config, err := parseFlags(*git, *local, *branch)
	if err != nil {
		fmt.Printf("%s[ERROR] %v%s\n", ColorRed, err, ColorReset)
		fmt.Println("Error: Please provide a template source using --git or --local flag.")
		fmt.Println("Usage: gogen --git=<template-repo-url> | --local=<template-path>] [--branch=<branch>]")
		os.Exit(1)
	}

	// 获取用户输入
	fmt.Printf("%s[INFO] Validating input...%s\n", ColorCyan, ColorReset)
	config = getUserInput(config)

	// 验证用户输入
	if err := validateInput(config); err != nil {
		fmt.Printf("%s[ERROR] %v%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}

	// 克隆模版并生成项目
	fmt.Printf("%s[INFO] Generating project '%s'...%s\n", ColorCyan, config.ProjectName, ColorReset)
	if err := generateProject(config); err != nil {
		fmt.Printf("%s[ERROR] Failed to generate project: %v%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}
	fmt.Printf("%s[SUCCESS] Project '%s' generated successfully!%s\n", ColorGreen, config.ProjectName, ColorReset)
}

// parseFlags 解析命令行参数并返回 Config
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

// extractRepoName 从 Git URL 中提取仓库名称
func extractRepoName(gitURL string) string {
	// 匹配路径的最后一段，去除 .git 和查询参数
	re := regexp.MustCompile(`([^/?]+)(?:\.git)?(?:\?.*)?$`)
	matches := re.FindStringSubmatch(gitURL)
	if len(matches) > 1 {
		return matches[1]
	}
	return "new_project" // 兜底默认值
}

// getUserInput prompts the user for input and returns a Config struct.
func getUserInput(config Config) Config {
	scanner := bufio.NewScanner(os.Stdin)

	// 设置默认项目名称
	var defaultName string
	if !config.IsLocal {
		defaultName = extractRepoName(config.TemplateSrc) // 从 Git URL 中提取
	} else {
		base := filepath.Base(config.TemplateSrc)
		reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
		defaultName = reg.ReplaceAllString(base, "_")
	}
	if config.ProjectName == "" {
		fmt.Printf("%s[INPUT] Enter project name (default: %s): %s", ColorCyan, defaultName, ColorReset)
		scanner.Scan()
		name := strings.TrimSpace(scanner.Text())
		if name == "" {
			config.ProjectName = defaultName
		} else {
			config.ProjectName = name
		}
	}
	// 确认操作
	fmt.Printf("%s[INPUT] Generate project '%s' from %s? (Y/n): %s", ColorCyan, config.ProjectName, config.TemplateSrc, ColorReset)
	scanner.Scan()
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if response == "n" { // 只有显式输入 'n' 才取消
		fmt.Printf("%s[INFO] Operation canceled.%s\n", ColorCyan, ColorReset)
		os.Exit(0)
	}

	return config
}

// validateInput validates the git URL and project name to prevent command injection.
func validateInput(config Config) error {
	if config.IsLocal {
		// 检查本地路径是否存在
		if _, err := os.Stat(config.TemplateSrc); os.IsNotExist(err) {
			return fmt.Errorf("local template path '%s' does not exist", config.TemplateSrc)
		}
	} else {
		// 验证 Git URL 格式
		validGitURL := regexp.MustCompile(`^(https://|git@).+$`)
		if !validGitURL.MatchString(config.TemplateSrc) {
			return fmt.Errorf("invalid Git URL format: %s", config.TemplateSrc)
		}
	}

	// 验证项目名称
	validProjectName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validProjectName.MatchString(config.ProjectName) {
		return fmt.Errorf("invalid project name '%s': only letters, numbers, '_', and '-' are allowed", config.ProjectName)
	}

	return nil
}

// generateProject clones the template repository and generates the project.
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

	err := processor.Process(config.TemplateSrc, config.ProjectName, config)
	if err != nil {
		// 清理失败时的临时文件
		if rmErr := os.RemoveAll(config.ProjectName); rmErr != nil {
			return fmt.Errorf("failed to process template: %v, cleanup failed: %v", err, rmErr)
		}
		return err
	}

	// 处理模板变量替换
	err = replaceTemplateVariables(config.ProjectName, config)
	if err != nil {
		if rmErr := os.RemoveAll(config.ProjectName); rmErr != nil {
			return fmt.Errorf("failed to replace variables: %v, cleanup failed: %v", err, rmErr)
		}
		return fmt.Errorf("error replacing template variables: %v", err)
	}

	return nil
}

// replaceTemplateVariables 替换模板文件中的变量
func replaceTemplateVariables(dir string, config Config) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// 只处理文件
			content, err := os.ReadFile(filepath.Clean(path))
			if err != nil {
				return err
			}

			// 简单的字符串替换
			newContent := strings.ReplaceAll(string(content), "{{project_name}}", config.ProjectName)

			// 写回文件
			if newContent != string(content) {
				err = os.WriteFile(path, []byte(newContent), info.Mode())
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// GitTemplateProcessor 处理 Git 模板
type GitTemplateProcessor struct {
	Branch string
}

func (p *GitTemplateProcessor) Process(src, dest string, config Config) error {
	if src != config.TemplateSrc {
		return fmt.Errorf("source path '%s' does not match config template source '%s'", src, config.TemplateSrc)
	}
	fmt.Printf("%s[PROGRESS] Cloning Git repository from '%s'...%s\n", ColorCyan, src, ColorReset)
	args := []string{"clone", src, dest}
	if p.Branch != "" {
		args = append(args, "--branch", p.Branch)
		fmt.Printf("%s[PROGRESS] Using branch '%s'...%s\n", ColorCyan, p.Branch, ColorReset)
	}
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// LocalTemplateProcessor 处理本地模板
type LocalTemplateProcessor struct{}

func (p *LocalTemplateProcessor) Process(src, dest string, config Config) error {
	if src != config.TemplateSrc {
		return fmt.Errorf("source path '%s' does not match config template source '%s'", src, config.TemplateSrc)
	}
	fmt.Printf("%s[PROGRESS] Copying local template from '%s'...%s\n", ColorCyan, src, ColorReset)
	totalFiles, err := countFiles(src)
	if err != nil {
		return err
	}
	return copyDir(src, dest, totalFiles)
}

// countFiles 计算目录中的文件数，用于进度反馈
func countFiles(dir string) (int, error) {
	var count int
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}

// copyDir 递归复制目录，提供进度反馈
func copyDir(src, dest string, totalFiles int) error {
	err := os.MkdirAll(dest, 0750)
	if err != nil {
		return fmt.Errorf("error creating destination directory: %v", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("error reading source directory: %v", err)
	}

	var copiedFiles int
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, destPath, totalFiles); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
			copiedFiles++
			fmt.Printf("%s[PROGRESS] Copied %d/%d files%s\n", ColorCyan, copiedFiles, totalFiles, ColorReset)
		}
	}
	return nil
}

// copyFile 复制文件
func copyFile(src, dest string) error {
	in, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer func() {
		if err := in.Close(); err != nil {
			fmt.Printf("%s[ERROR] Failed to close input file: %v%s\n", ColorRed, err, ColorReset)
		}
	}()

	out, err := os.Create(filepath.Clean(dest))
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			fmt.Printf("%s[ERROR] Failed to close output file: %v%s\n", ColorRed, err, ColorReset)
		}
	}()

	_, err = io.Copy(out, in)
	return err
}
