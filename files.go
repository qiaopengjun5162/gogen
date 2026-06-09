package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func copyDir(src, dest string, config Config, output io.Writer) error {
	totalFiles, err := countFiles(src)
	if err != nil {
		return err
	}

	var copiedFiles int
	targets := make(map[string]string)
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return os.MkdirAll(dest, 0750)
		}
		if entry.IsDir() && entry.Name() == gitDirName {
			return filepath.SkipDir
		}

		targetPath := filepath.Join(dest, replacePathVariables(relPath, config))
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if err := rememberTarget(targets, path, targetPath); err != nil {
				return err
			}
			return os.MkdirAll(targetPath, info.Mode().Perm())
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		if err := rememberTarget(targets, path, targetPath); err != nil {
			return err
		}
		if err := copyFile(path, targetPath, info.Mode()); err != nil {
			return err
		}

		copiedFiles++
		logProgress(output, "Copied "+formatProgress(copiedFiles, totalFiles)+" files")
		return nil
	})
}

func countFiles(dir string) (int, error) {
	var count int
	err := filepath.WalkDir(dir, func(_ string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && entry.Name() == gitDirName {
			return filepath.SkipDir
		}
		if !entry.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}

func copyFile(src, dest string, mode os.FileMode) (err error) {
	in, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := in.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	out, err := os.OpenFile(filepath.Clean(dest), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := out.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	_, err = io.Copy(out, in)
	return err
}

func replaceTemplateVariables(dir string, config Config) error {
	return filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == gitDirName {
				return filepath.SkipDir
			}
			return nil
		}

		content, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return err
		}
		if isBinary(content) {
			return nil
		}

		newContent := replaceVariables(string(content), config)
		if newContent == string(content) {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Clean(path), []byte(newContent), info.Mode().Perm())
	})
}

func isBinary(content []byte) bool {
	const sniffSize = 8000
	if len(content) > sniffSize {
		content = content[:sniffSize]
	}
	for _, b := range content {
		if b == 0 {
			return true
		}
	}
	return false
}

func formatProgress(done, total int) string {
	return strconv.Itoa(done) + "/" + strconv.Itoa(total)
}

func replaceVariables(content string, config Config) string {
	replacements := []string{"{{project_name}}", config.ProjectName}
	for key, value := range config.Vars {
		replacements = append(replacements, "{{"+key+"}}", value)
	}
	return strings.NewReplacer(replacements...).Replace(content)
}

func replacePathVariables(path string, config Config) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for index, part := range parts {
		parts[index] = replaceVariables(part, config)
	}
	return filepath.FromSlash(strings.Join(parts, "/"))
}

func rememberTarget(targets map[string]string, sourcePath, targetPath string) error {
	cleanTarget := filepath.Clean(targetPath)
	if existingSource, exists := targets[cleanTarget]; exists {
		return fmt.Errorf("template paths %q and %q both resolve to %q", existingSource, sourcePath, cleanTarget)
	}
	targets[cleanTarget] = sourcePath
	return nil
}
