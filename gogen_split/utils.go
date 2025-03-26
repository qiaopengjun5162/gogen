package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func handleError(err error, msg string) {
	if err != nil {
		if msg != "" {
			fmt.Printf("%s[ERROR] %s: %v%s\n", ColorRed, msg, err, ColorReset)
		} else {
			fmt.Printf("%s[ERROR] %v%s\n", ColorRed, err, ColorReset)
		}
		os.Exit(1)
	}
}

func logInfo(msg string) {
	fmt.Printf("%s[INFO] %s%s\n", ColorCyan, msg, ColorReset)
}

func logSuccess(msg string) {
	fmt.Printf("%s[SUCCESS] %s%s\n", ColorGreen, msg, ColorReset)
}

func logInput(msg string) {
	fmt.Printf("%s[INPUT] %s%s", ColorCyan, msg, ColorReset)
}

func logProgress(msg string) {
	fmt.Printf("%s[PROGRESS] %s%s\n", ColorCyan, msg, ColorReset)
}

func countFiles(dir string) (int, error) {
	var count int
	err := filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
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

func copyDir(src, dest string, totalFiles int) error {
	if err := os.MkdirAll(dest, 0750); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
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
			logProgress(fmt.Sprintf("Copied %d/%d files", copiedFiles, totalFiles))
		}
	}
	return nil
}

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

func replaceTemplateVariables(dir string, config Config) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(filepath.Clean(path))
			if err != nil {
				return err
			}
			newContent := strings.ReplaceAll(string(content), "{{project_name}}", config.ProjectName)
			if newContent != string(content) {
				return os.WriteFile(path, []byte(newContent), info.Mode())
			}
		}
		return nil
	})
}
