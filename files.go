package main

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func copyDir(src, dest string, output io.Writer) error {
	totalFiles, err := countFiles(src)
	if err != nil {
		return err
	}

	var copiedFiles int
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
		if entry.IsDir() && entry.Name() == ".git" {
			return filepath.SkipDir
		}

		targetPath := filepath.Join(dest, relPath)
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			return os.MkdirAll(targetPath, info.Mode().Perm())
		}

		info, err := entry.Info()
		if err != nil {
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
		if entry.IsDir() && entry.Name() == ".git" {
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
			if entry.Name() == ".git" {
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

		newContent := strings.ReplaceAll(string(content), "{{project_name}}", config.ProjectName)
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
