package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	GitCommit string
	GitDate   string
)

func main() {
	config, err := parseFlags(os.Args[1:])
	if err != nil {
		logError(os.Stderr, err)
		printUsage(os.Stderr)
		os.Exit(1)
	}
	if config.ShowVersion {
		printVersion(os.Stdout)
		return
	}

	logInfo(os.Stdout, "Validating input...")
	config, ok, err := promptForInput(config, os.Stdin, os.Stdout)
	if err != nil {
		logError(os.Stderr, err)
		os.Exit(1)
	}
	if !ok {
		logInfo(os.Stdout, "Operation canceled.")
		return
	}

	if err := validateInput(config); err != nil {
		logError(os.Stderr, err)
		os.Exit(1)
	}

	logInfo(os.Stdout, fmt.Sprintf("Generating project '%s'...", config.ProjectName))
	if err := generateProject(config, os.Stdout); err != nil {
		logError(os.Stderr, fmt.Errorf("failed to generate project: %w", err))
		os.Exit(1)
	}
	logSuccess(os.Stdout, fmt.Sprintf("Project '%s' generated successfully!", config.ProjectName))
}

func parseFlags(args []string) (Config, error) {
	flags := flag.NewFlagSet("gogen", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	git := flags.String("git", "", "Git repository URL for the template")
	local := flags.String("local", "", "Local path to the template")
	branch := flags.String("branch", "", "Git branch to clone")
	name := flags.String("name", "", "Project name to generate")
	yes := flags.Bool("yes", false, "Skip confirmation prompts")
	shortYes := flags.Bool("y", false, "Skip confirmation prompts")
	version := flags.Bool("version", false, "Print version information")

	if err := flags.Parse(args); err != nil {
		return Config{}, err
	}
	if *version {
		return Config{ShowVersion: true}, nil
	}
	if *git != "" && *local != "" {
		return Config{}, fmt.Errorf("cannot specify both --git and --local")
	}
	if *git == "" && *local == "" {
		return Config{}, fmt.Errorf("please provide a template source using --git or --local")
	}
	if *local != "" && *branch != "" {
		return Config{}, fmt.Errorf("--branch can only be used with --git")
	}

	source := *git
	if *local != "" {
		source = *local
	}

	return Config{
		ProjectName: *name,
		TemplateSrc: source,
		IsLocal:     *local != "",
		Branch:      *branch,
		Yes:         *yes || *shortYes,
	}, nil
}

func printUsage(output io.Writer) {
	fmt.Fprintln(output, "Usage: gogen --git=<template-repo-url> | --local=<template-path> [--branch=<branch>] [--name=<project-name>] [--yes]")
}

func printVersion(output io.Writer) {
	commit := GitCommit
	if commit == "" {
		commit = "unknown"
	}
	date := GitDate
	if date == "" {
		date = "unknown"
	}
	fmt.Fprintf(output, "gogen\ncommit: %s\nbuild date: %s\n", commit, date)
}
