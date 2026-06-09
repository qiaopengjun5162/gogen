package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFlagsUsesLocalTemplatePath(t *testing.T) {
	config, err := parseFlags([]string{"--local=/tmp/example-template"})
	if err != nil {
		t.Fatalf("parseFlags returned error: %v", err)
	}

	if !config.IsLocal {
		t.Fatal("expected local mode")
	}
	if config.TemplateSrc != "/tmp/example-template" {
		t.Fatalf("TemplateSrc = %q, want local path", config.TemplateSrc)
	}
}

func TestParseFlagsRejectsBranchWithLocalTemplate(t *testing.T) {
	_, err := parseFlags([]string{"--local=/tmp/example-template", "--branch=main"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseFlagsSupportsAutomationOptions(t *testing.T) {
	config, err := parseFlags([]string{
		"--git=https://github.com/qiaopengjun5162/gogen.git",
		"--branch=main",
		"--name=demo",
		"--yes",
	})
	if err != nil {
		t.Fatalf("parseFlags returned error: %v", err)
	}

	if config.ProjectName != "demo" {
		t.Fatalf("ProjectName = %q, want demo", config.ProjectName)
	}
	if !config.Yes {
		t.Fatal("expected Yes=true")
	}
	if config.Branch != "main" {
		t.Fatalf("Branch = %q, want main", config.Branch)
	}
}

func TestParseFlagsSupportsShortYes(t *testing.T) {
	config, err := parseFlags([]string{"--git=https://github.com/example/template", "-y"})
	if err != nil {
		t.Fatalf("parseFlags returned error: %v", err)
	}
	if !config.Yes {
		t.Fatal("expected Yes=true")
	}
}

func TestParseFlagsSupportsVersionWithoutTemplateSource(t *testing.T) {
	config, err := parseFlags([]string{"--version"})
	if err != nil {
		t.Fatalf("parseFlags returned error: %v", err)
	}
	if !config.ShowVersion {
		t.Fatal("expected ShowVersion=true")
	}
}

func TestDefaultProjectName(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{
			name: "git repository",
			config: Config{
				TemplateSrc: "https://github.com/qiaopengjun5162/gogen.git?ref=main",
			},
			want: "gogen",
		},
		{
			name: "local path",
			config: Config{
				TemplateSrc: "/tmp/go starter!",
				IsLocal:     true,
			},
			want: "go_starter_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultProjectName(tt.config); got != tt.want {
				t.Fatalf("defaultProjectName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGenerateProjectFromLocalTemplate(t *testing.T) {
	workDir := t.TempDir()
	templateDir := filepath.Join(workDir, "template")
	if err := os.MkdirAll(filepath.Join(templateDir, "cmd"), 0750); err != nil {
		t.Fatalf("create template dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, "README.md"), []byte("name={{project_name}}\n"), 0644); err != nil {
		t.Fatalf("write README: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, "cmd", "main.go"), []byte("package main\n"), 0644); err != nil {
		t.Fatalf("write nested file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, "binary.dat"), []byte{'h', 0, '{', '{', 'p'}, 0644); err != nil {
		t.Fatalf("write binary file: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(templateDir, ".git"), 0750); err != nil {
		t.Fatalf("create .git dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(templateDir, ".git", "config"), []byte("{{project_name}}"), 0644); err != nil {
		t.Fatalf("write git config: %v", err)
	}

	t.Chdir(workDir)

	var output bytes.Buffer
	err := generateProject(Config{
		ProjectName: "demo",
		TemplateSrc: templateDir,
		IsLocal:     true,
	}, &output)
	if err != nil {
		t.Fatalf("generateProject returned error: %v", err)
	}

	readme, err := os.ReadFile(filepath.Join(workDir, "demo", "README.md"))
	if err != nil {
		t.Fatalf("read generated README: %v", err)
	}
	if string(readme) != "name=demo\n" {
		t.Fatalf("README content = %q", string(readme))
	}

	binary, err := os.ReadFile(filepath.Join(workDir, "demo", "binary.dat"))
	if err != nil {
		t.Fatalf("read generated binary: %v", err)
	}
	if !bytes.Equal(binary, []byte{'h', 0, '{', '{', 'p'}) {
		t.Fatalf("binary file was changed: %q", binary)
	}

	if _, err := os.Stat(filepath.Join(workDir, "demo", ".git")); !os.IsNotExist(err) {
		t.Fatalf(".git directory should not be copied, stat error: %v", err)
	}
	if !strings.Contains(output.String(), "Copied 3/3 files") {
		t.Fatalf("progress output = %q", output.String())
	}
}

func TestPromptForInputCanCancel(t *testing.T) {
	config, ok, err := promptForInput(
		Config{TemplateSrc: "https://github.com/example/template"},
		strings.NewReader("\nn\n"),
		&bytes.Buffer{},
	)
	if err != nil {
		t.Fatalf("promptForInput returned error: %v", err)
	}
	if ok {
		t.Fatal("expected cancellation")
	}
	if config.ProjectName != "template" {
		t.Fatalf("ProjectName = %q, want template", config.ProjectName)
	}
}

func TestPromptForInputSkipsPromptsWithYes(t *testing.T) {
	config, ok, err := promptForInput(
		Config{
			TemplateSrc: "https://github.com/example/template",
			Yes:         true,
		},
		strings.NewReader(""),
		&bytes.Buffer{},
	)
	if err != nil {
		t.Fatalf("promptForInput returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected confirmation")
	}
	if config.ProjectName != "template" {
		t.Fatalf("ProjectName = %q, want template", config.ProjectName)
	}
}

func TestPromptForInputRequiresConfirmationWhenNotYes(t *testing.T) {
	_, _, err := promptForInput(
		Config{
			ProjectName: "demo",
			TemplateSrc: "https://github.com/example/template",
		},
		strings.NewReader(""),
		&bytes.Buffer{},
	)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("error = %v, want ErrUnexpectedEOF", err)
	}
}

func TestPrintVersionUsesFallbacks(t *testing.T) {
	var output bytes.Buffer
	printVersion(&output)

	got := output.String()
	if !strings.Contains(got, "commit: unknown") {
		t.Fatalf("version output = %q", got)
	}
	if !strings.Contains(got, "build date: unknown") {
		t.Fatalf("version output = %q", got)
	}
}
